package loan

import (
	"fmt"
	"strconv"
	"time"

	"errors"

	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	"bitbucket.org/go-mis/modules/installment"
	loanHistory "bitbucket.org/go-mis/modules/loan-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Loan{})
	services.BaseCrudInit(Loan{}, []Loan{})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchAll - fetchAll Loan data
func FetchAll(ctx *iris.Context) {
	branchID := ctx.Get("BRANCH_ID")

	totalData := TotalData{}
	queryTotalData := "SELECT DISTINCT COUNT(loan.*) AS \"totalRows\" "
	queryTotalData += "FROM loan "
	queryTotalData += "LEFT JOIN r_loan_sector ON r_loan_sector.\"loanId\" = loan.\"id\" "
	queryTotalData += "LEFT JOIN sector ON r_loan_sector.\"sectorId\" = sector.\"id\" "
	queryTotalData += "LEFT JOIN r_loan_borrower ON r_loan_borrower.\"borrowerId\" = loan.\"id\" "
	queryTotalData += "LEFT JOIN borrower ON r_loan_borrower.\"borrowerId\" = borrower.\"id\" "
	queryTotalData += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\""
	queryTotalData += "LEFT JOIN cif ON r_cif_borrower.\"cifId\" = cif.\"id\" LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	queryTotalData += "LEFT JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "
	queryTotalData += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	queryTotalData += "LEFT JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	queryTotalData += "LEFT JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.\"id\" "
	queryTotalData += "LEFT JOIN disbursement ON disbursement.\"id\" = r_loan_disbursement.\"disbursementId\" "
	queryTotalData += "WHERE branch.id = ? AND loan.\"deletedAt\" IS NULL AND loan.\"stage\" NOT IN ('END', 'END-EARLY') "

	if ctx.URLParam("search") != "" {
		queryTotalData += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	services.DBCPsql.Raw(queryTotalData, branchID).Find(&totalData)

	loans := []LoanFetch{}

	var limitPagination int64 = 10
	var offset int64 = 0

	query := "SELECT DISTINCT loan.*, "
	query += "sector.\"name\" as \"sector\", "
	query += "cif.\"name\" as \"borrower\", "
	query += "\"group\".\"name\" as \"group\", "
	query += "branch.\"name\" as \"branch\",  "
	query += "disbursement.\"disbursementDate\" "
	query += "FROM loan "
	query += "LEFT JOIN r_loan_sector ON r_loan_sector.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN sector ON r_loan_sector.\"sectorId\" = sector.\"id\" "
	query += "LEFT JOIN r_loan_borrower ON r_loan_borrower.\"borrowerId\" = loan.\"id\" "
	query += "LEFT JOIN borrower ON r_loan_borrower.\"borrowerId\" = borrower.\"id\" "
	query += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\" "
	query += "LEFT JOIN cif ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	query += "LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	query += "LEFT JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN disbursement ON disbursement.\"id\" = r_loan_disbursement.\"disbursementId\" "
	query += "WHERE branch.id = ? AND loan.\"deletedAt\" IS NULL AND loan.\"stage\" NOT IN ('END', 'END-EARLY') "

	if ctx.URLParam("search") != "" {
		query += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	if ctx.URLParam("limit") != "" {
		query += "LIMIT " + ctx.URLParam("limit") + " "
	} else {
		query += "LIMIT " + strconv.FormatInt(limitPagination, 10) + " "
	}

	if ctx.URLParam("page") != "" {
		offset, _ = strconv.ParseInt(ctx.URLParam("page"), 10, 64)
		query += "OFFSET " + strconv.FormatInt(offset, 10)
	} else {
		query += "OFFSET 0"
	}

	services.DBCPsql.Raw(query, branchID).Find(&loans)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      loans,
	})
}

// FetchDropping - Fetch loans which in ARCHIVE or DISBURSEMENT-FAILED stage
func FetchDropping(ctx *iris.Context) {
	loanData := []LoanDropping{}

	// ref: dropping.sql
	query := "SELECT loan.id, stage, cif_borrower.\"name\" AS borrower, \"group\".\"name\" AS \"group\", cif_investor.name AS investor "
	query += "FROM loan "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.id "
	query += "JOIN borrower ON borrower.id = r_loan_borrower.\"borrowerId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.id "
	query += "JOIN (SELECT * FROM cif WHERE \"deletedAt\" IS NULL) AS cif_borrower ON cif_borrower.id = r_cif_borrower.\"cifId\" "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.id "
	query += "JOIN \"group\" ON \"group\".id = r_loan_group.\"groupId\" "
	query += "JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"loanId\"= loan.id "
	query += "JOIN investor ON investor.id = r_investor_product_pricing_loan.\"investorId\" "
	query += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	query += "JOIN (SELECT * FROM cif WHERE \"deletedAt\" IS NULL) AS cif_investor ON cif_investor.id = r_cif_investor.\"cifId\" "
	query += "WHERE loan.\"deletedAt\" IS NULL AND (loan.stage = 'ARCHIVE' OR loan.stage = 'DISBURSEMENT-FAILED')"

	services.DBCPsql.Raw(query).Find(&loanData)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loanData,
	})
}

