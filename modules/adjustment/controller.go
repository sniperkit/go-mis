package adjustment

import (
	"strconv"
	"time"

	account_transaction_debit "bitbucket.org/go-mis/modules/account-transaction-debit"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

type ParamAdjustment struct {
	AccountTransactionDebitID uint64  `json:"accountTransactionDebitId"`
	AmountToAdjust            float64 `json:"amountToAdjust"`
	Remark                    string  `json:"remark"`
}

func Init() {
	services.DBCPsql.AutoMigrate(&Adjustment{})
	services.BaseCrudInit(Adjustment{}, []Adjustment{})
}

// SubmitAdjustment - submit adjustment
func SubmitAdjustment(ctx *iris.Context) {
	userMis := ctx.Get("USER_MIS").(userMis.UserMis)
	paramAdjustment := ParamAdjustment{}

	if err := ctx.ReadJSON(&paramAdjustment); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Todo: add account_type (debit or credit) conditional
	accountTransactionDebit := account_transaction_debit.AccountTransactionDebit{}
	services.DBCPsql.Table("account_transaction_debit").Where("id = ?", paramAdjustment.AccountTransactionDebitID).First(&accountTransactionDebit)

	adjustmentSchema := &Adjustment{
		Type:           accountTransactionDebit.Type,
		AmountBefore:   accountTransactionDebit.Amount,
		AmountToAdjust: paramAdjustment.AmountToAdjust,
		AmountAfter:    accountTransactionDebit.Amount + paramAdjustment.AmountToAdjust,
		Remark:         paramAdjustment.Remark,
	}
	services.DBCPsql.Create(adjustmentSchema)

	rAdjustmentSubmittedBy := &r.RAdjustmentSubmittedBy{
		AdjustmentId: adjustmentSchema.ID,
		UserMisId:    userMis.ID,
	}
	services.DBCPsql.Create(rAdjustmentSubmittedBy)

	rAdjustmentAccountTransactionDebit := &r.RAdjustmentAccountTransactionDebit{
		AdjustmentID:              adjustmentSchema.ID,
		AccountTransactionDebitID: paramAdjustment.AccountTransactionDebitID,
	}
	services.DBCPsql.Create(rAdjustmentAccountTransactionDebit)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
	})
}

type AdjustmentSchema struct {
	ID             uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name           string    `gorm:"column:name" json:"name"`
	Type           string    `gorm:"column:type" json:"type"` // INSTALLMENT, DISBURSEMENT, ACCOUNT-CREDIT, ACCOUNT-DEBIT
	AmountBefore   float64   `gorm:"column:amountBefore" json:"amountBefore"`
	AmountToAdjust float64   `gorm:"column:amountToAdjust" json:"amountToAdjust"`
	AmountAfter    float64   `gorm:"column:amountAfter" json:"amountAfter"`
	Remark         string    `gorm:"column:remark" json:"remark"`
	Stage          string    `gorm:"column:stage" json:"stage"`
	CreatedAt      time.Time `gorm:"column:createdAt" json:"createdAt"`
}

// GetAdjustment - get list of adjustment
func GetAdjustment(ctx *iris.Context) {
	// adjustmentSchema := []Adjustment{}
	// services.DBCPsql.Table("adjustment").Where("\"deletedAt\" IS NULL").Scan(&adjustmentSchema)
	query := "SELECT cif.\"name\", adjustment.*  "
	query += "FROM adjustment "
	query += "JOIN r_installment_adjustment ON r_installment_adjustment.\"adjustmentId\" = adjustment.id "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = r_installment_adjustment.\"installmentId\" "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = r_loan_installment.\"loanId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = r_loan_borrower.\"borrowerId\" "
	query += "JOIN cif ON cif.id = r_cif_borrower.\"cifId\" WHERE adjustment.\"deletedAt\" IS NULL "

	adjustmentSchema := []AdjustmentSchema{}
	services.DBCPsql.Raw(query).Scan(&adjustmentSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   adjustmentSchema,
	})
}

// GetAdjustmentDetail - get adjustment detail
func GetAdjustmentDetail(ctx *iris.Context) {
	adjustmentID := ctx.Param("adjustment_id")

	query := "SELECT cif.\"name\", adjustment.*  "
	query += "FROM adjustment "
	query += "JOIN r_installment_adjustment ON r_installment_adjustment.\"adjustmentId\" = adjustment.id "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = r_installment_adjustment.\"installmentId\" "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = r_loan_installment.\"loanId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = r_loan_borrower.\"borrowerId\" "
	query += "JOIN cif ON cif.id = r_cif_borrower.\"cifId\" WHERE adjustment.\"deletedAt\" IS NULL AND adjustment.id = ?"

	adjustmentSchema := AdjustmentSchema{}
	services.DBCPsql.Raw(query, adjustmentID).Scan(&adjustmentSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   adjustmentSchema,
	})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

type InReviewInstallment struct {
	BorrowerID      uint64    `gorm:"column:borrowerId" json:"borrowerId"`
	Borrower        string    `gorm:"column:borrower" json:"borrower"`
	InstallmentID   uint64    `gorm:"column:installmentId" json:"installmentId"`
	Type            string    `gorm:"column:type" json:"type"`
	Presence        string    `gorm:"column:presence" json:"presence"`
	PaidInstallment float64   `gorm:"column:paidInstallment" json:"paidInstallment"`
	Reserve         float64   `gorm:"column:reserve" json:"reserve"`
	Frequency       int32     `gorm:"column:frequency" json:"frequency"`
	CreatedAt       time.Time `gorm:"column:createdAt" json:"createdAt"`
}

