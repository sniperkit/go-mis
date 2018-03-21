package cashout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/account"
	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	cashoutHistory "bitbucket.org/go-mis/modules/cashout-history"
	"bitbucket.org/go-mis/modules/r"
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
		ErrorLogger("ReadJSON Error:", err)
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
		fmt.Println("send to mandiri")
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
			ErrorLogger("Call Go-Banking. HTTP REQUEST ERROR.", err)
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": err.Error(),
				"data":    iris.Map{},
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
			ErrorLogger("Return From Go-Banking. Failed to process request", errors.New("Body.Success == false"))
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": "Failed to process request.",
				"data":    iris.Map{},
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

		ErrorLogger("ERROR", errors.New("Bad Request"))
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "Bad request.",
		})
		return
	}
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