// UpdateStage - Update Stage Loan
func UpdateStage(ctx *iris.Context) {
	// Habib: logic dipindah ke executeUpdateStage supaya bisa di-reuse
	// loanData := Loan{}

	// loanID := ctx.Param("id")
	// stage := ctx.Param("stage")
	// services.DBCPsql.First(&loanData, loanID)
	// if loanData == (Loan{}) {
	// 	ctx.JSON(iris.StatusInternalServerError, iris.Map{
	// 		"status":  "error",
	// 		"message": "Can't find any loan detail.",
	// 	})
	// 	return
	// }

	// loanHistoryData := loanHistory.LoanHistory{StageFrom: loanData.Stage, StageTo: stage, Remark: "loanId=" + fmt.Sprintf("%v", loanData.ID) + " updated stage to " + stage}
	// services.DBCPsql.Table("loan_history").Create(&loanHistoryData)

	// rLoanHistory := r.RLoanHistory{LoanId: loanData.ID, LoanHistoryId: loanHistoryData.ID}
	// services.DBCPsql.Table("r_loan_history").Create(&rLoanHistory)

	// services.DBCPsql.Table("loan").Where("\"id\" = ?", loanData.ID).UpdateColumn("stage", stage)

	// ctx.JSON(iris.StatusOK, iris.Map{
	// 	"status":    "success",
	// 	"stageFrom": loanData.Stage,
	// 	"stageTo":   stage,
	// })

	loanID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	stage := ctx.Param("stage")

	loanStage, err := executeUpdateStage(loanID, stage)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Can't find any loan detail.",
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"stageFrom": loanStage,
		"stageTo":   stage,
	})
	return
}

func executeUpdateStage(id uint64, stage string) (string, error) {
	loanData := Loan{}
	services.DBCPsql.Where("id = ?", id).First(&loanData)
	fmt.Printf("%v", loanData)
	if loanData == (Loan{}) {
		return "", errors.New("loan not found")
	}

	loanHistoryData := loanHistory.LoanHistory{StageFrom: loanData.Stage, StageTo: stage, Remark: "loanId=" + fmt.Sprintf("%v", loanData.ID) + " updated stage to " + stage}
	services.DBCPsql.Table("loan_history").Create(&loanHistoryData)

	rLoanHistory := r.RLoanHistory{LoanId: loanData.ID, LoanHistoryId: loanHistoryData.ID}
	services.DBCPsql.Table("r_loan_history").Create(&rLoanHistory)

	services.DBCPsql.Table("loan").Where("\"id\" = ?", loanData.ID).UpdateColumn("stage", stage)

	return loanData.Stage, nil
}

