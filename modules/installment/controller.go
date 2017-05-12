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
	"github.com/jinzhu/gorm"
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

	query := `
		SELECT 
		branch."name" AS "branch",
		"group"."id" AS "groupId",
		"group"."name" AS "group",
		SUM(installment."paidInstallment") AS "totalPaidInstallment",
		SUM(installment.reserve) AS "totalReserve", 
		installment."createdAt"::date
		FROM installment
		JOIN r_loan_installment ON r_loan_installment."installmentId" = installment."id"
		JOIN loan ON loan."id" = r_loan_installment."loanId"
		JOIN r_loan_branch ON r_loan_branch."loanId" = loan."id"
		JOIN branch ON branch."id" = r_loan_branch."branchId" 
		JOIN r_loan_group ON r_loan_group."loanId" = loan."id"
		JOIN "group" ON "group"."id" = r_loan_group."groupId"
		WHERE installment.stage = ? AND branch.id = ?
		GROUP BY installment."createdAt"::date, branch."name", "group"."id", "group"."name"
		ORDER BY installment."createdAt"::date DESC, branch."name" ASC
	`
	err := services.DBCPsql.Raw(query, installmentType, branchID).Find(&installments).Error
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"data": err})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{"data": installments})
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
	err := services.DBCPsql.Raw(query, transactionDate, groupID, branchID).Scan(&installmentDetailSchema).Error
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"data": err})
		return
	}

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

func StoreInstallment(installmentId uint64, status string) {
	convertedInstallmentId := strconv.FormatUint(installmentId, 10)
	fmt.Println("[INFO] Storing installment. installmentID=" + convertedInstallmentId + " status=" + status)
	installmentSchema := Installment{}
	services.DBCPsql.Table("installment").Where("\"id\" = ?", installmentId).First(&installmentSchema)

	if installmentSchema.Stage != "PENDING" && installmentSchema.Stage != "IN-REVIEW" && installmentSchema.Stage != "APPROVE" {
		fmt.Println("Current installment stage is NEITHER 'PENDING' NOR 'IN-REVIEW' nor 'APPROVE'. System cannot continue to process your request. installmentId=" + convertedInstallmentId)
		return
	}

	if status == "REJECT" {
		UpdateStageInstallmentApproveOrReject(installmentId, installmentSchema.Stage, status)
		fmt.Println("Installment data has been rejected. installmentId=" + convertedInstallmentId)
		return
	}

	if status == "IN-REVIEW" {
		UpdateStageInstallmentApproveOrReject(installmentId, installmentSchema.Stage, status)
		fmt.Println("Installment data will be reviewed. installmentId=" + convertedInstallmentId)
		return
	}

	if status == "APPROVE" {
		UpdateStageInstallmentApproveOrReject(installmentId, installmentSchema.Stage, status)
		fmt.Println("Installment data has been approved. Waiting worker. installmentId=" + convertedInstallmentId)
		return
	}

	/*
	*		UPDATE STATUS TO `PROCESSING`, ONCE THE CALCULATION IS DONE, THEN UPDATE STATUS TO `SUCCESS`
	 */

	UpdateStageInstallmentApproveOrReject(installmentId, installmentSchema.Stage, "PROCESSING")

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

	rAccountTransactionDebitInstallmentData := r.RAccountTransactionDebitInstallment{InstallmentId: installmentId, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
	services.DBCPsql.Table("r_account_transaction_debit_installment").Create(&rAccountTransactionDebitInstallmentData)

	totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(loanInvestorAccountIDSchema.AccountID)
	totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(loanInvestorAccountIDSchema.AccountID)

	totalBalance := totalDebit - totalCredit
	services.DBCPsql.Table("account").Where("id = ?", loanInvestorAccountIDSchema.AccountID).Updates(account.Account{TotalDebit: totalDebit, TotalCredit: totalCredit, TotalBalance: totalBalance})

	fmt.Println("Calculation process has been done. installmentId=" + convertedInstallmentId)

	/*
	*		CALCULATION IS DONE, UPDATE INSTALLMENT STATUS FROM `PROCESSING` TO `SUCCESS`
	*/

	if err := UpdateLoanStage(installmentSchema, loanSchema.ID, services.DBCPsql); err != nil {
		fmt.Printf("Skip Update loan %d to END is failed. Error = %s\n", loanSchema.ID, err.Error())
	}

	UpdateStageInstallmentApproveOrReject(installmentId, "PROCESSING", status)
}

// UpdateStageInstallmentApproveOrReject - Update installment stage
func UpdateStageInstallmentApproveOrReject(installmentId uint64, stageFrom string, status string) {
	convertedInstallmentID := strconv.FormatUint(installmentId, 10)
	fmt.Println("Updating status to " + status + ". installmentId=" + convertedInstallmentID)

	installmentHistorySchema := &installmentHistory.InstallmentHistory{StageFrom: stageFrom, StageTo: status}
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

	go StoreInstallment(installmentID, status)

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
			// go StoreInstallment(item.InstallmentID, status)
			StoreInstallment(item.InstallmentID, status)
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
		// go StoreInstallment(item.InstallmentID, status)
		StoreInstallment(item.InstallmentID, "SUCCESS")
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"message": "Your request has been received. It might need take a while to process your request.",
		},
	})
}

