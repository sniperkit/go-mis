package installment

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/go-mis/modules/account"
	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	installmentHistory "bitbucket.org/go-mis/modules/installment-history"
	loanHistory "bitbucket.org/go-mis/modules/loan-history"
	"bitbucket.org/go-mis/modules/r"
	systemParameter "bitbucket.org/go-mis/modules/system-parameter"
	"bitbucket.org/go-mis/services"
	"github.com/jinzhu/gorm"
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
	query += "WHERE installment.stage = 'PENDING' AND branch.id = ? "
	query += "AND installment.\"deletedAt\" IS NULL "
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
		WHERE installment.stage = ? AND branch.id = ? AND installment."deletedAt" IS NULL
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
	stage := ctx.Param("stage")
	transactionDate := ctx.Param("transaction_date")

	query := `
			select g.name,cif.name as "borrowerName", sum(i."paidInstallment") "repayment",sum(i.reserve) "tabungan",sum(i."paidInstallment"+i.reserve) "total",
		bow.id as "borrowerId",
		i.id as "installmentId",
		i.type,          
    i."paidInstallment", 
    i.penalty,       
    i.reserve, 
    i.presence, 
    i.frequency, 
    i.stage, 
        sum(i.cash_on_hand) "cashOnHand",
        sum(i.cash_on_reserve) "cashOnReserve",
        coalesce(sum(
                case
                when frequency >= 3 then l.installment+((plafond/tenor)*(frequency-1))
                when frequency >0 then l.installment*frequency
                when frequency = 0 then 0
                end
                ),0) "projectionRepayment",
                coalesce(sum(
                case
                when plafond < 0 then 0
                when plafond <= 3000000 then 3000
                when plafond > 3000000 and plafond <= 5000000 then 4000
                when plafond > 5000000 and plafond <= 7000000 then 5000
                when plafond > 7000000 and plafond <= 9000000 then 6000
                when plafond > 9000000 and plafond <= 11000000 then 7000
                else 8000
                end
                ),0) "projectionTabungan",
        coalesce(sum(case
        when d.stage = 'SUCCESS' then plafond end
        ),0) "totalCair",
        coalesce(sum(case
        when d.stage = 'FAILED' then plafond end
        ),0) "totalGagalDropping",
        split_part(string_agg(i.stage,'| '),'|',1) "status"
        from loan l join r_loan_group rlg on l.id = rlg."loanId"
        join r_loan_borrower bow on bow."loanId"=l.id
        join r_cif_borrower rcif using("borrowerId") 
        join cif on cif.id = rcif."cifId"
        join "group" g on g.id = rlg."groupId"
        join r_group_agent rga on g.id = rga."groupId"
        join agent a on a.id = rga."agentId"
        join r_loan_branch rlb on rlb."loanId" = l.id
        join branch b on b.id = rlb."branchId"
        join r_loan_installment rli on rli."loanId" = l.id
        join installment i on i.id = rli."installmentId"
        join r_loan_disbursement rld on rld."loanId" = l.id
        join disbursement d on d.id = rld."disbursementId"
        where l."deletedAt" isnull and b.id= ? and coalesce(i."transactionDate",i."createdAt")::date = ?
        and l.stage = 'INSTALLMENT' and i.stage= ? and g.id=?
		group by l.id, i.id, bow.id, g.name, cif.name,i.type,i."paidInstallment", i.penalty, i.reserve, i.presence, i.frequency, i.stage, i.cash_on_hand, i.cash_on_reserve 
			`

	installmentDetailSchema := []InstallmentDetail{}
	err := services.DBCPsql.Raw(query, branchID, transactionDate, stage, groupID).Scan(&installmentDetailSchema).Error
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