func GetLoanDetail(ctx *iris.Context) {
	loanObj := Loan{}
	borrowerObj := LoanBorrowerProfile{}
	installmentObj := []installment.Installment{}

	loanId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		fmt.Println(err)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Something went wrong. Please try again later.",
		})
		return
	}

	services.DBCPsql.Table("loan").Where("\"id\" = ?", loanId).First(&loanObj)

	queryBorrowerObj := "SELECT cif.\"cifNumber\", cif.\"name\", \"group\".\"name\" AS \"group\", area.\"name\" AS \"area\", branch.\"name\" AS \"branch\" "
	queryBorrowerObj += "FROM loan "
	queryBorrowerObj += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.\"id\" "
	queryBorrowerObj += "JOIN borrower ON borrower.\"id\" = r_loan_borrower.\"borrowerId\" "
	queryBorrowerObj += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\" "
	queryBorrowerObj += "JOIN cif ON cif.\"id\" = r_cif_borrower.\"cifId\" "
	queryBorrowerObj += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	queryBorrowerObj += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
	queryBorrowerObj += "JOIN r_loan_area ON r_loan_area.\"loanId\" = loan.\"id\" "
	queryBorrowerObj += "JOIN area ON area.\"id\" = r_loan_area.\"areaId\" "
	queryBorrowerObj += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	queryBorrowerObj += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\" "
	queryBorrowerObj += "WHERE loan.\"id\" = ?"

	services.DBCPsql.Raw(queryBorrowerObj, loanId).Find(&borrowerObj)

	queryInstallmentObj := "SELECT * "
	queryInstallmentObj += "FROM installment "
	queryInstallmentObj += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	queryInstallmentObj += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
	queryInstallmentObj += "WHERE loan.\"id\" = ?"

	services.DBCPsql.Raw(queryInstallmentObj, loanId).Find(&installmentObj)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"loan":        loanObj,
			"borrower":    borrowerObj,
			"installment": installmentObj,
		},
	})
}

// RefundAndChangeStageTo - refund investor balance and change loan stage
func RefundAndChangeStageTo(ctx *iris.Context) {
	loanID, _ := strconv.ParseUint(ctx.Param("loan_id"), 10, 64)
	stage := ctx.Param("stage")

	loanStage, err := executeUpdateStage(loanID, stage)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Can't find any loan detail.",
		})
		return
	}

	// get loan_id, investor_id, account_id, plafond
	refundBase := RefundBase{}
	// ref: refund-base.sql
	queryRefundBase := "SELECT loan.id AS loan_id, investor.id AS investor_id, account.id AS account_id, loan.plafond "
	queryRefundBase += "FROM loan "
	queryRefundBase += "JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"loanId\" = loan.id "
	queryRefundBase += "JOIN investor ON investor.id = r_investor_product_pricing_loan.\"investorId\" "
	queryRefundBase += "JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.id "
	queryRefundBase += "JOIN account ON account.id = r_account_investor.\"accountId\" "
	queryRefundBase += "WHERE loan.id = ? "
	services.DBCPsql.Raw(queryRefundBase, loanID).First(&refundBase)

	// add new account_transaction_debit entry
	transaction := accountTransactionDebit.AccountTransactionDebit{
		Type:            "REFUND",
		TransactionDate: time.Now(),
		Amount:          refundBase.Plafond,
		Remark:          "",
	}
	services.DBCPsql.Table("account_transaction_debit").Create(&transaction)

	// connect the entry to investor account
	rTransaction := r.RAccountTransactionDebit{
		AccountId:                 refundBase.AccountID,
		AccountTransactionDebitId: transaction.ID,
	}
	services.DBCPsql.Table("r_account_transaction_debit").Create(&rTransaction)

	// calculate account balance and save it to account

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":      "success",
		"stageFrom":   loanStage,
		"stageTo":     stage,
		"refundBase":  refundBase,
		"transaction": transaction,
	})
}
