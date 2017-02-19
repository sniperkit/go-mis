package installment

import (
	"strconv"
	"strings"
	"time"

	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	installmentHistory "bitbucket.org/go-mis/modules/installment-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Installment{})
	services.BaseCrudInit(Installment{}, []Installment{})
}

// FetchAll - fetchAll installment data
func FetchAll(ctx *iris.Context) {
	installments := []InstallmentFetch{}

	query := "SELECT branch.\"name\" AS \"branch\", \"group\".\"id\" AS \"groupId\", \"group\".\"name\" AS \"group\", SUM(installment.\"paidInstallment\") AS \"totalPaidInstallment\", installment.\"createdAt\"::date "
	query += "FROM installment "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	query += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\"  "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
	query += "GROUP BY installment.\"createdAt\"::date, branch.\"name\", \"group\".\"id\", \"group\".\"name\" "
	query += "ORDER BY installment.\"createdAt\"::date DESC, branch.\"name\" ASC LIMIT 5"

	services.DBCPsql.Raw(query).Find(&installments)
	ctx.JSON(iris.StatusOK, iris.Map{"data": installments})
}

// SubmitInstallment - submit installment data
func SubmitInstallment(ctx *iris.Context) {
	installment := Installment{}
	tempLoanID := ctx.Param("loan_id")

	loanID, err := strconv.ParseUint(tempLoanID, 10, 64)

	if err != nil {
		ctx.JSON(iris.StatusExpectationFailed, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = ctx.ReadJSON(&installment)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	services.DBCPsql.Table("installment").Create(&installment)

	installmentHistoryData := installmentHistory.InstallmentHistory{StageFrom: installment.Stage, StageTo: installment.Stage}
	services.DBCPsql.Table("r_installment_history").Create(&installmentHistoryData)

	rInstallmentHistoryData := r.RInstallmentHistory{InstallmentId: installment.ID, InstallmentHistoryId: installmentHistoryData.ID}
	services.DBCPsql.Table("r_installment_history").Create(&rInstallmentHistoryData)

	accountTransactionDebitData := accountTransactionDebit.AccountTransactionDebit{Type: "INSTALLMENT", Amount: installment.PaidInstallment}
	services.DBCPsql.Table("account_transaction_debit").Create(&accountTransactionDebitData)

	rLoanInvestorProductPricing := r.RInvestorProductPricingLoan{}
	services.DBCPsql.Table("r_investor_product_pricing_loan").Where("\"loanId\" = ?", loanID).First(&rLoanInvestorProductPricing)

	rAccountInvestor := r.RAccountInvestor{}
	services.DBCPsql.Table("r_account_investor").Where("\"investorId\" = ?", rLoanInvestorProductPricing.InvestorId).First(&rAccountInvestor)

	rAccountTransactionDebitData := r.RAccountTransactionDebit{AccountId: rAccountInvestor.AccountId, AccountTransactionDebitId: accountTransactionDebitData.ID}
	go services.DBCPsql.Table("r_account_transaction_debit").Create(&rAccountTransactionDebitData)

	loanInstallmentData := r.RLoanInstallment{LoanId: loanID, InstallmentId: installment.ID}
	go services.DBCPsql.Table("r_loan_installment").Create(&loanInstallmentData)

	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status": "success",
		"data":   installment,
	})

}

//GetInstallmentByGroupIDAndTransactionDate - get list of installment by group and transaction date
func GetInstallmentByGroupIDAndTransactionDate(ctx *iris.Context) {
	groupID := ctx.Param("group_id")
	transactionDate := ctx.Param("transaction_date")

	query := "SELECT "
	query += "\"group\".\"id\" as \"groupId\", \"group\".\"name\" as \"groupName\","
	query += "installment.\"id\" as \"installmentId\", installment.\"type\", installment.\"paidInstallment\", installment.\"penalty\", installment.\"reserve\", installment.\"presence\", installment.\"frequency\", installment.\"stage\" "
	query += "FROM installment "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	query += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\"  "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
	query += "WHERE installment.\"createdAt\"::date = ? AND \"group\".\"id\" = ?"

	installmentDetailSchema := []InstallmentDetail{}
	services.DBCPsql.Raw(query, transactionDate, groupID).Scan(&installmentDetailSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   installmentDetailSchema,
	})
}

type LoanInvestorAccountID struct {
	LoanID     uint64 `gorm:"column:loanId" json:"loanId"`
	InvestorID uint64 `gorm:"column:investorId" json:"investorId"`
	AccountID  uint64 `gorm:"column:accountId" json:"accountId"`
}

type AccountTransactionDebitAndCredit struct {
	TotalDebit  float64 `gorm:"column:totalDebit" json:"totalDebit"`
	TotalCredit float64 `gorm:"column:totalCredit" json:"totalCredit"`
}

