package investorCheck

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"strings"

	"bitbucket.org/go-mis/modules/cif"
	email "bitbucket.org/go-mis/modules/email"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/modules/investor"
	"net/http"
	"bitbucket.org/go-mis/config"
)

type totalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchDatatables -  fetch data based on parameters sent by datatables
func FetchDatatables(ctx *iris.Context) {
	var dtTables DataTable
	var totalRows int
	investors := []InvestorCheck{}
	limit := ctx.URLParam("limit")
	page := ctx.URLParam("page")
	filterBy := ctx.URLParam("filterBy")

	// Get object data tables from url
	dtTables = ParseDatatableURI(string(ctx.URI().FullURI()))
	search := dtTables.Search.Value
	orderBy := "NAME"
	orderDir := "ASC"

	fmt.Println("Search keyword: ", search)

	if len(dtTables.Columns) > 0 && len(dtTables.OrderInfo) > 0 {
		orderBy = dtTables.Columns[dtTables.OrderInfo[0].Column].Data
		orderDir = dtTables.OrderInfo[0].Dir
	}
	fmt.Println("Orderby: ", orderBy)
	fmt.Println("Orderdir: ", orderDir)
	query := ` SELECT cif.id, cif."name", 
					cif."phoneNo", 
					cif."idCardNo", 
					investor."bankAccountName",
					cif."taxCardNo", 
					cif."idCardFilename", 
					cif."taxCardFilename", 
					cif."idCardNo", 
					cif."isValidated", 
					cif."taxCardNo", 
					array_to_string(array_agg(virtual_account."bankName"),',') as "virtualAccountBankName", 
					array_to_string(array_agg(virtual_account."virtualAccountNo"),',') as "virtualAccountNumber", 
					investor."investorNo", 
					investor."createdAt",
					cif."isActivated",
					cif."isVerified",
					cif."isDeclined",
					cif."username",
					cif."activationDate",
					cif."declinedDate",
					cif."username",
					cif."activationDate",
					cif."declinedDate",
					count(*) OVER() AS full_count
				FROM investor 
					LEFT JOIN r_investor_virtual_account ON r_investor_virtual_account."investorId" = investor.id 
					LEFT JOIN virtual_account ON virtual_account.id = r_investor_virtual_account."vaId" 
					JOIN r_cif_investor ON r_cif_investor."investorId" = investor.id 
					JOIN cif ON cif.id = r_cif_investor."cifId" 
				WHERE cif."isVerified" = FALSE 
				AND cif."idCardFilename" IS NOT NULL 
				
				AND cif."deletedAt" isnull AND virtual_account."deletedAt" isnull `

	if len(strings.TrimSpace(filterBy)) > 0 {
		switch strings.ToUpper(filterBy) {
		case "ACTIVATED":
			query += ` AND cif."isActivated" = TRUE AND (cif."isDeclined" = FALSE OR cif."isDeclined" ISNULL) `
		case "DECLINED":
			query += ` AND cif."isActivated" = FALSE AND (cif."isDeclined" = TRUE OR cif."isDeclined" IS NOT NULL ) `
		default:
			query += ` AND cif."isActivated" = TRUE `
		}
	}

	if len(strings.TrimSpace(search)) > 0 {
		query += ` AND (cif.name ~* '` + search + `' OR investor."investorNo"::text ~* '` + search + `' OR cif."idCardNo" ~* '` + search + `' OR  
					cif."taxCardNo" ~* '` + search + `' OR cif."username"  ~* '` + search + `') `
	} else {
		query += "AND cif.name ~* '" + search + "' "
	}

	groupedBy := ` group by cif."id", cif."name", cif."phoneNo", cif."idCardNo", "bankAccountName", cif."taxCardNo",
	cif."idCardNo", cif."taxCardNo", cif."idCardFilename", cif."taxCardFilename", cif."isValidated",
	investor."investorNo", investor."createdAt" `
	query += groupedBy
	if len(strings.TrimSpace(orderBy)) > 0 && len(strings.TrimSpace(orderDir)) > 0 {
		switch strings.ToUpper(orderBy) {
		case "INVESTORNO":
			query += ` ORDER BY investor."investorNo" ` + orderDir
		case "NAME":
			query += ` ORDER BY cif."name" ` + orderDir
		case "IDCARDNO":
			query += ` ORDER BY cif."idCardNo" ` + orderDir
		case "EMAIL":
			query += ` ORDER BY cif.username ` + orderDir
		case "ACTIVATIONDATE":
			query += ` ORDER BY cif."activationDate" ` + orderDir
		case "DECLINEDDATE":
			query += ` ORDER BY cif."declinedDate" ` + orderDir
		case "STATUS":
			query += ` ORDER BY ( 
							CASE
								 WHEN cif."isDeclined" THEN 1 
								 WHEN NOT cif."isDeclined" THEN 2
								 ELSE 3
							END 
						) ` + orderDir
		case "CREATEDAT":
			query += ` ORDER BY investor."createdAt" ` + orderDir
		default:
			query += ` ORDER BY cif."name" ASC `
		}
	}
	if len(strings.TrimSpace(limit)) == 0 {
		query += ` LIMIT 10 `
	} else {
		query += ` LIMIT ` + limit
	}
	if len(strings.TrimSpace(page)) == 0 {
		query += ` OFFSET 0 `
	} else {
		query += ` OFFSET ` + page
	}
	services.DBCPsql.Raw(query).Scan(&investors)
	for idx := range investors {
		declined := investors[idx].IsDeclined
		if declined != nil && *investors[idx].IsDeclined {
			investors[idx].Status = "declined"
		} else {
			investors[idx].Status = "activated"
		}
	}
	if len(investors) > 0 {
		totalRows = investors[0].RowsFullCount
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"data":      investors,
		"totalRows": totalRows,
	})
}

