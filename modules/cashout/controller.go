package cashout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Cashout{})
	services.BaseCrudInit(Cashout{}, []Cashout{})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

func FetchDatatables(ctx *iris.Context) {
	var dtTables services.DataTable
	var totalRows int
	stage := ctx.URLParam("stage")
	limit := ctx.URLParam("limit")
	page := ctx.URLParam("page")
	submenu := ctx.URLParam("submenu")

	investorName := ctx.URLParam("investorName")
	dateSendToMandiri := ctx.URLParam("dateSTM")

	cashoutInvestors := []CashoutInvestor{}

	dtTables = services.ParseDatatableURI(string(ctx.URI().FullURI()))
	search := dtTables.Search.Value

	orderBy := "TRANSACTIONDATE"
	orderDir := "DESC"

	if len(dtTables.Columns) > 0 && len(dtTables.OrderInfo) > 0 {
		orderBy = dtTables.Columns[dtTables.OrderInfo[0].Column].Data
		orderDir = dtTables.OrderInfo[0].Dir
	}

	query := `
		SELECT r_cif_investor."cifId", r_cif_investor."investorId", r_account_transaction_credit_cashout."cashoutId", 
		r_account_investor."accountId", cif.name AS "investorName", cif."username" as "username", cashout."cashoutId" AS "cashoutNo", cashout.amount, 
		account."totalDebit", account."totalCredit", account."totalBalance", account_transaction_credit."type",  account_transaction_credit."transactionDate", 
		account_transaction_credit.remark, cashout.stage, cashout."sentToMandiriAt" as "sentToMandiriAt", count(*) OVER() AS full_count
		FROM cashout
		JOIN r_account_transaction_credit_cashout ON r_account_transaction_credit_cashout."cashoutId" = cashout.id
		JOIN r_account_transaction_credit ON r_account_transaction_credit."accountTransactionCreditId" = r_account_transaction_credit_cashout."accountTransactionCreditId"
		JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit_cashout."accountTransactionCreditId"
		JOIN account ON account.id = r_account_transaction_credit."accountId"
		JOIN r_account_investor ON r_account_investor."accountId" = account.id
		JOIN r_cif_investor ON r_cif_investor."investorId" = r_account_investor."investorId"
		JOIN cif ON cif.id = r_cif_investor."cifId" 
	`

	if stage == "" || stage == "ALL" {
		// query += "where stage not like 'SUCCESS' and stage not like 'CANCEL_AND_REFUND' "
		if submenu == "on-progress" {
			query += "where stage not like 'SUCCESS' and stage not like 'CANCEL%' "
		} else if submenu == "completed" && len(investorName) == 0 && len(dateSendToMandiri) == 0 {
			query += "where stage like 'SUCCESS' or stage like 'CANCEL%' "
		}

	} else if stage != "" && stage != "ALL" {
		// query += strings.Replace("WHERE stage ='?'", "?", stage, -1)
		query += ` WHERE stage LIKE '%` + stage + `%'`
	}

	if len(strings.TrimSpace(search)) > 0 {
		query += ` AND (cif.name LIKE '%` + search + `%' OR investor."investorNo"::text LIKE '%` + search + `%' OR cif."idCardNo" LIKE '%` + search + `%' OR  
					cif."taxCardNo" LIKE '%` + search + `%' OR cif."username" LIKE '%` + search + `%') `
	}

	/*
		if len(investorName) > 0 {
			query += ` AND cif.name LIKE '%` + investorName + `%'`
		}
	*/

	if len(investorName) > 0 {
		if submenu == "completed" {
			query += `where stage like 'SUCCESS' AND cif.name LIKE '%` + investorName + `%' or stage like 'CANCEL%' AND cif.name LIKE '%` + investorName + `%' `
		} else {
			query += ` AND cif.name LIKE '%` + investorName + `%'`
		}
	}

	if len(dateSendToMandiri) > 0 {
		query += ` AND cashout."sentToMandiriAt"::date = '` + dateSendToMandiri + `'`
	}

	/*
		if len(investorName) > 0 && len(stageId) > 0 && len(dateSendToMandiri) > 0 {
			query += ` AND (cif.name LIKE '%` + investorName + `%' AND cashout."stage" LIKE '%` +
				stageId + `%' AND account_transaction_credit."transactionDate" = '` + dateSendToMandiri + `') `
		}
	*/

	if len(strings.TrimSpace(orderBy)) > 0 && len(strings.TrimSpace(orderDir)) > 0 {
		switch strings.ToUpper(orderBy) {
		case "CASHOUTID":
			query += ` ORDER BY r_account_transaction_credit_cashout."cashoutId" ` + orderDir
		case "CASHOUTNO":
			query += ` ORDER BY cashout."cashoutId" ` + orderDir
		case "TRANSACTIONDATE":
			query += ` ORDER BY account_transaction_credit."transactionDate" ` + orderDir
		case "SENTTOMANDIRIAT":
			query += ` ORDER BY cashout."sentToMandiriAt" ` + orderDir
		case "INVESTORNAME":
			query += ` ORDER BY cif.name ` + orderDir
		case "TOTALDEBIT":
			query += ` ORDER BY account."totalDebit" ` + orderDir
		case "TOTALCREDIT":
			query += ` ORDER BY account."totalCredit" ` + orderDir
		case "CURRENTBALANCE":
			query += ` ORDER BY account."totalBalance" ` + orderDir
		case "CASHOUTAMOUNT":
			query += ` ORDER BY cashout.amount ` + orderDir
		case "STAGE":
			query += ` ORDER BY cashout.stage ` + orderDir
		default:
			query += ``
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

	services.DBCPsql.Raw(query).Find(&cashoutInvestors)

	if len(cashoutInvestors) > 0 {
		totalRows = cashoutInvestors[0].RowsFullCount
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalRows,
		"data":      cashoutInvestors,
	})
}

// FetchAll - fetch cashout data
func FetchAll(ctx *iris.Context) {
	var totalRows int
	investorName := ctx.URLParam("investorName")
	stage := ctx.URLParam("stage")
	dateSendToMandiri := ctx.URLParam("dateSTM")
	cashoutId := ctx.URLParam("cashoutId")
	submenu := ctx.URLParam("submenu")

	cashoutInvestors := []CashoutInvestor{}

	query := `
		SELECT r_cif_investor."cifId", r_cif_investor."investorId", r_account_transaction_credit_cashout."cashoutId", 
		r_account_investor."accountId", cif.name AS "investorName", cif."username" as "username", cashout."cashoutId" AS "cashoutNo", cashout.amount, 
		account."totalDebit", account."totalCredit", account."totalBalance", account_transaction_credit."type",  account_transaction_credit."transactionDate", 
		account_transaction_credit.remark, cashout.stage, cashout."sentToMandiriAt" as "sentToMandiriAt",  count(*) OVER() AS full_count
		FROM cashout
		JOIN r_account_transaction_credit_cashout ON r_account_transaction_credit_cashout."cashoutId" = cashout.id
		JOIN r_account_transaction_credit ON r_account_transaction_credit."accountTransactionCreditId" = r_account_transaction_credit_cashout."accountTransactionCreditId"
		JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit_cashout."accountTransactionCreditId"
		JOIN account ON account.id = r_account_transaction_credit."accountId"
		JOIN r_account_investor ON r_account_investor."accountId" = account.id
		JOIN r_cif_investor ON r_cif_investor."investorId" = r_account_investor."investorId"
		JOIN cif ON cif.id = r_cif_investor."cifId" 
	`

	if stage == "" || stage == "ALL" {
		// query += "where stage not like 'SUCCESS'"
		if submenu == "on-progress" {
			query += "where stage not like 'SUCCESS' and stage not like 'CANCEL%' "
		} else if submenu == "completed" && len(investorName) == 0 && len(dateSendToMandiri) == 0 {
			query += "where stage like 'SUCCESS' or stage like 'CANCEL%' "
		}

	} else if stage != "" && stage != "ALL" {
		query += ` WHERE stage LIKE '%` + stage + `%'`
		// query += strings.Replace("WHERE stage ='?'", "?", stage, -1)
	}

	/*
		if len(investorName) > 0 {
			query += ` AND cif.name LIKE '%` + investorName + `%'`
		}
	*/

	if len(investorName) > 0 {
		if submenu == "completed" {
			query += `where stage like 'SUCCESS' AND cif.name LIKE '%` + investorName + `%' or stage like 'CANCEL%' AND cif.name LIKE '%` + investorName + `%' `
		} else {
			query += ` AND cif.name LIKE '%` + investorName + `%'`
		}
	}

	if len(dateSendToMandiri) > 0 {
		query += ` AND cashout."sentToMandiriAt"::date = '` + dateSendToMandiri + `'`
	}

	if cashoutId != "" {
		query += ` AND cashout."id" = ` + cashoutId

	}

	if len(cashoutInvestors) > 0 {
		totalRows = cashoutInvestors[0].RowsFullCount
	}

	if err := services.DBCPsql.Raw(query).Find(&cashoutInvestors).Error; err != nil {
		log.Println(err)
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalRows,
		"data":      cashoutInvestors,
	})
}

