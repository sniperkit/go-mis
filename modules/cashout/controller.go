package cashout

import (
	"strconv"
	"time"
	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/account"
	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	cashoutHistory "bitbucket.org/go-mis/modules/cashout-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
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

	investorName := ctx.URLParam("investorName")
	stageId := ctx.URLParam("stage")
	dateSendToMandiri := ctx.URLParam("dateSTM")

	cashoutInvestors := []CashoutInvestor{}

	dtTables = services.ParseDatatableURI(string(ctx.URI().FullURI()))
	search := dtTables.Search.Value
	orderBy := "NAME"
	orderDir := "ASC"

	if len(dtTables.Columns) > 0 && len(dtTables.OrderInfo) > 0 {
		orderBy = dtTables.Columns[dtTables.OrderInfo[0].Column].Data
		orderDir = dtTables.OrderInfo[0].Dir
	}

	query := `
		SELECT r_cif_investor."cifId", r_cif_investor."investorId", r_account_transaction_credit_cashout."cashoutId", 
		r_account_investor."accountId", cif.name AS "investorName", cashout."cashoutId" AS "cashoutNo", cashout.amount, 
		account."totalDebit", account."totalCredit", account."totalBalance", account_transaction_credit."type",  account_transaction_credit."transactionDate", 
		account_transaction_credit.remark, cashout.stage, count(*) OVER() AS full_count
		FROM cashout
		JOIN r_account_transaction_credit_cashout ON r_account_transaction_credit_cashout."cashoutId" = cashout.id
		JOIN r_account_transaction_credit ON r_account_transaction_credit."accountTransactionCreditId" = r_account_transaction_credit_cashout."accountTransactionCreditId"
		JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit_cashout."accountTransactionCreditId"
		JOIN account ON account.id = r_account_transaction_credit."accountId"
		JOIN r_account_investor ON r_account_investor."accountId" = account.id
		JOIN r_cif_investor ON r_cif_investor."investorId" = r_account_investor."investorId"
		JOIN cif ON cif.id = r_cif_investor."cifId" 
	`

	if stage == "" {
		query += "where stage not like 'SUCCESS'"	
	} else if stage != "" && stage != "ALL" {
		query += strings.Replace("WHERE stage ='?'","?", stage, -1)
	}  

	if len(strings.TrimSpace(search)) > 0 {
		query += ` AND (cif.name ~* '` + search + `' OR investor."investorNo"::text ~* '` + search + `' OR cif."idCardNo" ~* '` + search + `' OR  
					cif."taxCardNo" ~* '` + search + `' OR cif."username"  ~* '` + search + `') `
	}

	if len(investorName) > 0 && len(stageId) > 0 && len(dateSendToMandiri) > 0 {
		query += ` AND (cif.name ~* '` + investorName + `' AND cashout."stage" ~* '` + 
			stageId + `' AND account_transaction_credit."transactionDate" ~* '` + dateSendToMandiri + `') `
	}

	if len(strings.TrimSpace(orderBy)) > 0 && len(strings.TrimSpace(orderDir)) > 0 {
		switch strings.ToUpper(orderBy) {
		case "CASHOUTID":
			query += ` ORDER BY r_account_transaction_credit_cashout."cashoutId" ` + orderDir
		case "CASHOUTNO":
			query += ` ORDER BY cashout."cashoutId" ` + orderDir
		case "REQUESTDATE":
			query += ` ORDER BY account_transaction_credit."transactionDate" ` + orderDir
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
	cashoutInvestors := []CashoutInvestor{}
	stage := ctx.URLParam("stage")

	// query := "SELECT cashout.*, "
	// query += "cif.\"name\" as \"investorName\", "
	// query += "account.\"totalBalance\" as \"totalBalance\" "
	// // query += "r_cif_borrower.\"borrowerId\" IS NOT NULL as \"isBorrower\", r_cif_investor.\"investorId\" IS NOT NULL as \"isInvestor\" "
	// query += "FROM cashout "
	// query += "LEFT JOIN r_investor_cashout ON r_investor_cashout.\"cashoutId\" = cashout.\"id\" "
	// query += "LEFT JOIN investor ON investor.\"id\" = r_investor_cashout.\"investorId\" "
	// query += "LEFT JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.\"id\" "
	// query += "LEFT JOIN cif ON r_cif_investor.\"cifId\" = cif.\"id\" "
	// query += "LEFT JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.\"id\" "
	// query += "LEFT JOIN account ON r_account_investor.\"accountId\" = account.\"id\" "

	query := "SELECT r_cif_investor.\"cifId\", r_cif_investor.\"investorId\", r_account_transaction_credit_cashout.\"cashoutId\", r_account_investor.\"accountId\", cif.name AS \"investorName\", cashout.\"cashoutId\" AS \"cashoutNo\", cashout.amount, account.\"totalDebit\", account.\"totalCredit\", account.\"totalBalance\", account_transaction_credit.\"type\",  account_transaction_credit.\"transactionDate\", account_transaction_credit.remark, cashout.stage "
	query += "FROM cashout "
	query += "JOIN r_account_transaction_credit_cashout ON r_account_transaction_credit_cashout.\"cashoutId\" = cashout.id "
	query += "JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountTransactionCreditId\" = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	query += "JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	query += "JOIN account ON account.id = r_account_transaction_credit.\"accountId\" "
	query += "JOIN r_account_investor ON r_account_investor.\"accountId\" = account.id "
	query += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = r_account_investor.\"investorId\" "
	query += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "

	queryTotalData := "SELECT count(*) AS \"totalRows\" "
	queryTotalData += "FROM cashout "
	queryTotalData += "JOIN r_account_transaction_credit_cashout ON r_account_transaction_credit_cashout.\"cashoutId\" = cashout.id "
	queryTotalData += "JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountTransactionCreditId\" = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	queryTotalData += "JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	queryTotalData += "JOIN account ON account.id = r_account_transaction_credit.\"accountId\" "
	queryTotalData += "JOIN r_account_investor ON r_account_investor.\"accountId\" = account.id "
	queryTotalData += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = r_account_investor.\"investorId\" "
	queryTotalData += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	
	if stage != "" {
		query += strings.Replace("WHERE stage = '?'", "?", stage, -1)
		queryTotalData += strings.Replace("WHERE stage = '?'", "?", stage, -1)
	}
 
	// if len(strings.TrimSpace(stage)) == 0 {
	// 	services.DBCPsql.Raw(query).Find(&cashoutInvestors)
	// } else {
	// 	query += "WHERE cashout.\"stage\" = ? "
	// 	services.DBCPsql.Raw(query, stage).Find(&cashoutInvestors)
	// }

	totalData := TotalData{}

	services.DBCPsql.Raw(query).Find(&cashoutInvestors)
	services.DBCPsql.Raw(queryTotalData).Find(&totalData)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      cashoutInvestors,
	})
}