// PendingInstallmentSchema - struct - pending installment
type PendingInstallmentSchema struct {
	BranchID                     uint64  `gorm:"column:branchId" json:"branchId"`
	Branch                       string  `gorm:"column:branch" json:"branch"`
	GroupID                      uint64  `gorm:"column:groupId" json:"groupId"`
	Group                        string  `gorm:"column:group" json:"group"`
	TotalAmountInstallmentPerDay float64 `gorm:"column:totalInstallmentPerDay" json:"totalAmountInstallmentPerDay"`
	TotalPaidInstallment         float64 `gorm:"column:paidInstallment" json:"totalPaidInstallment"`
}

// GetPendingInstallment - get pending installment
func GetPendingInstallment(ctx *iris.Context) {
	branchID := ctx.Param("branch_id")
	scheduleDay := ctx.Param("schedule_day")

	query := `
		select b."groupIdeal",a.* from (
		SELECT count(distinct(loan.id)),branch.id AS "branchId", branch."name" AS "branch", 
		"group".id AS "groupId", "group"."name" AS "group", 
		sum((loan.installment*installment.frequency)) AS "totalInstallmentPerDay", sum(installment."paidInstallment") AS "paidInstallment", 
		array_agg(loan.id) AS "loanIds",
		array_agg(installment.id) AS "installmentIds"
		FROM branch
		JOIN r_loan_branch ON r_loan_branch."branchId" = branch.id
		JOIN loan ON loan.id = r_loan_branch."loanId"
		JOIN r_loan_group ON r_loan_group."loanId" = loan.id
		JOIN "group" ON "group"."id" = r_loan_group."groupId"
		LEFT JOIN r_loan_installment ON r_loan_installment."loanId" = loan.id
		LEFT JOIN installment ON installment.id = r_loan_installment."installmentId"
		WHERE branch.id = ?
		AND "group"."scheduleDay" = ?
		AND installment.stage = 'PENDING'
		AND loan."deletedAt" IS NULL 
		AND "group"."deletedAt" IS NULL
		AND "branch"."deletedAt" IS NULL
		AND installment."deletedAt" IS null and loan."stage" = 'INSTALLMENT'
		GROUP BY "group".id, branch.id) a join
		(select g.id "gId", count(distinct(l.id)) "groupIdeal" from loan l join r_loan_group rlg on l.id = rlg."loanId" join "group" g on g.id = rlg."groupId" where g."scheduleDay" = 'Rabu' and l.stage = 'INSTALLMENT' group by "gId") b on b."gId"=a."groupId"
	`

	pendingInstallmentSchema := []PendingInstallmentSchema{}
	services.DBCPsql.Raw(query, branchID, scheduleDay).Scan(&pendingInstallmentSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   pendingInstallmentSchema,
	})
}

type PendingInstallmentDetailSchema struct {
	BorrowerNo           uint64  `gorm:"column:borrowerNo" json:"borrowerNo"`
	Borrower             string  `gorm:"column:borrower" json:"borrower"`
	LoanID               uint64  `gorm:"column:loanId" json:"loanId"`
	BaseInstallment      float64 `gorm:"column:baseInstallment" json:"baseInstallment"`
	ProjectedInstallment float64 `gorm:"column:projectedInstallment" json:"projectedInstallment"`
	PaidInstallment      float64 `gorm:"column:paidInstallment" json:"paidInstallment"`
	Frequency            uint32  `gorm:"column:frequency" json:"frequency"`
	Reserve              float64 `gorm:"column:reserve" json:"reserve"`
	InstallmentID        uint64  `gorm:"column:installmentId" json:"installmentId"`
	Type                 string  `gorm:"column:type" json:"type"`
	CreatedAt            string  `gorm:"column:installmentCreatedAt" json:"installmentCreatedAt"`
}

