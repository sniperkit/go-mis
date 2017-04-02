package loanOrder

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&LoanOrder{})
	services.BaseCrudInit(LoanOrder{}, []LoanOrder{})
}

func FetchAllPendingWaiting(ctx *iris.Context) {
	loansOrderPendingWaiting := []LoanOrderInvestorPendingWaiting{}

	query := "select lo.id , c.name, a.\"totalBalance\",\"orderNo\",sum(plafond) as totalPlafond from loan l join r_loan_order rlo on l.id = rlo.\"loanId\" "
	query += "join loan_order lo on lo.id = rlo.\"loanOrderId\" "
	query += "join r_investor_product_pricing_loan rippl on rippl.\"loanId\" = l.id "
	query += "join investor i on i.id = rippl.\"investorId\" "
	query += "join r_cif_investor rci on rci.\"investorId\" = i.id "
	query += "join cif c on c.id = rci.\"cifId\" "
	query += "join r_account_investor rai on rai.\"investorId\" = i.id "
	query += "join account a on a.id = rai.\"accountId\" "
	query += "where lo.remark = 'PENDING' or lo.remark = 'WAITING PAYMENT' "
	query += "and a.\"deletedAt\" isnull and l.\"deletedAt\" isnull "
	query += "and lo.\"deletedAt\" isnull and c.\"deletedAt\" isnull "
	query += "and i.\"deletedAt\" isnull "
	query += "group by c.name,a.\"totalBalance\",\"orderNo\",lo.id order by lo.id desc "

	services.DBCPsql.Raw(query).Find(&loansOrderPendingWaiting)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loansOrderPendingWaiting,
	})
}

func Reject(ctx *iris.Context) {
	orderNo := ctx.Param("orderNo")

	queryUpdateLoanStage := "update loan set stage = 'MARKETPLACE' where id in (select l.id from loan l join r_loan_order rlo on l.id = rlo.\"loanId\" join loan_order lo on lo.id = rlo.\"loanOrderId\" where lo.\"orderNo\"='" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateLoanStage)

	queryUpdateVouher := "update voucher set \"deletedAt\" = current_timestamp where id in (select v.id from voucher v join r_loan_order_voucher rlov on rlov.\"voucherId\" = v.id  join loan_order lo on lo.id = rlov.\"loanOrderId\" where \"orderNo\" = '" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateVouher)
	queryUpdateRLoanOrderVouher := "update r_loan_order_voucher set \"deletedAt\" = current_timestamp where id in( select rlov.id from r_loan_order_voucher rlov join loan_order lo on lo.id = rlov.\"loanOrderId\" where \"orderNo\" = '" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateRLoanOrderVouher)

	queryUpdateLoanOrderRemark := "update loan_order set remark = 'FAILED' where \"orderNo\" = '" + orderNo + "'"
	services.DBCPsql.Exec(queryUpdateLoanOrderRemark)

	queryUpdateLoanOrderDeleted := "update loan_order set \"deletedAt\" = current_timestamp where \"orderNo\" = '" + orderNo + "';"
	services.DBCPsql.Exec(queryUpdateLoanOrderDeleted)

	queryUpdateRLoanOrderDeleted := "update r_loan_order set \"deletedAt\" = current_timestamp where id in (select rlo.id from r_loan_order rlo join loan_order lo on lo.id = rlo.\"loanOrderId\" where \"orderNo\" = '" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateRLoanOrderDeleted)

	queryUpdateRipplInvestorID := "update r_investor_product_pricing_loan set \"investorId\" = null where \"loanId\" in (select l.id from loan l join r_loan_order rlo on l.id = rlo.\"loanId\" join loan_order lo on lo.id = rlo.\"loanOrderId\" where lo.\"orderNo\"='" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateRipplInvestorID)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}

func Accept(ctx *iris.Context) {
	orderNo := ctx.Param("orderNo")

	queryUpdateLoanOrderRemark := "update loan_order set remark = 'SUCCESS' where \"orderNo\" = '" + orderNo + "'"
	services.DBCPsql.Exec(queryUpdateLoanOrderRemark)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}