func StoreInstallment(db *gorm.DB, installmentId uint64, status string) error {
	convertedInstallmentId := strconv.FormatUint(installmentId, 10)
	fmt.Println("[INFO] Storing installment. installmentID=" + convertedInstallmentId + " status=" + status)
	installmentSchema := Installment{}
	db.Table("installment").Where("\"id\" = ?", installmentId).First(&installmentSchema)

	if installmentSchema.Stage != "TELLER" && installmentSchema.Stage != "AGENT" && installmentSchema.Stage != "PENDING" && installmentSchema.Stage != "IN-REVIEW" && installmentSchema.Stage != "APPROVE" {
		return errors.New("Current installment stage is NEITHER 'PENDING' NOR 'IN-REVIEW' nor 'APPROVE'. System cannot continue to process your request. installmentId=" + convertedInstallmentId)
	}

	if strings.ToUpper(status) == "REJECT" || strings.ToUpper(status) == "IN-REVIEW" || strings.ToUpper(status) == "APPROVE" || strings.ToUpper(status) == "AGENT" || strings.ToUpper(status) == "TELLER" {
		log.Println("Installment data has been", status, ". Waiting worker. installmentId=", convertedInstallmentId)
		UpdateStageInstallmentApproveOrReject(db, installmentId, installmentSchema.Stage, status)
		return nil
	}

	/*
	*		UPDATE STATUS TO `PROCESSING`, ONCE THE CALCULATION IS DONE, THEN UPDATE STATUS TO `SUCCESS`
	 */

	UpdateStageInstallmentApproveOrReject(db, installmentId, installmentSchema.Stage, "PROCESSING")

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
	er := db.Raw(queryGetAccountInvestor, installmentId).Scan(&loanInvestorAccountIDSchema).Error
	if er != nil {
		return er
	}

	loanSchema := LoanSchema{}
	db.Table("loan").Where("id = ?", loanInvestorAccountIDSchema.LoanID).Scan(&loanSchema)

	// Recheck paidInstallment and update to END/END EARLY if true
	if err := UpdateLoanStage(installmentSchema, loanSchema.ID, services.DBCPsql); err != nil {
		return fmt.Errorf("Error on Update Loan Stage. Error = %s\n", loanSchema.ID, err.Error())
	}

	// accountTransactionDebitAmount := frequency * (plafond / tenor) + ((paidInstallment - (frequency * (plafond/tenor))) * pplROI);
	freq := float64(installmentSchema.Frequency)
	plafond := loanSchema.Plafond
	tenor := float64(loanSchema.Tenor)
	paidInstallment := installmentSchema.PaidInstallment
	pplROI := loanInvestorAccountIDSchema.PPLROI

	accountTransactionDebitAmount := freq*(plafond/tenor) + ((paidInstallment - (freq * (plafond / tenor))) * pplROI)

	accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "INSTALLMENT", TransactionDate: time.Now(), Amount: accountTransactionDebitAmount}
	db.Table("account_transaction_debit").Create(accountTransactionDebitSchema)

	rAccountTransactionDebit := &r.RAccountTransactionDebit{AccountId: loanInvestorAccountIDSchema.AccountID, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
	db.Table("r_account_transaction_debit").Create(rAccountTransactionDebit)

	rAccountTransactionDebitInstallmentData := r.RAccountTransactionDebitInstallment{InstallmentId: installmentId, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
	db.Table("r_account_transaction_debit_installment").Create(&rAccountTransactionDebitInstallmentData)

	totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(loanInvestorAccountIDSchema.AccountID)
	totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(loanInvestorAccountIDSchema.AccountID)

	totalBalance := totalDebit - totalCredit
	db.Table("account").Where("id = ?", loanInvestorAccountIDSchema.AccountID).Updates(account.Account{TotalDebit: totalDebit, TotalCredit: totalCredit, TotalBalance: totalBalance})

	fmt.Println("Calculation process has been done. installmentId=" + convertedInstallmentId)

	/*
	*		CALCULATION IS DONE, UPDATE INSTALLMENT STATUS FROM `PROCESSING` TO `SUCCESS`
	 */

	UpdateStageInstallmentApproveOrReject(db, installmentId, "PROCESSING", status)
	return nil
}

// UpdateStageInstallmentApproveOrReject - Update installment stage
func UpdateStageInstallmentApproveOrReject(db *gorm.DB, installmentId uint64, stageFrom string, status string) error {
	var err error
	convertedInstallmentID := strconv.FormatUint(installmentId, 10)
	fmt.Println("Updating status to " + status + ". installmentId=" + convertedInstallmentID)

	installmentHistorySchema := &installmentHistory.InstallmentHistory{StageFrom: stageFrom, StageTo: status}
	if err = db.Table("installment_history").Create(installmentHistorySchema).Error; err != nil {
		return err
	}

	installmentHistoryID := installmentHistorySchema.ID

	rInstallmentHistorySchema := &r.RInstallmentHistory{InstallmentId: installmentId, InstallmentHistoryId: installmentHistoryID}
	if err = db.Table("r_installment_history").Create(rInstallmentHistorySchema).Error; err != nil {
		return err
	}

	if err = db.Table("installment").Where("\"id\" = ?", installmentId).UpdateColumn("stage", status).Error; err != nil {
		return err
	}

	fmt.Println("Done. Updated status to " + status + ". installmentId=" + convertedInstallmentID)
	return nil
}

// SubmitInstallmentByInstallmentIDWithStatus - approve or reject installment by installment_id
func SubmitInstallmentByInstallmentIDWithStatus(ctx *iris.Context) {
	installmentID, _ := strconv.ParseUint(ctx.Param("installment_id"), 10, 64)
	status := strings.ToUpper(ctx.Param("status"))

	go func() {
		db := services.DBCPsql.Begin()
		err := StoreInstallment(db, installmentID, status)
		if err != nil {
			ProcessErrorAndRollback(ctx, db, err.Error())
			return
		}
		db.Commit()
	}()

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
	stageTo := strings.ToUpper(ctx.Param("stageTo"))
	stageFrom := ctx.Param("stageFrom")

	if strings.ToLower(stageTo) == "agent" || strings.ToLower(stageTo) == "teller" || strings.ToLower(stageTo) == "pending" || strings.ToLower(stageTo) == "approve" || strings.ToLower(stageTo) == "reject" || strings.ToLower(stageTo) == "in-review" || strings.ToLower(stageTo) == "success" {
		query := "SELECT "
		query += "\"group\".\"id\" as \"groupId\", \"group\".\"name\" as \"groupName\","
		query += "installment.\"id\" as \"installmentId\", installment.\"type\", installment.\"paidInstallment\", installment.\"penalty\", installment.\"reserve\", installment.\"presence\", installment.\"frequency\", installment.\"stage\", branch.\"id\" "
		query += "FROM installment "
		query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
		query += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
		query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
		query += "JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\"  "
		query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
		query += "JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
		db := services.DBCPsql.Begin()
		installmentDetailSchema := []InstallmentDetail{}
		if strings.ToLower(stageTo) == "success" {
			query += "WHERE installment.\"stage\" = 'APPROVE'"
			db.Raw(query).Scan(&installmentDetailSchema)
		} else {
			query += "WHERE installment.\"createdAt\"::date = ? AND \"group\".\"id\" = ? AND installment.\"stage\" != 'APPROVE'"
			db.Raw(query, transactionDate, groupID).Scan(&installmentDetailSchema)
		}

		for _, item := range installmentDetailSchema {
			// go StoreInstallment(item.InstallmentID, status)
			err := StoreInstallment(db, item.InstallmentID, stageTo)
			if err != nil {
				fmt.Println(err)
				ProcessErrorAndRollback(ctx, db, err.Error())
				return
			}
		}
		db.Commit()

		// write to go-log

		tempGid, _ := strconv.Atoi(groupID)
		gid := uint64(tempGid)
		inst := struct {
			GroupID     uint64
			Date        string
			StageFrom   string
			StageTo     string
			Installment []InstallmentDetail
		}{
			GroupID:     gid,
			Date:        transactionDate,
			StageFrom:   stageFrom,
			StageTo:     stageTo,
			Installment: installmentDetailSchema,
		}
		_ = services.PostToLog(services.GetLog(gid, inst, stageTo))

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
	db := services.DBCPsql.Begin()
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
	db.Raw(query).Scan(&installmentDetailSchema)

	for _, item := range installmentDetailSchema {
		// go StoreInstallment(item.InstallmentID, status)
		err := StoreInstallment(db, item.InstallmentID, "SUCCESS")
		if err != nil {
			ProcessErrorAndRollback(ctx, db, err.Error())
			return
		}
	}
	db.Commit()
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

type SimpleLoan struct {
	ID          string  `gorm:"column:id"`
	Plafond     int32   `gorm:"column:plafond"`
	Installment float64 `gorm:"column:installment"`
	Frequency   int32   `gorm:"column:frequency"`
	Tenor       int32   `gorm:"column:tenor"`
	Rate        float32 `gorm:"column:rate"`
	Stage       string  `gorm:"column:stage"`
}

/**
 * sorry there are some calculation
 * in here not only Update LoanStage
 *
 *
 *
 */
func UpdateLoanStage(installment Installment, loanId uint64, db *gorm.DB) error {
	var loan = SimpleLoan{}
	query := `
	SELECT 
	loan.id as id, loan.plafond as plafond, loan.installment as installment, SUM(frequency) as frequency, tenor, rate, loan.stage as stage 
	FROM loan 
	LEFT JOIN r_loan_installment on loan.id = "loanId" AND r_loan_installment."deletedAt" isnull 
	LEFT JOIN installment on installment.id = "installmentId" AND installment.stage = 'SUCCESS' AND installment."deletedAt" isnull
	WHERE loan.id = ?
	GROUP BY loan.id
	`

	if err := db.Raw(query, loanId).Scan(&loan).Error; err != nil {
		return err
	}

	if installment.Type == "MENINGGAL" {

		stageTo := "MENINGGAL"

		if err := db.Table("loan").Where("id = ?", loanId).UpdateColumn("stage", stageTo).Error; err != nil {
			return err
		}

		loanHistoryData := loanHistory.LoanHistory{StageFrom: loan.Stage, StageTo: stageTo, Remark: fmt.Sprintf("Automatic update stage %s loanId = %d", stageTo, loanId)}
		if err := db.Table("loan_history").Create(&loanHistoryData).Error; err != nil {
			return err
		}

		rLoanHistory := r.RLoanHistory{LoanId: loanId, LoanHistoryId: loanHistoryData.ID}
		if err := db.Table("r_loan_history").Create(&rLoanHistory).Error; err != nil {
			return err
		}
		return nil
	}

	if loan.Frequency+installment.Frequency < loan.Tenor {
		// frequency is below tenor so dont go on
		return nil
	}

	stageTo, calculationError := GetStageTo(installment, loan)

	if err := db.Table("loan").Where("id = ?", loanId).UpdateColumn("stage", stageTo).Error; err != nil {
		return err
	}

	loanHistoryData := loanHistory.LoanHistory{StageFrom: loan.Stage, StageTo: stageTo, Remark: fmt.Sprintf("Automatic update stage %s loanId = %d", stageTo, loanId)}
	if err := db.Table("loan_history").Create(&loanHistoryData).Error; err != nil {
		return err
	}

	rLoanHistory := r.RLoanHistory{LoanId: loanId, LoanHistoryId: loanHistoryData.ID}
	if err := db.Table("r_loan_history").Create(&rLoanHistory).Error; err != nil {
		return err
	}

	// supposed not to go here
	//return error.New("Somethings is wrong")
	return calculationError
}

func GetStageTo(installment Installment, loan SimpleLoan) (string, error) {

	if installment.Frequency == 1 {
		return "END", nil
	}

	installmentProfit := int32(installment.PaidInstallment) - (loan.Plafond / loan.Tenor * installment.Frequency)

	profit1x := int32(float32(loan.Plafond)*loan.Rate) / loan.Tenor

	if installmentProfit == profit1x {
		return "END-EARLY", nil
	}

	if installmentProfit == profit1x*installment.Frequency {
		return "END", nil
	}

	return "END-PENDING", errors.New("Calculation End or End Early not match")
}

// FindByBranchAndDate - Filter Installment by branch and date
func FindByBranchAndDate(branchID uint64, transactionDate string) ([]Installment, error) {
	if branchID < 0 {
		return nil, errors.New("Branch ID can not be empty")
	}
	if len(strings.Trim(transactionDate, " ")) == 0 {
		return nil, errors.New("Transaction date can not be empty")
	}
	installemnts := make([]Installment, 0)
	query := `select installment.id,

					installment.type,
					installment.presence,
					installment."paidInstallment",
					installment.penalty,
					installment.reserve,
					installment.frequency,
					installment.stage,
					installment."transactionDate",
					installment."createdAt"
			FROM installment,
					r_loan_installment,
					loan,
					branch,
					r_loan_branch
			WHERE installment.id = r_loan_installment."installmentId" AND
			loan.id = r_loan_installment."loanId" AND
			loan.id = r_loan_branch."loanId" AND
			branch.id = r_loan_branch."branchId" AND
			installment."deletedAt" is null AND
			UPPER(installment.stage) = 'TELLER' AND
			branch.id = ? AND
			installment."createdAt"::date = ?`

	if err := services.DBCPsql.Raw(query, branchID, transactionDate).Scan(&installemnts).Error; err != nil {
		log.Println("#ERROR: ", err.Error())
		return nil, errors.New("Unable to retrieve installments")
	}
	return installemnts, nil
}

func ProcessErrorAndRollback(ctx *iris.Context, db *gorm.DB, message string) {
	log.Println("#Error", message)
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status":  "error",
		"message": message,
	})
}

// GetPendingInstallmentNew - Get data pending Installment
// Route: /api/v2/installment-pending/get/:currentStage/:branchId/:date
func GetPendingInstallmentNew(ctx *iris.Context) {
	bId := ctx.Param("branchId")
	intBid, _ := strconv.Atoi(bId)
	branchID := uint64(intBid)
	dateParam := ctx.Param("date")
	// Check branchID, if equal to 0 return error message to client
	if branchID == 0 {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":       iris.StatusBadRequest,
			"errorMessage": "Invalid Branch ID",
		})
		return
	}
	if len(strings.Trim(dateParam, " ")) == 0 {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":       iris.StatusBadRequest,
			"errorMessage": "Date can not be empty",
		})
		return
	}
	if !systemParameter.IsAllowedBackdate(dateParam) {
		log.Println("#ERROR: Not Allowed back date")
		ctx.JSON(405, iris.Map{
			"message":      "Not Allowed",
			"errorMessage": "View back date is not allowed",
		})
		return
	}
	res := GetDataPendingInstallment(ctx.Param("currentStage"), branchID, dateParam)
	notes, err := services.GetNotes(services.ConstructNotesGroupId(branchID, dateParam))
	if err != nil || len(notes) > 0 {
		borrowerNotes := services.GetBorrowerNotes(notes)
		majelisNotes := services.GetMajelisNotes(notes)
		if borrowerNotes != nil {
			res.BorrowerNotes = borrowerNotes
		}
		if borrowerNotes != nil {
			res.MajelisNotes = majelisNotes
		}
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   res,
	})
}

