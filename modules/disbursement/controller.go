package disbursement

import (
	"strconv"
	"strings"

	disbursementHistory "bitbucket.org/go-mis/modules/disbursement-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
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
	// query += "AND loan.\"submittedLoanDate\" != '1900-01-00'  "
	// query += "AND loan.\"submittedLoanDate\" != '#N/A'  "
	// query += "AND loan.\"submittedLoanDate\" != '' "
	query += "AND to_char(DATE(loan.\"submittedLoanDate\"), 'YYYY') = to_char(DATE(now()), 'YYYY') "
	query += "AND to_char(DATE(disbursement.\"disbursementDate\"), 'YYYY-MM-DD') >= to_char(DATE(now()), 'YYYY-MM-DD') "
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
	query := "SELECT \"group\".id AS \"groupId\", \"group\".\"name\" AS \"groupName\", branch.\"name\" AS \"branchName\", borrower.\"borrowerNo\" cif.\"name\" AS \"borrower\", loan.id AS \"loanId\", loan.plafond, disbursement.\"disbursementDate\"::date, disbursement.stage "
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
	query += "WHERE disbursement.stage IN ('PENDING', 'FAILED') "
	query += "AND loan.\"submittedLoanDate\" IS NOT NULL  "
	// query += "AND loan.\"submittedLoanDate\" != '1900-01-00'   "
	// query += "AND loan.\"submittedLoanDate\" != '#N/A'   "
	// query += "AND loan.\"submittedLoanDate\" != ''  "
	query += "AND to_char(DATE(loan.\"submittedLoanDate\"), 'YYYY') = to_char(DATE(now()), 'YYYY')  "
	query += "AND to_char(DATE(disbursement.\"disbursementDate\"), 'YYYY-MM-DD') >= to_char(DATE(now()), 'YYYY-MM-DD')  "
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