// UpdateStage - update cashout stage
func UpdateStage(ctx *iris.Context) {
	cashoutInvestorSchema := CashoutInvestor{}

	if err := ctx.ReadJSON(&cashoutInvestorSchema); err != nil {
		ErrorLogger("ReadJSON Error:", err)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Failed to read param",
		})
		return
	}

	_, err := strconv.ParseUint(ctx.Param("cashout_id"), 10, 64)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Invalid cashout ID",
		})
		return
	}
	stage := ctx.Param("stage")

	params := map[string]string{"cashoutId": ctx.Param("cashout_id")}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(params)

	var url string = config.GoWithdrawalPath + "/api/v1/cashout/update/" + ctx.Param("cashout_id") + "/stage/" + stage

	req, err := http.NewRequest("PUT", url, b)
	// req.Header.Set("X-Auth-Token", "AMARTHA123")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    iris.Map{},
		})
		return
	}
	defer resp.Body.Close()

	var body struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Success bool        `json:"success"`
	}

	json.NewDecoder(resp.Body).Decode(&body)
	fmt.Println(body)

	if !body.Success {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":  "error",
			"message": body.Message,
		})
		return
	}

	// failed send request to mandiri
	ss := strings.Split(body.Message, " ")
	if ss[len(ss)-1] == "FAILED_SEND_REQUEST_TO_MANDIRI" {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":  "error",
			"message": "Failed send request to mandiri. " + body.Message,
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   body.Message,
	})
	return

}

// emergency error logger
func ErrorLogger(desc string, logError error) {
	startTime := time.Now()
	// logfilename
	filename := "./Go-Banking-Error.log"
	// write log.
	t := fmt.Sprintf("Process: %s. Time: %v, Error: %v.\n", desc, startTime, logError)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	_, err = f.WriteString(t)
	if err != nil {
		fmt.Println(err)
	}
}