func GetDataPendingInstallment(currentStage string, branchId uint64, now string) PendingInstallment {
	var pendingInstallment PendingInstallment
	queryResult := []PendingRawInstallmentData{}
	query := `select g.id as "groupId", a.fullname,g.name, sum(i."paidInstallment") "repayment",sum(i.reserve) "tabungan",sum(i."paidInstallment"+i.reserve) "total",
				sum(i.cash_on_hand) "cashOnHand",
				sum(i.cash_on_reserve) "cashOnReserve",
				coalesce(sum(
                case
                when frequency >= 3 then l.installment+((plafond/tenor)*(frequency-1))
                when frequency >0 then l.installment*frequency
                when frequency = 0 then 0
                end
                ),0) "projectionRepayment",
                coalesce(sum(
                case
                when plafond < 0 then 0
                when plafond <= 3000000 then 3000
                when plafond > 3000000 and plafond <= 5000000 then 4000
                when plafond > 5000000 and plafond <= 7000000 then 5000
                when plafond > 7000000 and plafond <= 9000000 then 6000
                when plafond > 9000000 and plafond <= 11000000 then 7000
                else 8000
                end
                ),0) "projectionTabungan",
				coalesce(sum(case
                when d."disbursementDate"::date = current_date then plafond end
                ),0) "totalCairProj",
                coalesce(sum(case
                when d.stage = 'SUCCESS' and d."disbursementDate"::date = current_date then plafond end
                ),0) "totalCair",
				coalesce(sum(case
				when d.stage = 'FAILED' then plafond end
				),0) "totalGagalDropping",
				split_part(string_agg(i.stage,'| '),'|',1) "status"
				from loan l join r_loan_group rlg on l.id = rlg."loanId"
				join "group" g on g.id = rlg."groupId"
				join r_group_agent rga on g.id = rga."groupId"
				join agent a on a.id = rga."agentId"
				join r_loan_branch rlb on rlb."loanId" = l.id
				join branch b on b.id = rlb."branchId"
				join r_loan_installment rli on rli."loanId" = l.id
				join installment i on i.id = rli."installmentId"
				join r_loan_disbursement rld on rld."loanId" = l.id
				join disbursement d on d.id = rld."disbursementId"`
	if currentStage == "in-review" {
		parseNow, _ := time.Parse("2006-01-02", now)
		yesterday := parseNow.AddDate(0, 0, -1).Format("2006-01-02")
		query += `where l."deletedAt" isnull and b.id= ? and coalesce(i."transactionDate",i."createdAt")::date <= ? and coalesce(i."transactionDate",i."createdAt")::date >= ? and l.stage = 'INSTALLMENT'
				group by g.name, a.fullname, g.id
				order by a.fullname`
		services.DBCPsql.Raw(query, branchId, now, yesterday).Scan(&queryResult)
	} else {
		query += `where l."deletedAt" isnull and b.id= ? and coalesce(i."transactionDate",i."createdAt")::date = ? and l.stage = 'INSTALLMENT'
				group by g.name, a.fullname, g.id
				order by a.fullname`
		services.DBCPsql.Raw(query, branchId, now).Scan(&queryResult)
	}

	res := []PendingInstallmentData{}
	agents := map[string]bool{"": false}
	for _, val := range queryResult {
		if agents[val.Fullname] == false {
			agents[val.Fullname] = true
			res = append(res, PendingInstallmentData{Agent: val.Fullname})
		}
	}
	majelisIDs := make([]MajelisId, len(res))
	for idx, rval := range res {
		var totalRepaymentAct float64
		var totalRepaymentProj float64
		var totalRepaymentCoh float64
		var totalTabunganAct float64
		var totalTabunganProj float64
		var totalTabunganCoh float64
		var totalActualAgent float64
		var totalProjectionAgent float64
		var totalCohAgent float64
		var totalPencairanAgent float64
		var totalPencairanProjAgent float64
		var totalGagalDroppingAgent float64
		m := []Majelis{}
		for _, qrval := range queryResult {
			if rval.Agent == qrval.Fullname {
				m = append(m, Majelis{
					GroupId:             qrval.GroupId,
					Name:                qrval.Name,
					Repayment:           qrval.Repayment,
					Tabungan:            qrval.Tabungan,
					TotalActual:         qrval.Total,
					TotalProyeksi:       qrval.ProjectionRepayment + qrval.ProjectionTabungan,
					TotalCoh:            qrval.CashOnHand + qrval.CashOnReserve,
					TotalCair:           qrval.TotalCair,
					TotalCairProj:       qrval.TotalCairProj,
					TotalGagalDropping:  qrval.TotalGagalDropping,
					Status:              qrval.Status,
					CashOnHand:          qrval.CashOnHand,
					CashOnReserve:       qrval.CashOnReserve,
					ProjectionRepayment: qrval.ProjectionRepayment,
					ProjectionTabungan:  qrval.ProjectionTabungan,
				})
				majelisIDs = append(majelisIDs, MajelisId{GroupId: qrval.GroupId, Name: qrval.Name})
				totalRepaymentAct += qrval.Repayment
				totalRepaymentProj += qrval.ProjectionRepayment
				totalRepaymentCoh += qrval.CashOnHand
				totalTabunganAct += qrval.Tabungan
				totalTabunganProj += qrval.ProjectionTabungan
				totalTabunganCoh += qrval.CashOnReserve
				totalActualAgent += qrval.Total
				totalProjectionAgent += qrval.ProjectionRepayment + qrval.ProjectionTabungan
				totalCohAgent += qrval.CashOnHand + qrval.CashOnReserve
				totalPencairanAgent += qrval.TotalCair
				totalPencairanProjAgent += qrval.TotalCairProj
				totalGagalDroppingAgent += qrval.TotalGagalDropping
			}
		}
		res[idx].Majelis = m
		res[idx].TotalActualRepayment = totalRepaymentAct
		res[idx].TotalProjectionRepayment = totalRepaymentProj
		res[idx].TotalCohRepayment = totalRepaymentCoh
		res[idx].TotalActualTabungan = totalTabunganAct
		res[idx].TotalProjectionTabungan = totalTabunganProj
		res[idx].TotalCohTabungan = totalTabunganCoh
		res[idx].TotalActualAgent = totalActualAgent
		res[idx].TotalProjectionAgent = totalProjectionAgent
		res[idx].TotalCohAgent = totalCohAgent
		res[idx].TotalPencairanAgent = totalPencairanAgent
		res[idx].TotalPencairanProjAgent = totalPencairanProjAgent
		res[idx].TotalGagalDroppingAgent = totalGagalDroppingAgent
	}
	pendingInstallment.ListMajelis = majelisIDs
	pendingInstallment.PendingInstallmentData = res
	return pendingInstallment
}
