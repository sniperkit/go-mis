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

// FetchAll - fetchAll Disbursement data
func FetchAll(ctx *iris.Context) {
	disbursements := []DisbursementFetch{}

	query := "SELECT disbursement.\"id\", disbursement.\"disbursementDate\", disbursement.\"stage\", "
	query += "loan.\"submittedLoanDate\", loan.\"plafond\", \"group\".\"name\" as \"group\", branch.\"name\" as \"branch\", "
	query += "cif.\"name\" as \"borrower\" "
	query += "FROM disbursement "
	query += "JOIN r_loan_disbursement ON r_loan_disbursement.\"disbursementId\" = disbursement.\"id\" "
	query += "JOIN loan ON r_loan_disbursement.\"loanId\" = loan.\"id\" "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.\"id\" "
	query += "JOIN borrower ON borrower.\"id\" = r_loan_borrower.\"borrowerId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\" "
	query += "JOIN cif ON cif.\"id\" = r_cif_borrower.\"cifId\" "
	query += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	query += "LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "
	query += "WHERE disbursement.\"stage\" IN ('PENDING', 'FAILED') "

	services.DBCPsql.Raw(query).Find(&disbursements)
	ctx.JSON(iris.StatusOK, iris.Map{"data": disbursements})
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

	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status":    "success",
		"stageFrom": disbursementData.Stage,
		"stageTo":   stage,
	})

}
