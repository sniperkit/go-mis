package installment

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/go-mis/modules/account"
	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
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
// Habib : logicnya sudah bisa di handle sama FetchByType dgn parameter pending
func FetchAll(ctx *iris.Context) {
	branchID := ctx.Get("BRANCH_ID")
	installments := []InstallmentFetch{}

	query := "SELECT branch.\"name\" AS \"branch\", \"group\".\"id\" AS \"groupId\", \"group\".\"name\" AS \"group\", SUM(installment.\"paidInstallment\") AS \"totalPaidInstallment\", installment.\"createdAt\"::date "
	query += "FROM installment "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	query += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\"  "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
	query += "WHERE installment.stage = 'PENDING' AND branch.id = ?"
	query += "GROUP BY installment.\"createdAt\"::date, branch.\"name\", \"group\".\"id\", \"group\".\"name\" "
	query += "ORDER BY installment.\"createdAt\"::date DESC, branch.\"name\" ASC"

	services.DBCPsql.Raw(query, branchID).Find(&installments)
	ctx.JSON(iris.StatusOK, iris.Map{"data": installments})
}

// FetchByType - fetch installment data by type ["PENDING", "IN-REVIEW"]
func FetchByType(ctx *iris.Context) {
	installmentType := strings.ToUpper(ctx.Param("type"))
	branchID := ctx.Get("BRANCH_ID")
	installments := []InstallmentFetch{}

	query := "SELECT branch.\"name\" AS \"branch\", \"group\".\"id\" AS \"groupId\", \"group\".\"name\" AS \"group\", SUM(installment.\"paidInstallment\") AS \"totalPaidInstallment\", installment.\"createdAt\"::date "
	query += "FROM installment "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	query += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\"  "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
	query += "WHERE installment.stage = ? AND branch.id = ?"
	query += "GROUP BY installment.\"createdAt\"::date, branch.\"name\", \"group\".\"id\", \"group\".\"name\" "
	query += "ORDER BY installment.\"createdAt\"::date DESC, branch.\"name\" ASC"

	services.DBCPsql.Raw(query, installmentType, branchID).Find(&installments)
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
	branchID := ctx.Get("BRANCH_ID")
	groupID := ctx.Param("group_id")
	transactionDate := ctx.Param("transaction_date")

	query := `
		SELECT 
		"group".id as "groupId", 
		"group".name as "groupName",
		cif.name as "cifName",
		borrower."borrowerNo",
		loan.id "loanId",
		installment."id" as "installmentId", 
		installment.type, 
		installment."paidInstallment", 
		installment.penalty, 
		installment.reserve, 
		installment.presence, 
		installment.frequency, 
		installment.stage 

		FROM installment 

		JOIN r_loan_installment ON installment.id = r_loan_installment."installmentId"
		JOIN loan               ON loan.id        = r_loan_installment."loanId"
		JOIN r_loan_branch      ON loan.id        = r_loan_branch."loanId"
		JOIN r_loan_group       ON loan.id        = r_loan_group."loanId"
		JOIN "group"            ON "group".id     = r_loan_group."groupId"
		JOIN r_loan_borrower    ON loan.id        = r_loan_borrower."loanId"
		JOIN borrower           ON borrower.id    = r_loan_borrower."borrowerId"
		JOIN r_cif_borrower     ON borrower.id    = r_cif_borrower."borrowerId"
		JOIN cif                ON cif.id         = r_cif_borrower."cifId"

		WHERE 

		installment."createdAt"::date = ? 
		AND r_loan_group."groupId" = ? 
		AND r_loan_branch."branchId" = ?
	`

	installmentDetailSchema := []InstallmentDetail{}
	services.DBCPsql.Raw(query, transactionDate, groupID, branchID).Scan(&installmentDetailSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   installmentDetailSchema,
	})
}

type LoanInvestorAccountID struct {
	LoanID     uint64  `gorm:"column:loanId" json:"loanId"`
	InvestorID uint64  `gorm:"column:investorId" json:"investorId"`
	AccountID  uint64  `gorm:"column:accountId" json:"accountId"`
	PPLROI     float64 `gorm:"column:pplROI" json:"pplROI"`
}