// GetPendingInstallmentDetail - get pending installment detail
func GetPendingInstallmentDetail(ctx *iris.Context) {
	groupID := ctx.Param("group_id")

	query := `
		SELECT
		borrower."borrowerNo", cif.name AS "borrower",
		loan.id AS "loanId", loan.installment AS "baseInstallment", (loan.installment*installment.frequency) AS "projectedInstallment",
		installment."paidInstallment", installment.frequency, installment.reserve, installment.id AS "installmentId", installment."type", installment."createdAt"::date AS "installmentCreatedAt"
		FROM installment
		JOIN r_loan_installment ON r_loan_installment."installmentId" = installment.id
		JOIN r_loan_group ON r_loan_group."loanId" = r_loan_installment."loanId"
		JOIN loan ON loan.id = r_loan_installment."loanId"
		JOIN r_loan_borrower ON r_loan_borrower."loanId" = loan.id
		JOIN borrower ON borrower.id = r_loan_borrower."borrowerId"
		JOIN r_cif_borrower ON r_cif_borrower."borrowerId" = borrower.id
		JOIN cif ON cif.id = r_cif_borrower."cifId"
		WHERE r_loan_group."groupId" = ? AND installment.stage = 'PENDING' AND 
		installment."deletedAt" IS NULL AND loan."deletedAt" IS NULL AND borrower."deletedAt" IS NULL AND cif."deletedAt" IS null
	`

	pendingInstallmentDetailSchema := []PendingInstallmentDetailSchema{}
	services.DBCPsql.Raw(query, groupID).Scan(&pendingInstallmentDetailSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   pendingInstallmentDetailSchema,
	})
}

type UpdateInstallmentJSON struct {
	PaidInstallment float64 `json:"paidInstallment"`
	Reserve         float64 `json:"reserve"`
}

// UpdateInstallmentByInstallmentID - update installment data
func UpdateInstallmentByInstallmentID(ctx *iris.Context) {
	installmentID := ctx.Param("installment_id")

	updateInstallmentJSON := UpdateInstallmentJSON{}
	if err := ctx.ReadJSON(&updateInstallmentJSON); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err,
		})
		return
	}

	fmt.Printf("%+v", updateInstallmentJSON)

	query := `
		UPDATE installment
		SET "paidInstallment" = ?, reserve = ?, "updatedAt" = current_timestamp
		WHERE installment.id = ?
	`

	if err := services.DBCPsql.Exec(query, updateInstallmentJSON.PaidInstallment, updateInstallmentJSON.Reserve, installmentID).Error; err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err,
		})
		return
	}

	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})

}

func UpdateLoanStage(installment Installment, loanId uint64, db *gorm.DB) error {
	var loan = struct {
		ID        string `gorm:"column:id"`
		Frequency int32 `gorm:"column:frequency"`
		Tenor     int32 `gorm:"column:tenor"`
	}{}

	query := ` SELECT loan.id as id, SUM(frequency) as frequency, tenor FROM loan JOIN r_loan_installment on loan.id = "loanId" JOIN installment on installment.id = "installmentId" WHERE loan.id = ? AND installment.stage = 'SUCCESS' AND installment."deletedAt" isnull AND r_loan_installment."deletedAt" isnull GROUP BY loan.id`

	if err := db.Raw(query, loanId).Scan(&loan).Error; err != nil {
		return err
	}

	if loan.Frequency + installment.Frequency < loan.Tenor {
		// frequency is below tenor so dont go on
		return nil
	}

	if installment.Frequency < 3 {
		if err := db.Table("loan").Where("id = ?", loanId).UpdateColumn("stage", "END").Error; err != nil {
			return err
		}
		return nil
	}

	if installment.Frequency >= 3 {
		if err := db.Table("loan").Where("id = ?", loanId).UpdateColumn("stage", "END-EARLY").Error; err != nil {
			return err
		}
		return nil
	}

	// supposed not to go here
	//return error.New("Somethings is wrong")
	return nil
}
