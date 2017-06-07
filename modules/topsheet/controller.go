package topsheet

import (
	"fmt"

	"bitbucket.org/go-mis/modules/installment"
	"bitbucket.org/go-mis/services"
	"bitbucket.org/go-worker-topsheet/model"
	iris "gopkg.in/kataras/iris.v4"
)

type Reserve struct {
	LoanID  uint64  `gorm:"column:loanId"`
	Reserve float64 `gorm:"column:reserve"`
}

func addReserveInfo(currentTopsheet []CurrentTopsheetSchema) []CurrentTopsheetSchema {
	length := len(currentTopsheet)
	installmentIds := make([]uint64, length)
	for i, current := range currentTopsheet {
		installmentIds[i] = current.LatestInstallmentID
	}

	query := `
select loan.id as "loanId", installment.reserve
from installment
join r_loan_installment on r_loan_installment."installmentId" = installment.id
join loan on loan.id = r_loan_installment."loanId"
join r_loan_group on r_loan_group."loanId" = loan.id
join r_group_agent on r_group_agent."groupId" = r_loan_group."groupId"
where installment.id in (?) `

	reserveList := []Reserve{}
	services.DBCPsql.Raw(query, installmentIds).Scan(&reserveList)

	for i, current := range currentTopsheet {
		tempLoanID := current.LoanID
		for _, currentReserve := range reserveList {
			if tempLoanID == currentReserve.LoanID {
				currentTopsheet[i].LatestReserve = currentReserve.Reserve
			}
		}
	}

	return currentTopsheet
}

// GenerateTopsheet is a function to generate topsheet
func GenerateTopsheet(ctx *iris.Context) {
	query := `SELECT agent."fullname" as "agentName", "group"."id" as "groupId", "group"."name" as "groupName", "group"."scheduleDay", "group"."scheduleTime", r_group_branch."branchId",
		loan."id" as "loanId", loan."creditScoreGrade", loan."creditScoreValue", loan."tenor", loan."rate", loan."installment", loan."plafond", loan."subgroup", loan."submittedLoanDate",
		borrower."borrowerNo" as "borrowerNo", cif."name" as "borrowerName", disbursement."disbursementDate",
		SUM(installment."frequency") as "frequency",
		SUM(installment."reserve") as "totalReserve",
		MAX(installment.id) as "latestInstallmentId",
		COUNT(CASE WHEN (installment."presence" = 'HADIR') THEN 1 END) as "totalHadir",
		COUNT(CASE WHEN (installment."presence" = 'ALFA') THEN 1 END) as "totalAlfa",
		COUNT(CASE WHEN (installment."presence" = 'CUTI') THEN 1 END) as "totalCuti",
		COUNT(CASE WHEN (installment."presence" = 'SAKIT') THEN 1 END) as "totalSakit",
		COUNT(CASE WHEN (installment."presence" = 'IZIN') THEN 1 END) as "totalIzin"
		FROM "group" INNER JOIN r_loan_group ON r_loan_group."groupId" = "group"."id"
		INNER JOIN loan ON "loan"."id" = r_loan_group."loanId"
		INNER JOIN r_group_agent ON r_group_agent."groupId" = "group"."id"
		INNER JOIN r_group_branch ON r_group_branch."groupId" = "group"."id"
		INNER JOIN agent ON agent."id" = r_group_agent."agentId"
		INNER JOIN r_loan_borrower ON r_loan_borrower."loanId" = loan."id"
		INNER JOIN borrower ON borrower."id" = r_loan_borrower."borrowerId"
		INNER JOIN r_cif_borrower ON r_cif_borrower."borrowerId" = borrower."id"
		INNER JOIN cif ON cif."id" = r_cif_borrower."cifId"
		LEFT JOIN r_loan_installment ON r_loan_installment."loanId" = loan."id"
		LEFT JOIN installment ON installment."id" = r_loan_installment."installmentId" and installment.stage = 'SUCCESS'
		INNER JOIN r_loan_disbursement ON r_loan_disbursement."loanId" = loan."id"
		INNER JOIN disbursement ON disbursement."id" = r_loan_disbursement."disbursementId"
		WHERE loan."deletedAt" IS NULL 
		AND loan.stage = 'INSTALLMENT'
		AND "group".id = ?
		GROUP BY loan."id", agent."id", agent."fullname", "group"."id", cif."name", borrower."borrowerNo", disbursement."disbursementDate", r_group_branch."branchId"
		ORDER BY "group"."name", loan.subgroup, cif."name" ASC
	`

	topsheetSchema := []CurrentTopsheetSchema{}
	services.DBCPsql.Raw(query, ctx.Param("group_id")).Scan(&topsheetSchema)

	topsheetSchema = addReserveInfo(topsheetSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   topsheetSchema,
	})
}

type JSONTopsheet struct {
	Topsheet []TopsheetFormSchema `json:topsheet`
}

// SubmitTopsheet is a function to insert topsheet
func SubmitTopsheet(ctx *iris.Context) {
	jsonTopsheet := JSONTopsheet{}

	fmt.Println(jsonTopsheet)

	if err := ctx.ReadJSON(&jsonTopsheet); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	fmt.Println("masuk sini")

	topsheetForm := jsonTopsheet.Topsheet

	for _, item := range topsheetForm {
		installmentSchema := &installment.Installment{
			Type:            item.Type,
			Presence:        item.Presence,
			PaidInstallment: item.PaidInstallment,
			Penalty:         item.Penalty,
			Reserve:         item.Reserve,
			Frequency:       item.Frequency,
			Stage:           item.Stage,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.CreatedAt,
			TransactionDate: &item.CreatedAt,
		}

		services.DBCPsql.Table("installment").Create(installmentSchema)
		services.DBCPsql.Table("installment").Where("id = ?", installmentSchema.ID).UpdateColumns(map[string]interface{}{"createdAt": item.CreatedAt, "updatedAt": item.CreatedAt})

		StoreRLoanInstallment(item.LoanID, installmentSchema.ID)
		StoreToInstallmentHistory(installmentSchema.ID, "PENDING", "PENDING")
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   topsheetForm,
	})
}

// StoreRLoanInstallment - Store relation between loan and installment
func StoreRLoanInstallment(loanID uint64, installmentID uint64) {
	rLoanInstallmentSchema := &model.RLoanInstallmentSchema{LoanID: loanID, InstallmentID: installmentID}
	services.DBCPsql.Table("r_loan_installment").Create(rLoanInstallmentSchema)
}

// StoreToInstallmentHistory - Store new installmet history
func StoreToInstallmentHistory(installmentID uint64, stageFrom string, stageTo string) {
	installmentHistorySchema := &model.InstallmentHistorySchema{StageFrom: stageFrom, StageTo: stageTo}
	services.DBCPsql.Table("installment_history").Create(installmentHistorySchema)

	StoreToRInstallmentHistory(installmentID, installmentHistorySchema.ID)
}

// StoreToRInstallmentHistory - Store installment id and also installment history Id
func StoreToRInstallmentHistory(installmentID uint64, installmentHistoryID uint64) {
	rInstallmentHistorySchema := &model.RInstallmentHistorySchema{InstallmentID: installmentID, InstallmentHistoryID: installmentHistoryID}
	services.DBCPsql.Table("r_installment_history").Create(rInstallmentHistorySchema)
}