// Validate - verify the selected investor
func Validate(ctx *iris.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	// status type: verified or declined
	status := ctx.Param("status")

	cifSchema := cif.Cif{}
	services.DBCPsql.Table("cif").Where("id = ?", id).Scan(&cifSchema)

	if status == "validated" {
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isValidated", true)
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isVerified", true)
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isDeclined", false)
		services.DBCPsql.Exec(`update cif set "validationDate"=current_timestamp where id=?`,id)

		// get investor id
		inv := &r.RCifInvestor{}
		services.DBCPsql.Table("r_cif_investor").Where("\"cifId\" = ?", id).Scan(&inv)

		// get investor id
		invNo := &InvestorNumber{}
		services.DBCPsql.Raw(`select id,"investorNo" from investor where id = ?`, inv.InvestorId).Scan(invNo)

		// get virtual account
		rInvVa := []r.RInvestorVirtualAccount{}
		services.DBCPsql.Table("r_investor_virtual_account").Where("\"investorId\" = ?", inv.InvestorId).Scan(&rInvVa)

		vaData := make(map[string]string)

		vaData["MANDIRI"] = "88000" + strconv.Itoa(invNo.InvestorNo)
		vaData["MANDIRI_HOLDER"] = cifSchema.Name

		vaData["BCA"] = "10036" + strconv.Itoa(invNo.InvestorNo)
		vaData["BCA_HOLDER"] = cifSchema.Name

		// get investor data
		investorSchema := investor.Investor{}
		services.DBCPsql.Table("investor").Where("id = ?", inv.InvestorId).Scan(&investorSchema)

		// send create BCA VA request
		params := strings.NewReader(`{"investorNo":` + strconv.Itoa(invNo.InvestorNo) + `}`)
		request, err := http.NewRequest("POST", config.GoBankingPath+`/bca/register-va`, params)
		if err != nil {
			fmt.Println(err)
		}

		request.Header.Set("X-Auth-Token", "AMARTHA123")

		client := &http.Client{}
		_, errResp := client.Do(request)
		if errResp != nil {
			fmt.Println(errResp)
		}

		if cifSchema.Username != "" {
			fmt.Println("Sending email..")
			//go email.SendEmailVerificationSuccess(cifSchema.Username, cifSchema.Name, vaData["BCA"], vaData["BCA_HOLDER"], vaData["MANDIRI"], vaData["MANDIRI_HOLDER"])
			go email.SendEmailVerificationSuccess("wuri.wulandari@amartha.com", cifSchema.Name, vaData["BCA"], vaData["BCA_HOLDER"], vaData["MANDIRI"], vaData["MANDIRI_HOLDER"])
		}

		if cifSchema.PhoneNo != "" {
			// send sms notification
			fmt.Println("Sending sms ... ")
			message := "Selamat data Anda sudah terverifikasi. Silakan login ke dashboard Anda dan mulai berinvestasi. www.amartha.com \n\nAmartha"
			//services.SendSMS(cifSchema.PhoneNo, message)
			services.SendSMS("+628119780880", message)
		}

	} else if strings.ToUpper(status) == "DECLINED" {
		// set isDeclined true
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isValidated", false)
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isVerified", false)
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isDeclined", true)
		date:=time.Now()
		fmt.Println("DateInvestorCheck",date)
		services.DBCPsql.Exec(`update cif set "declinedDate"=current_timestamp where id=?`,id)

		// Decline will send an email to investor
		payload := struct {
			Reasons []string `json:"reasons"`
		}{}
		ctx.ReadJSON(&payload)
		fmt.Println("Decline reasons: ", payload.Reasons)
		//go email.SendEmailVerificationFailed(cifSchema.Username, cifSchema.Name, payload.Reasons)
		go email.SendEmailVerificationFailed("wuri.wulandari@amartha.com", cifSchema.Name, payload.Reasons)
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":             "success",
		"verificationStatus": status,
	})
}

func rightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen] + padStr
}

func ParseDatatableURI(fullURI string) DataTable {
	var dtTables DataTable
	u, _ := url.Parse(fullURI)
	q := u.Query()
	for k, v := range q {
		if len(strings.TrimSpace(v[0])) == 0 {
			err := json.Unmarshal([]byte(k), &dtTables)
			if err == nil {
				return dtTables
			}

		}
	}

	return dtTables
}
