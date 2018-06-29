package disbursement

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	disbursementHistory "bitbucket.org/go-mis/modules/disbursement-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/modules/feature-flag"
	"bitbucket.org/go-mis/modules/utility"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Disbursement{})
	services.BaseCrudInit(Disbursement{}, []Disbursement{})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchAll - fetchAll Disbursement data
func FetchAll(ctx *iris.Context) {
	branchID := ctx.Get("BRANCH_ID")
	disbursementFetchSchema := []DisbursementFetch{}

	query := "SELECT \"group\".id AS \"groupId\", \"group\".\"name\" AS \"group\", branch.id AS \"branchId\", branch.\"name\" AS \"branch\", SUM(loan.plafond) AS \"plafond\", loan.\"submittedLoanDate\"::date AS \"submittedLoanDate\", disbursement.\"disbursementDate\"::date AS \"disbursementDate\" "
	query += "FROM \"group\" "
	query += "JOIN r_group_branch ON r_group_branch.\"groupId\" = \"group\".id "
	query += "JOIN branch ON branch.id = r_group_branch.\"branchId\" "
	query += "JOIN r_loan_group ON r_loan_group.\"groupId\" = \"group\".id "
	query += "JOIN loan ON loan.id = r_loan_group.\"loanId\" "
	query += "JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.id "
	query += "JOIN disbursement ON disbursement.id = r_loan_disbursement.\"disbursementId\" "
	query += "WHERE disbursement.stage IN ('PENDING', 'FAILED') "
	query += "AND loan.\"submittedLoanDate\" IS NOT NULL "

	query += "AND DATE(disbursement.\"disbursementDate\") >= 'now()' "
	query += "AND branch.id = ? "
	query += "GROUP BY \"group\".id, branch.id, branch.\"name\", loan.\"submittedLoanDate\", disbursement.\"disbursementDate\" "
	query += "ORDER BY disbursement.\"disbursementDate\" ASC, \"group\".\"name\" ASC "

	services.DBCPsql.Raw(query, branchID).Scan(&disbursementFetchSchema)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   disbursementFetchSchema,
	})
}

func GetDisbursementDetailByGroup(ctx *iris.Context) {
	query := ""

	bid := utility.ParseBranchIDFromContext(ctx.Get("BRANCH_ID"))
	if feature_flag.Control.IsEnabledForBranchID("new-disbursement", bid) {
		query += "SELECT investor.id AS \"investorId\", \"group\".id AS \"groupId\", \"group\".\"name\" AS \"groupName\", branch.\"name\" AS \"branchName\", borrower.\"borrowerNo\", cif.\"name\" AS \"borrower\", loan.id AS \"loanId\", loan.plafond, disbursement.\"disbursementDate\"::date, disbursement.stage, loan.stage AS \"loanStage\", borrower.\"lwk1Date\", borrower.\"lwk2Date\", borrower.\"upkDate\", borrower.\"id\" as \"borrowerId\", "
		query += "case when loan.\"isLWK\" = true and loan.\"isUPK\" = true then true else false end as \"akadAvailable\" "
	} else {
		query += "SELECT investor.id AS \"investorId\", \"group\".id AS \"groupId\", \"group\".\"name\" AS \"groupName\", branch.\"name\" AS \"branchName\", borrower.\"borrowerNo\", cif.\"name\" AS \"borrower\", loan.id AS \"loanId\", loan.plafond, disbursement.\"disbursementDate\"::date, disbursement.stage, loan.stage AS \"loanStage\" "
	}

	query += "FROM \"group\" "
	query += "JOIN r_group_branch ON r_group_branch.\"groupId\" = \"group\".id "
	query += "JOIN branch ON branch.id = r_group_branch.\"branchId\" "
	query += "JOIN r_loan_group ON r_loan_group.\"groupId\" = \"group\".id "
	query += "JOIN loan ON loan.id = r_loan_group.\"loanId\" "
	query += "JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.id "
	query += "JOIN disbursement ON disbursement.id = r_loan_disbursement.\"disbursementId\" "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.id "
	query += "JOIN borrower ON borrower.id = r_loan_borrower.\"borrowerId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.id "
	query += "JOIN cif ON cif.id = r_cif_borrower.\"cifId\" "
	query += "JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"loanId\" = loan.id "
	query += "LEFT JOIN investor ON investor.\"id\" = r_investor_product_pricing_loan.\"investorId\" "
	query += "WHERE disbursement.stage IN ('PENDING', 'FAILED') "
	query += "AND loan.\"submittedLoanDate\" IS NOT NULL  "

	query += "AND DATE(disbursement.\"disbursementDate\") >= 'now()' "
	query += "AND branch.id = ? "
	query += "AND \"group\".id = ? "
	query += "AND disbursement.\"disbursementDate\"::date = ? "

	branchID := ctx.Param("branch_id")
	groupID := ctx.Param("group_id")
	disbursementDate := ctx.Param("disbursement_date")

	disbursementDetailByGroupSchema := []DisbursementDetailByGroup{}
	services.DBCPsql.Raw(query, branchID, groupID, disbursementDate).Scan(&disbursementDetailByGroupSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   disbursementDetailByGroupSchema,
	})
}