// GetInReviewInstallment - get list of in-review installment
func GetInReviewInstallment(ctx *iris.Context) {
	totalData := TotalData{}

	query := "SELECT r_cif_borrower.\"borrowerId\", cif.\"name\" AS \"borrower\", installment.id AS \"installmentId\", installment.\"type\", installment.presence, installment.\"paidInstallment\", installment.reserve, installment.frequency, installment.\"createdAt\" "
	query += "FROM installment "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.id "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = r_loan_installment.\"loanId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = r_loan_borrower.\"borrowerId\" "
	query += "JOIN cif ON cif.id = r_cif_borrower.\"cifId\" WHERE installment.\"deletedAt\" IS NULL AND installment.stage = 'IN-REVIEW' AND (installment.\"createdAt\" BETWEEN ? AND ?) "

	queryTotal := "SELECT COUNT(*) AS \"totalRows\" "
	queryTotal += "FROM installment "
	queryTotal += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.id "
	queryTotal += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = r_loan_installment.\"loanId\" "
	queryTotal += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = r_loan_borrower.\"borrowerId\" "
	queryTotal += "JOIN cif ON cif.id = r_cif_borrower.\"cifId\" WHERE installment.\"deletedAt\" IS NULL AND installment.stage = 'IN-REVIEW' AND (installment.\"createdAt\" BETWEEN ? AND ?) "

	if ctx.URLParam("search") != "" {
		query += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
		queryTotal += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	if ctx.URLParam("LIMIT") != "" {
		query += "LIMIT " + ctx.URLParam("LIMIT")
		queryTotal += "LIMIT " + ctx.URLParam("LIMIT")
	} else {
		query += "LIMIT 500 "
		queryTotal += "LIMIT 500 "
	}

	startDate := ctx.Param("start_date") + " 00:00:00"
	endDate := ctx.Param("end_date") + " 00:00:00"

	installmentSchema := []InReviewInstallment{}

	services.DBCPsql.Raw(queryTotal, startDate, endDate).Find(&totalData)
	services.DBCPsql.Raw(query, startDate, endDate).Scan(&installmentSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      installmentSchema,
	})

}

type paramSubmitAdjustment struct {
	AmountBefore   float64 `json:"previousAmount"`
	AmountToAdjust float64 `json:"amountToAdjust"`
	AmountAfter    float64 `json:"resultAmount"`
	Remark         string  `json:"remark"`
}

//SetAdjustmentForInstallment
func SetAdjustmentForInstallment(ctx *iris.Context) {
	installmentID, _ := strconv.ParseUint(ctx.Param("installment_id"), 10, 64)

	psaSchema := paramSubmitAdjustment{}
	if err := ctx.ReadJSON(&psaSchema); err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "Bad request param.",
		})
		return
	}

	adjustmentSchema := &Adjustment{
		Type:           "INSTALLMENT",
		AmountBefore:   psaSchema.AmountBefore,
		AmountToAdjust: psaSchema.AmountToAdjust,
		AmountAfter:    psaSchema.AmountAfter,
		Remark:         psaSchema.Remark,
		Stage:          "PENDING",
	}

	services.DBCPsql.Table("adjustment").Create(adjustmentSchema)

	rInstallmentAdjustmentSchema := &r.RInstallmentAdjustment{
		AdjustmentID:  adjustmentSchema.ID,
		InstallmentID: installmentID,
	}

	services.DBCPsql.Table("r_installment_adjustment").Create(rInstallmentAdjustmentSchema)

	userMis := ctx.Get("USER_MIS").(userMis.UserMis)
	rAdjustmentSubmittedBySchema := &r.RAdjustmentSubmittedBy{AdjustmentId: adjustmentSchema.ID, UserMisId: userMis.ID}
	services.DBCPsql.Table("r_adjustment_submitted_by").Create(rAdjustmentSubmittedBySchema)

	// services.DBCPsql.Table("installment").Where("id = ?", installmentID).Update("paidInstallment", psaSchema.AmountAfter)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}

type UpdateAdjustmentInstallmentSchema struct {
	AmountAfter   float64 `gorm:"column:amountAfter" json:"amountAfter"`
	InstallmentID uint64  `gorm:"column:installmentId" json:"installmentId"`
}

// UpdateAdjustmentAndInstallment - Update adjustment stage and update installment amount
func UpdateAdjustmentAndInstallment(ctx *iris.Context) {
	userMis := ctx.Get("USER_MIS").(userMis.UserMis)
	adjustmentID, _ := strconv.ParseUint(ctx.Param("adjustment_id"), 10, 64)
	stage := ctx.URLParam("stage")

	query := "SELECT adjustment.*, r_installment_adjustment.\"installmentId\" FROM adjustment "
	query += "JOIN r_installment_adjustment ON r_installment_adjustment.\"adjustmentId\" = adjustment.id "
	query += "WHERE adjustment.id = ? AND adjustment.\"deletedAt\" IS NULL LIMIT 1"

	adjustmentSchema := UpdateAdjustmentInstallmentSchema{}
	services.DBCPsql.Raw(query, adjustmentID).Scan(&adjustmentSchema)

	if stage == "APPROVE" {
		services.DBCPsql.Table("installment").Where("id = ?", adjustmentSchema.InstallmentID).Update("paidInstallment", adjustmentSchema.AmountAfter)
	} else {
		services.DBCPsql.Table("r_installment_adjustment").Where("\"adjustmentId\" = ?", adjustmentID).Update("deletedAt", time.Now())
	}

	rAdjustmentApprovedBySchema := &r.RAdjustmentApprovedBy{AdjustmentId: adjustmentID, UserMisId: userMis.ID}
	services.DBCPsql.Table("r_adjustment_approved_by").Create(rAdjustmentApprovedBySchema)

	services.DBCPsql.Table("adjustment").Where("id = ?", adjustmentID).Update("stage", stage)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}
