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
	va "bitbucket.org/go-mis/modules/virtual-account"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

type totalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchDatatables -  fetch data based on parameters sent by datatables
func FetchDatatables(ctx *iris.Context) {
	var dtTables DataTable
	investors := []InvestorCheck{}
	limit := ctx.URLParam("limit")
	page := ctx.URLParam("page")
	filterBy := ctx.URLParam("filterBy")

	// Get object data tables from url
	dtTables = ParseDatatableURI(string(ctx.URI().FullURI()))

	search := dtTables.Search.Value
	orderBy := dtTables.Columns[dtTables.OrderInfo[0].Column].Data
	orderDir := dtTables.OrderInfo[0].Dir

	fmt.Println("Search: ", search)
	fmt.Println("OrderBy: ", orderBy)
	fmt.Println("OrderDir: ", orderDir)
	fmt.Println("FilterBy: ", filterBy)
	fmt.Println("Limit: ", limit)
	fmt.Println("Page: ", page)

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
					cif."declinedDate"
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

	if search != "" {
		query += ` AND (cif.name ILIKE '%` + search + `%' OR investor."investorNo"::text ILIKE '%` + search + `%' OR cif."idCardNo" ILIKE '%` + search + `%' OR  
					cif."taxCardNo" ILIKE '%` + search + `' OR cif."username" ILIKE '%` + search + `%') `
	}

	if search != "" {
		query += "AND cif.name ~* '" + ctx.URLParam("search") + "' "
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
			query += `ORDER BY cif.username ` + orderDir
		case "ACTIVATIONDATE":
			query += `ORDER BY cif."activationDate" ` + orderDir
		case "DECLINEDATE":
			query += `ORDER BY cif."declinedDate" ` + orderDir
		case "REGISTRATIONDATE":
			query += `ORDER BY investor."createdAt" ` + orderDir
		default:
			query += ` ORDER BY cif."name" ASC`
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

	for idx, val := range investors {
		d := val.IsDeclined

		if d != nil && *d {
			investors[idx].Status = "declined"
		} else {
			investors[idx].Status = "activated"
		}
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"data":      investors,
		"totalRows": len(investors),
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

		/** FOR PRODUCTION, PLEASE UNCOMMENT
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
			go email.SendEmailVerificationSuccess(cifSchema.Username, cifSchema.Name, vaData["BCA"], vaData["BCA_HOLDER"], vaData["MANDIRI"], vaData["MANDIRI_HOLDER"])
		}

		if cifSchema.PhoneNo != "" {
			// send sms notification
			fmt.Println("Sending sms ... ")
			twilio := services.InitTwilio()
			message := "Selamat data Anda sudah terverifikasi. Silakan login ke dashboard Anda dan mulai berinvestasi. www.amartha.com"
			twilio.SetParam(cifSchema.PhoneNo, message)
			// twilio.SetParam("+628992548716", message)
			twilio.SendSMS()
		}
		*/

	} else {
		// set isDeclined true
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isValidated", false)
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isVerified", false)
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isDeclined", true)

		services.DBCPsql.Table("cif").Where("id = ?", id).Update("declinedDate", time.Now())

		/* FOR PROUDCTION, pleas uncomment
		email.SendEmailVerificationFailed(cifSchema.Username, cifSchema.Name, "Declined")
		*/
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

// Verified - verify the selected investor
func Verified(ctx *iris.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	// status type: verified or declined
	cifSchema := cif.Cif{}
	services.DBCPsql.Table("cif").Where("id = ?", id).Scan(&cifSchema)

	if *cifSchema.IsValidated == true {
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isVerified", true)

		// get investor id
		inv := &r.RCifInvestor{}
		services.DBCPsql.Table("r_cif_investor").Where("\"cifId\" = ?", id).Scan(&inv)

		// get virtual account
		rInvVa := []r.RInvestorVirtualAccount{}
		services.DBCPsql.Table("r_investor_virtual_account").Where("\"investorId\" = ?", inv.InvestorId).Scan(&rInvVa)

		vaObj := &va.VirtualAccount{}
		userVa := []va.VirtualAccount{}
		for _, val := range rInvVa {
			services.DBCPsql.Table("virtual_account").Where("\"id\" = ?", val.VirtualAccountId).Scan(&vaObj)
			userVa = append(userVa, *vaObj)
		}

		vaData := make(map[string]string)
		for _, val := range userVa {
			if val.BankName == "BRI" {
				vaData["BRI"] = val.VirtualAccountNo
				vaData["BRI_HOLDER"] = val.VirtualAccountName
			} else if val.BankName == "BCA" {
				vaData["BCA"] = val.VirtualAccountNo
				vaData["BCA_HOLDER"] = val.VirtualAccountName
			}
		}

		if cifSchema.Username != "" {
			fmt.Println("Sending email..")
			go email.SendEmailVerificationSuccess(cifSchema.Username, cifSchema.Name, vaData["BCA"], vaData["BCA_HOLDER"], vaData["BRI"], vaData["BRI_HOLDER"])
			// sendgrid := email.Sendgrid{}
			// sendgrid.SetFrom("Amartha", "no-reply@amartha.com")
			// sendgrid.SetTo(cifSchema.Name, cifSchema.Username)
			// sendgrid.SetSubject(cifSchema.Name + ", Verifikasi Data Anda Berhasil")
			// sendgrid.VerifiedBodyEmail("VERIFIED_DATA", cifSchema.Name, cifSchema.Username, vaData)
			// sendgrid.SendEmail()
		}

		if cifSchema.PhoneNo != "" {
			// send sms notification
			fmt.Println("Sending sms ... ")
			twilio := services.InitTwilio()
			message := "Selamat data Anda sudah terverifikasi. Silakan login ke dashboard Anda dan mulai berinvestasi. www.amartha.com \n\nAmartha"
			twilio.SetParam(cifSchema.PhoneNo, message)
			twilio.SendSMS()
		}

		ctx.JSON(iris.StatusOK, iris.Map{
			"status":             "success",
			"verificationStatus": "verified",
		})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":             "success",
			"verificationStatus": "verification failed investor not validated",
		})
	}
}

func ParseDatatableURI(fullURI string) DataTable {
	var dtTables DataTable
	u, _ := url.Parse(fullURI)
	q := u.Query()
	for k, v := range q {
		if len(strings.TrimSpace(v[0])) == 0 {
			if err := json.Unmarshal([]byte(k), &dtTables); err == nil {
				return dtTables
			}
		}
	}
	return dtTables
}