// UpdateDisbursementStage - Update disbursement stage and loan stage
func UpdateDisbursementStage(ctx *iris.Context) {
	stage := strings.ToLower(ctx.Param("stage"))
	loanID, errConvLoanID := strconv.ParseUint(ctx.Param("loan_id"), 10, 64)

	if errConvLoanID != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": errConvLoanID.Error(),
		})
		return
	}

	query := "SELECT disbursement.* FROM disbursement "
	query += "INNER JOIN r_loan_disbursement ON r_loan_disbursement.\"disbursementId\" = disbursement.id "
	query += "WHERE r_loan_disbursement.\"loanId\" = ?"

	disbursementSchema := Disbursement{}
	services.DBCPsql.Raw(query, loanID).Scan(&disbursementSchema)

	if stage == "success" {
		disbursementHistorySchema := &disbursementHistory.DisbursementHistory{StageFrom: "PENDING", StageTo: "SUCCESS"}
		services.DBCPsql.Table("disbursement_history").Create(disbursementHistorySchema)

		rDisbursementHistorySchema := &r.RDisbursementHistory{DisbursementId: disbursementSchema.ID, DisbursementHistoryId: disbursementHistorySchema.ID}
		services.DBCPsql.Table("r_disbursement_history").Create(rDisbursementHistorySchema)

		services.DBCPsql.Table("loan").Where("id = ?", loanID).UpdateColumn("stage", "INSTALLMENT")
		services.DBCPsql.Table("disbursement").Where("id = ?", disbursementSchema.ID).UpdateColumn("stage", "SUCCESS")

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
		return
	} else if stage == "failed" {
		jsonDisbursementStage := DisbursementStageInput{}
		if err := ctx.ReadJSON(&jsonDisbursementStage); err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		jsonDisbursementStage.UpdateDateValue()

		disbursementHistorySchema := &disbursementHistory.DisbursementHistory{StageFrom: "PENDING", StageTo: "FAILED", LastDisbursementDate: jsonDisbursementStage.LastDisbursementDate, NextDisbursementDate: jsonDisbursementStage.NextDisbursementDate}
		services.DBCPsql.Table("disbursement_history").Create(disbursementHistorySchema)

		rDisbursementHistorySchema := &r.RDisbursementHistory{DisbursementId: disbursementSchema.ID, DisbursementHistoryId: disbursementHistorySchema.ID}
		services.DBCPsql.Table("r_disbursement_history").Create(rDisbursementHistorySchema)

		services.DBCPsql.Table("loan").Where("id = ?", loanID).UpdateColumn("stage", "DROPPING-FAILED")
		services.DBCPsql.Table("disbursement").Where("id = ?", disbursementSchema.ID).UpdateColumn("disbursementDate", jsonDisbursementStage.NextDisbursementDate)

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	} else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "Invalid request.",
		})
	}
}

