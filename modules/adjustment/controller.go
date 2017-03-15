package adjustment

import (
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