func storeInstallment(installmentId uint64, status string) {
	// fmt.Println("[INFO] Storing installment. installmentID=" + strconv.FormatUint(installmentId, 64) + " status=" + status)
	installmentSchema := Installment{}
	services.DBCPsql.Table("installment").Where("\"id\" = ?", installmentId).First(&installmentSchema)

	if installmentSchema.Stage != "PENDING" {
		// ctx.JSON(iris.StatusBadRequest, iris.Map{
		// 	"status":  "error",
		// 	"message": "Current installment stage is NOT 'PENDING'. System cannot continue to process your request.",
		// })
		return
	}

	installmentHistorySchema := &installmentHistory.InstallmentHistory{StageFrom: "PENDING", StageTo: status}
	services.DBCPsql.Table("installment_history").Create(installmentHistorySchema)

	installmentHistoryID := installmentHistorySchema.ID

	rInstallmentHistorySchema := &r.RInstallmentHistory{InstallmentId: installmentId, InstallmentHistoryId: installmentHistoryID}
	services.DBCPsql.Table("r_installment_history").Create(rInstallmentHistorySchema)

	services.DBCPsql.Table("installment").Where("\"id\" = ?", installmentId).UpdateColumn("stage", status)

	if status == "REJECT" {
		// ctx.JSON(iris.StatusOK, iris.Map{
		// 	"status": "success",
		// 	"data":   iris.Map{"message": "Installment data has been rejected."},
		// })
		return
	}

	accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "INSTALLMENT", TransactionDate: time.Now(), Amount: installmentSchema.PaidInstallment}
	services.DBCPsql.Table("account_transaction_debit").Create(accountTransactionDebitSchema)

	queryGetAccountInvestor := "SELECT r_loan_installment.\"loanId\", r_investor_product_pricing_loan.\"investorId\", r_account_investor.\"accountId\" "
	queryGetAccountInvestor += "FROM installment "
	queryGetAccountInvestor += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	queryGetAccountInvestor += "JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"loanId\" = r_loan_installment.\"loanId\" "
	queryGetAccountInvestor += "JOIN r_account_investor ON r_account_investor.\"investorId\" = r_investor_product_pricing_loan.\"investorId\" "
	queryGetAccountInvestor += "WHERE installment.\"id\" = ?"

	loanInvestorAccountIDSchema := LoanInvestorAccountID{}
	services.DBCPsql.Raw(queryGetAccountInvestor, installmentId).Scan(&loanInvestorAccountIDSchema)

	rAccountTransactionDebit := &r.RAccountTransactionDebit{AccountId: loanInvestorAccountIDSchema.AccountID, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
	services.DBCPsql.Table("r_account_transaction_debit").Create(rAccountTransactionDebit)

	querySumDebitAndCredit := "SELECT SUM(account_transaction_debit.\"amount\") as \"totalDebit\", SUM(account_transaction_credit.\"amount\")  as \"totalCredit\" "
	querySumDebitAndCredit += "FROM account "
	querySumDebitAndCredit += "JOIN r_account_transaction_debit ON r_account_transaction_debit.\"accountId\" = account.\"id\" "
	querySumDebitAndCredit += "JOIN account_transaction_debit ON account_transaction_debit.\"id\" = r_account_transaction_debit.\"accountTransactionDebitId\" "
	querySumDebitAndCredit += "JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountId\" = account.\"id\" "
	querySumDebitAndCredit += "JOIN account_transaction_credit ON account_transaction_credit.\"id\" = r_account_transaction_credit.\"accountTransactionCreditId\" "
	querySumDebitAndCredit += "WHERE account.\"id\" = ?"

	accountTransactionDebitAndCreditSchema := AccountTransactionDebitAndCredit{}
	services.DBCPsql.Raw(querySumDebitAndCredit, loanInvestorAccountIDSchema.AccountID).Scan(&accountTransactionDebitAndCreditSchema)

	totalBalance := accountTransactionDebitAndCreditSchema.TotalDebit - accountTransactionDebitAndCreditSchema.TotalCredit
	services.DBCPsql.Table("account").Exec("UPDATE account SET \"totalDebit\" = ?, \"totalCredit\" = ?, \"totalBalance\" = ? WHERE \"id\" = ?", accountTransactionDebitAndCreditSchema.TotalDebit, accountTransactionDebitAndCreditSchema.TotalCredit, totalBalance, loanInvestorAccountIDSchema.AccountID)

	// ctx.JSON(iris.StatusOK, iris.Map{
	// 	"status": "success",
	// 	"data": iris.Map{
	// 		"message": "Installment has been updated to " + status,
	// 	},
	// })
}

// SubmitInstallmentByInstallmentIDWithStatus - approve or reject installment by installment_id
func SubmitInstallmentByInstallmentIDWithStatus(ctx *iris.Context) {
	installmentID, _ := strconv.ParseUint(ctx.Param("installment_id"), 10, 64)
	status := strings.ToUpper(ctx.Param("status"))

	go storeInstallment(installmentID, status)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"message": "Your request has been received. It might need take a while to process your request.",
		},
	})
}

// SubmitInstallmentByGroupIDAndTransactionDateWithStatus - approve or reject installment per group
func SubmitInstallmentByGroupIDAndTransactionDateWithStatus(ctx *iris.Context) {
	groupID := ctx.Param("group_id")
	transactionDate := ctx.Param("transaction_date")
	status := strings.ToUpper(ctx.Param("status"))

	if strings.ToLower(ctx.Param("status")) == "approve" || strings.ToLower(ctx.Param("status")) == "reject" {
		query := "SELECT "
		query += "\"group\".\"id\" as \"groupId\", \"group\".\"name\" as \"groupName\","
		query += "installment.\"id\" as \"installmentId\", installment.\"type\", installment.\"paidInstallment\", installment.\"penalty\", installment.\"reserve\", installment.\"presence\", installment.\"frequency\", installment.\"stage\" "
		query += "FROM installment "
		query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
		query += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
		query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
		query += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\"  "
		query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
		query += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
		query += "WHERE installment.\"createdAt\"::date = ? AND \"group\".\"id\" = ? AND installment.\"stage\" != 'APPROVE'"

		installmentDetailSchema := []InstallmentDetail{}
		services.DBCPsql.Raw(query, transactionDate, groupID).Scan(&installmentDetailSchema)

		for _, item := range installmentDetailSchema {
			go storeInstallment(item.InstallmentID, status)
		}

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data": iris.Map{
				"message": "Your request has been received. It might need take a while to process your request.",
			},
		})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data": iris.Map{
				"message": "Invalid status.",
			},
		})
	}
}