// UpdateStage - update cashout stage
func UpdateStage(ctx *iris.Context) {
	cashoutInvestorSchema := CashoutInvestor{}

	if err := ctx.ReadJSON(&cashoutInvestorSchema); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Failed to read param",
		})
		return
	}

	cashoutID, _ := strconv.ParseUint(ctx.Param("cashout_id"), 10, 64)
	stage := ctx.Param("stage")

	cashoutHistoryObj := &cashoutHistory.CashoutHistory{StageFrom: "PENDING", StageTo: stage}
	services.DBCPsql.Create(cashoutHistoryObj)

	rCashoutHistoryObj := &r.RCashoutHistory{CashoutId: cashoutID, CashoutHistoryId: cashoutHistoryObj.ID}
	services.DBCPsql.Create(rCashoutHistoryObj)

	services.DBCPsql.Table("cashout").Where("cashout.\"id\" = ?", cashoutID).UpdateColumn("stage", stage)

	cashoutNo := strconv.FormatUint(cashoutInvestorSchema.CashoutNo, 10)

	if stage == "SEND-TO-MANDIRI" {
		cashoutSchema := Cashout{}
		services.DBCPsql.Table("cashout").Where("cashout.\"id\" = ?", cashoutID).Scan(&cashoutSchema)

		params := map[string]string{"cashoutId": cashoutSchema.CashoutID}

		b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(params)

    var url string = config.GoBankingPath + `/mandiri/payment`
    req, err := http.NewRequest("POST", url, b)
    req.Header.Set("X-Auth-Token", "AMARTHA123")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      fmt.Println(err)
      ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status": "error",
				"message": "Failed to process request.",
				"data":   iris.Map{},
			})
      return
    }
    defer resp.Body.Close()

    var body struct {
    	Success bool `json:"success"`
    }
    fmt.Println(resp)

		json.NewDecoder(resp.Body).Decode(&body)

		fmt.Println("---------")
		fmt.Printf("%+v", body)

		if body.Success == false {
			cashoutHistoryObj := &cashoutHistory.CashoutHistory{StageFrom: "SEND-TO-MANDIRI", StageTo: "FAILED-PROCESS-MANDIRI"}
			services.DBCPsql.Create(cashoutHistoryObj)

			rCashoutHistoryObj := &r.RCashoutHistory{CashoutId: cashoutID, CashoutHistoryId: cashoutHistoryObj.ID}
			services.DBCPsql.Create(rCashoutHistoryObj)

			services.DBCPsql.Table("cashout").Where("cashout.\"id\" = ?", cashoutID).UpdateColumn("stage", "FAILED-PROCESS-MANDIRI")
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status": "error",
				"message": "Failed to process request.",
				"data":   iris.Map{},
			})
			return
		}

		services.DBCPsql.Table("account_transaction_credit").Where("remark = ?", cashoutInvestorSchema.Remark).Update("remark", "CASHOUT #"+cashoutNo+" HAS BEEN SUBMITTED.")

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	} else if stage == "SUCCESS" {
		services.DBCPsql.Table("account_transaction_credit").Where("remark = ?", cashoutInvestorSchema.Remark).Update("remark", "CASHOUT #"+cashoutNo+" SUCCESS")
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	} else if stage == "CANCEL" {
		services.DBCPsql.Table("account_transaction_credit").Where("remark = ?", cashoutInvestorSchema.Remark).Update("remark", "CASHOUT #"+cashoutNo+" CANCEL")

		accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "REFUND", Remark: "REFUND CASHOUT #" + cashoutNo, Amount: cashoutInvestorSchema.Amount, TransactionDate: time.Now()}
		services.DBCPsql.Table("account_transaction_debit").Create(accountTransactionDebitSchema)

		rAccountTransactionDebitSchema := &r.RAccountTransactionDebit{AccountId: cashoutInvestorSchema.AccountID, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
		services.DBCPsql.Table("r_account_transaction_debit").Create(rAccountTransactionDebitSchema)

		totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(cashoutInvestorSchema.AccountID)
		totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(cashoutInvestorSchema.AccountID)

		totalBalance := totalDebit - totalCredit

		services.DBCPsql.Table("account").Where("id = ?", cashoutInvestorSchema.AccountID).Updates(account.Account{TotalDebit: totalDebit, TotalCredit: totalCredit, TotalBalance: totalBalance})

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	} else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "Bad request.",
		})
		return
	}
}