type AccountTransactionDebitAndCredit struct {
	TotalDebit  float64 `gorm:"column:totalDebit" json:"totalDebit"`
	TotalCredit float64 `gorm:"column:totalCredit" json:"totalCredit"`
}

type LoanSchema struct {
	ID                   uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanPeriod           int64      `gorm:"column:loanPeriod" json:"loanPeriod"`
	AgreementType        string     `gorm:"column:agreementType" json:"agreementType"`
	Subgroup             string     `gorm:"column:subgroup" json:"subgrop"`
	Purpose              string     `gorm:"column:purpose" json:"purpose"`
	URLPic1              string     `gorm:"column:urlPic1" json:"urlPic1"`
	URLPic2              string     `gorm:"column:urlPic2" json:"urlPic2"`
	SubmittedLoanDate    string     `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	SubmittedPlafond     float64    `gorm:"column:submittedPlafond" json:"submittedPlafond"`
	SubmittedTenor       int64      `gorm:"column:submittedTenor" json:"submittedTenor"`
	SubmittedInstallment float64    `gorm:"column:submittedInstallment" json:"submittedInstallment"`
	CreditScoreGrade     string     `gorm:"column:creditScoreGrade" json:"creditScoreGrade"`
	CreditScoreValue     float64    `gorm:"column:creditScoreValue" json:"creditScoreValue"`
	Tenor                uint64     `gorm:"column:tenor" json:"tenor"`
	Rate                 float64    `gorm:"column:rate" json:"rate"`
	Installment          float64    `gorm:"column:installment" json:"installment"`
	Plafond              float64    `gorm:"column:plafond" json:"plafond"`
	GroupReserve         float64    `gorm:"column:groupReserve" json:"groupReserve"`
	Stage                string     `gorm:"column:stage" json:"stage"`
	IsLWK                bool       `gorm:"column:isLWK" json:"isLWK" sql:"default:false"`
	IsUPK                bool       `gorm:"column:isUPK" json:"IsUPK" sql:"default:false"`
	CreatedAt            time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt            *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

func storeInstallment(installmentId uint64, status string) {
	convertedInstallmentId := strconv.FormatUint(installmentId, 10)
	fmt.Println("[INFO] Storing installment. installmentID=" + convertedInstallmentId + " status=" + status)
	installmentSchema := Installment{}
	services.DBCPsql.Table("installment").Where("\"id\" = ?", installmentId).First(&installmentSchema)

	if installmentSchema.Stage != "PENDING" && installmentSchema.Stage != "IN-REVIEW" && installmentSchema.Stage != "APPROVE" {
		// ctx.JSON(iris.StatusBadRequest, iris.Map{
		// 	"status":  "error",
		// 	"message": "Current installment stage is NOT 'PENDING'. System cannot continue to process your request.",
		// })
		fmt.Println("Current installment stage is NEITHER 'PENDING' NOR 'IN-REVIEW'. System cannot continue to process your request. installmentId=" + convertedInstallmentId)
		return
	}

	if status == "REJECT" {
		// ctx.JSON(iris.StatusOK, iris.Map{
		// 	"status": "success",
		// 	"data":   iris.Map{"message": "Installment data has been rejected."},
		// })
		UpdateStageInstallmentApproveOrReject(installmentId, status)
		fmt.Println("Installment data has been rejected. installmentId=" + convertedInstallmentId)
		return
	}

	if status == "IN-REVIEW" {
		UpdateStageInstallmentApproveOrReject(installmentId, status)
		fmt.Println("Installment data will be reviewed. installmentId=" + convertedInstallmentId)
		return
	}

	if status == "APPROVE" {
		UpdateStageInstallmentApproveOrReject(installmentId, status)
		fmt.Println("Installment data has been approved. Waiting worker. installmentId=" + convertedInstallmentId)
		return
	}

	/*
	*		UPDATE STATUS TO `PROCESSING`, ONCE THE CALCULATION IS DONE, THEN UPDATE STATUS TO `SUCCESS`
	 */

	UpdateStageInstallmentApproveOrReject(installmentId, "PROCESSING")

	/*
	*		START CALCULATION PROCESS
	 */

	fmt.Println("Start calculation process. installmentId=" + convertedInstallmentId)

	queryGetAccountInvestor := `SELECT r_loan_installment."loanId", r_investor_product_pricing_loan."investorId", r_account_investor."accountId", product_pricing."returnOfInvestment" as "pplROI" 
	FROM installment 
	JOIN r_loan_installment ON r_loan_installment."installmentId" = installment."id" 
	JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan."loanId" = r_loan_installment."loanId" 
	JOIN r_account_investor ON r_account_investor."investorId" = r_investor_product_pricing_loan."investorId"
	join product_pricing on product_pricing.id = r_investor_product_pricing_loan."productPricingId"
	WHERE installment."id" = ?`

	loanInvestorAccountIDSchema := LoanInvestorAccountID{}
	er := services.DBCPsql.Raw(queryGetAccountInvestor, installmentId).Scan(&loanInvestorAccountIDSchema).Error
	if er != nil {
		fmt.Println(er)
		return
	}

	loanSchema := LoanSchema{}
	services.DBCPsql.Table("loan").Where("id = ?", loanInvestorAccountIDSchema.LoanID).Scan(&loanSchema)

	// accountTransactionDebitAmount := frequency * (plafond / tenor) + ((paidInstallment - (frequency * (plafond/tenor))) * pplROI);
	freq := float64(installmentSchema.Frequency)
	plafond := loanSchema.Plafond
	tenor := float64(loanSchema.Tenor)
	paidInstallment := installmentSchema.PaidInstallment
	pplROI := loanInvestorAccountIDSchema.PPLROI

	accountTransactionDebitAmount := freq*(plafond/tenor) + ((paidInstallment - (freq * (plafond / tenor))) * pplROI)

	accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "INSTALLMENT", TransactionDate: time.Now(), Amount: accountTransactionDebitAmount}
	services.DBCPsql.Table("account_transaction_debit").Create(accountTransactionDebitSchema)

	rAccountTransactionDebit := &r.RAccountTransactionDebit{AccountId: loanInvestorAccountIDSchema.AccountID, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
	services.DBCPsql.Table("r_account_transaction_debit").Create(rAccountTransactionDebit)

	// querySumDebitAndCredit := "SELECT SUM(account_transaction_debit.\"amount\") as \"totalDebit\", SUM(account_transaction_credit.\"amount\")  as \"totalCredit\" "
	// querySumDebitAndCredit += "FROM account "
	// querySumDebitAndCredit += "LEFT JOIN r_account_transaction_debit ON r_account_transaction_debit.\"accountId\" = account.\"id\" "
	// querySumDebitAndCredit += "LEFT JOIN account_transaction_debit ON account_transaction_debit.\"id\" = r_account_transaction_debit.\"accountTransactionDebitId\" "
	// querySumDebitAndCredit += "LEFT JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountId\" = account.\"id\" "
	// querySumDebitAndCredit += "LEFT JOIN account_transaction_credit ON account_transaction_credit.\"id\" = r_account_transaction_credit.\"accountTransactionCreditId\" "
	// querySumDebitAndCredit += "WHERE account.\"id\" = ?"

	// accountTransactionDebitAndCreditSchema := AccountTransactionDebitAndCredit{}
	// services.DBCPsql.Raw(querySumDebitAndCredit, loanInvestorAccountIDSchema.AccountID).Scan(&accountTransactionDebitAndCreditSchema)

	// totalBalance := accountTransactionDebitAndCreditSchema.TotalDebit - accountTransactionDebitAndCreditSchema.TotalCredit
	// services.DBCPsql.Table("account").Exec("UPDATE account SET \"totalDebit\" = ?, \"totalCredit\" = ?, \"totalBalance\" = ? WHERE \"id\" = ?", accountTransactionDebitAndCreditSchema.TotalDebit, accountTransactionDebitAndCreditSchema.TotalCredit, totalBalance, loanInvestorAccountIDSchema.AccountID)

	totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(loanInvestorAccountIDSchema.AccountID)
	totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(loanInvestorAccountIDSchema.AccountID)

	totalBalance := totalDebit - totalCredit
	services.DBCPsql.Table("account").Where("id = ?", loanInvestorAccountIDSchema.AccountID).Updates(account.Account{TotalDebit: totalDebit, TotalCredit: totalCredit, TotalBalance: totalBalance})

	fmt.Println("Calculation process has been done. installmentId=" + convertedInstallmentId)

	/*
	*		CALCULATION IS DONE, UPDATE INSTALLMENT STATUS FROM `PROCESSING` TO `APPROVE`
	 */

	UpdateStageInstallmentApproveOrReject(installmentId, status)
	// ctx.JSON(iris.StatusOK, iris.Map{
	// 	"status": "success",
	// 	"data": iris.Map{
	// 		"message": "Installment has been updated to " + status,
	// 	},
	// })
}

// UpdateStageInstallmentApproveOrReject - Update installment stage
func UpdateStageInstallmentApproveOrReject(installmentId uint64, status string) {
	convertedInstallmentID := strconv.FormatUint(installmentId, 10)
	fmt.Println("Updating status to " + status + ". installmentId=" + convertedInstallmentID)

	installmentHistorySchema := &installmentHistory.InstallmentHistory{StageFrom: "PENDING", StageTo: status}
	services.DBCPsql.Table("installment_history").Create(installmentHistorySchema)

	installmentHistoryID := installmentHistorySchema.ID

	rInstallmentHistorySchema := &r.RInstallmentHistory{InstallmentId: installmentId, InstallmentHistoryId: installmentHistoryID}
	services.DBCPsql.Table("r_installment_history").Create(rInstallmentHistorySchema)

	services.DBCPsql.Table("installment").Where("\"id\" = ?", installmentId).UpdateColumn("stage", status)

	fmt.Println("Done. Updated status to " + status + ". installmentId=" + convertedInstallmentID)
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

	if strings.ToLower(ctx.Param("status")) == "approve" || strings.ToLower(ctx.Param("status")) == "reject" || strings.ToLower(ctx.Param("status")) == "in-review" || strings.ToLower(ctx.Param("status")) == "success" {
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

		installmentDetailSchema := []InstallmentDetail{}
		if strings.ToLower(ctx.Param("status")) == "success" {
			query += "WHERE installment.\"stage\" = 'APPROVE'"
			services.DBCPsql.Raw(query).Scan(&installmentDetailSchema)
		} else {
			query += "WHERE installment.\"createdAt\"::date = ? AND \"group\".\"id\" = ? AND installment.\"stage\" != 'APPROVE'"
			services.DBCPsql.Raw(query, transactionDate, groupID).Scan(&installmentDetailSchema)
		}

		for _, item := range installmentDetailSchema {
			// go storeInstallment(item.InstallmentID, status)
			storeInstallment(item.InstallmentID, status)
		}

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data": iris.Map{
				"message": "Your request has been received. It might need take a while to process your request.",
			},
		})
	} else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status": "error",
			"data": iris.Map{
				"message": "Invalid status.",
			},
		})
	}
}

//
func SubmitInstallmentByGroupIDAndTransactionDateWithStatusAndInstallmentId(ctx *iris.Context) {
	key := ctx.URLParam("ais")

	if key == "" {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{
			"status": "error",
			"data": iris.Map{
				"message": "Unauthorized access.",
			},
		})
		return
	}

	query := "SELECT  "
	query += "\"group\".\"id\" as \"groupId\", \"group\".\"name\" as \"groupName\", "
	query += "installment.\"id\" as \"installmentId\", installment.\"type\", installment.\"paidInstallment\", installment.\"penalty\", installment.\"reserve\", installment.\"presence\", installment.\"frequency\", installment.\"stage\"  "
	query += "FROM installment  "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\"  "
	query += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\"  "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\"  "
	query += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\"   "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\"  "
	query += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\"  "
	query += "WHERE installment.id = 979763 AND installment.stage = 'APPROVE' "

	installmentDetailSchema := []InstallmentDetail{}
	services.DBCPsql.Raw(query).Scan(&installmentDetailSchema)

	for _, item := range installmentDetailSchema {
		// go storeInstallment(item.InstallmentID, status)
		storeInstallment(item.InstallmentID, "SUCCESS")
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"message": "Your request has been received. It might need take a while to process your request.",
		},
	})
}