// UpdateStage - updateStage Disbursement data
func UpdateStage(ctx *iris.Context) {
	disbursementStageInput := DisbursementStageInput{}
	tempLoanID := ctx.Param("loan_id")
	stage := ctx.Param("stage")

	loanID, err := strconv.ParseUint(tempLoanID, 10, 64)

	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = ctx.ReadJSON(&disbursementStageInput)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	disbursementStageInput.UpdateDateValue()

	rLoanDisbursementData := r.RLoanDisbursement{}
	services.DBCPsql.Table("r_loan_disbursement").Where("\"loanId\" = ?", loanID).First(&rLoanDisbursementData)

	if rLoanDisbursementData == (r.RLoanDisbursement{}) {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Can't find any loan detail.",
		})
		return
	}

	disbursementData := Disbursement{}
	services.DBCPsql.First(&disbursementData, rLoanDisbursementData.DisbursementId)

	go services.DBCPsql.Table("disbursement").Where("\"id\" = ?", disbursementData.ID).UpdateColumn("stage", stage)

	disbursementHistoryData := disbursementHistory.DisbursementHistory{StageFrom: disbursementData.Stage, StageTo: stage,
		Remark: disbursementStageInput.Remark, LastDisbursementDate: disbursementStageInput.LastDisbursementDate,
		NextDisbursementDate: disbursementStageInput.NextDisbursementDate}

	if strings.EqualFold(stage, "SUCCESS") {
		disbursementHistoryData = disbursementHistory.DisbursementHistory{StageFrom: disbursementData.Stage, StageTo: stage,
			Remark: disbursementStageInput.Remark, LastDisbursementDate: disbursementStageInput.LastDisbursementDate}
	}
	services.DBCPsql.Table("disbursement_history").Create(&disbursementHistoryData)

	rDisbursementHistoryData := r.RDisbursementHistory{DisbursementId: disbursementData.ID, DisbursementHistoryId: disbursementHistoryData.ID}
	go services.DBCPsql.Table("r_disbursement_history").Create(&rDisbursementHistoryData)
	//r_loan_disbursement

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"stageFrom": disbursementData.Stage,
		"stageTo":   stage,
	})

}

func UpdateDisbursementDate(ctx *iris.Context) {
	query := `with upd as (insert into disbursement_history ("stageFrom","stageTo",remark, "lastDisbursementDate", "nextDisbursementDate", "createdAt", "updatedAt")
	select 'PENDING', 'PENDING', concat('MANUAL ASSIGN FOR DISBURSEMENT ID = ', d.id::varchar), ?, ?, current_timestamp, current_timestamp
	from loan l
	join r_loan_disbursement rld on l.id = rld."loanId"
	join disbursement d on d.id = rld."disbursementId"
	where l.id = ?
	returning id, split_part(remark,' ',7)::int "disbursementId", "nextDisbursementDate"),
	upd2 as (update disbursement set "disbursementDate"= foo."nextDisbursementDate" from ( select "disbursementId", "nextDisbursementDate" from upd ) foo where foo."disbursementId" = disbursement.id)
	insert into r_disbursement_history ("disbursementId", "disbursementHistoryId", "createdAt", "updatedAt") select "disbursementId", id, current_timestamp, current_timestamp
	from upd`

	loanId := ctx.Param("loan_id")
	lastDisbursementDate := ctx.Param("last_disb_date") + " 00:00:00"
	nextDisbursementDate := ctx.Param("next_disb_date") + " 00:00:00"

	services.DBCPsql.Exec(query, lastDisbursementDate, nextDisbursementDate, loanId)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}

// set LWK/UPK date
func SetLWKUPKDate(ctx *iris.Context) {
	payload := struct {
		DateType   string   `json:"dateType"`
		LwkUpkDate string   `json:"lwkUpkDate"`
		BorrowerId []uint64 `json:"borrowerId"`
	}{}
	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// parse string
	layout := "2006-01-02"
	str := payload.LwkUpkDate
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}

	// update borrower
	query := ""
	if payload.DateType == "lwk1" {
		query = `update borrower set "lwk1Date"=? where id in (?)`
	} else if payload.DateType == "lwk2" {
		query = `update borrower set "lwk2Date"=? where id in (?)`
	} else if payload.DateType == "upk" {
		query = `update borrower set "upkDate"=? where id in (?)`
	} else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	services.DBCPsql.Exec(query, t, payload.BorrowerId)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})

}
