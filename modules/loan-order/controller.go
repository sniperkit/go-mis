package loanOrder

import (
	"fmt"
	"strconv"
	"time"

	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/voucher"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&LoanOrder{})
	services.BaseCrudInit(LoanOrder{}, []LoanOrder{})
}

func FetchAll(ctx *iris.Context) {

	query := "select lo.id , c.name, c.username, a.\"totalBalance\",\"orderNo\",sum(l.plafond) as \"totalPlafond\", "
	query += "case when rlov.id is not null then TRUE else FALSE end \"usingVoucher\", "
	query += "case when rlov.id is not null then v.amount else 0 end \"voucherAmount\" "
	query += "from loan l join r_loan_order rlo on l.id = rlo.\"loanId\" "
	query += "join loan_order lo on lo.id = rlo.\"loanOrderId\" "
	query += "join r_investor_product_pricing_loan rippl on rippl.\"loanId\" = l.id "
	query += "join investor i on i.id = rippl.\"investorId\" "
	query += "join r_cif_investor rci on rci.\"investorId\" = i.id "
	query += "join cif c on c.id = rci.\"cifId\" "
	query += "join r_account_investor rai on rai.\"investorId\" = i.id "
	query += "join account a on a.id = rai.\"accountId\" "
	query += "left join r_loan_order_voucher rlov on rlov.\"loanOrderId\" = lo.id "
	query += "left join voucher v on v.id = rlov.\"voucherId\" "
	query += "where lo.remark = 'PENDING' "
	query += "and a.\"deletedAt\" isnull and l.\"deletedAt\" isnull "
	query += "and lo.\"deletedAt\" isnull and c.\"deletedAt\" isnull "
	query += "and i.\"deletedAt\" isnull "
	query += "group by c.name, c.username,a.\"totalBalance\",\"orderNo\",lo.id, rlov.id, v.amount order by lo.id desc"

	loanOrderSchema := []LoanOrderList{}
	e := services.DBCPsql.Raw(query).Scan(&loanOrderSchema).Error
	if e != nil {
		fmt.Println(e)
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loanOrderSchema,
	})
}

func FetchSingle(ctx *iris.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	query := `select i.id, c.username, c.name, lo."orderNo", l.id "loanId", acc."totalBalance", l.plafond, lo.remark
from investor i join r_account_investor rai on i.id = rai."investorId" join account acc on acc.id = rai."accountId"
join r_cif_investor rci on i.id=rci."investorId" join cif c on c.id=rci."cifId"
join r_investor_product_pricing_loan rippl on i.id = rippl."investorId" join loan l on l.id=rippl."loanId"
join r_loan_order rlo on l.id = rlo."loanId" join loan_order lo on lo.id = rlo."loanOrderId"
where lo.remark = 'PENDING' and i.id = ?`

	loanOrderSchema := []LoanOrderDetail{}
	e := services.DBCPsql.Raw(query, id).Scan(&loanOrderSchema).Error
	if e != nil {
		fmt.Println(e)
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loanOrderSchema,
	})
}

// fungsi-fungsi dewa
func AcceptLoanOrder(ctx *iris.Context) {
	// seting order no
	orderNo := ctx.Param("orderNo")
	// get loanid
	loans := GetLoans(orderNo)
	// account
	accountId := GetAccountId(orderNo)
	fmt.Println("susatu")
	fmt.Printf("%v", loans)
	fmt.Printf("%v", accountId)
	fmt.Println("habis")

	// update success

	var voucherAmount float64 = 0.0
	voucherData := voucher.ChekVoucherByOrderNo(orderNo)
	if voucherData != (voucher.Voucher{}) {
		voucherAmount = voucherData.Amount
	}

	totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(accountId)
	totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(accountId)

	totalBalance := (totalDebit + voucherAmount) - totalCredit
	if totalBalance < 0 {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":  "error",
			"message": "totalBalance not enought",
			"data":    iris.Map{},
		})
	} else {
		UpdateSuccess(orderNo)
		CheckVoucherAndInsertToDebit(accountId, orderNo)
		UpdateCredit(loans, accountId)
		UpdateAccount2(orderNo, accountId)
		insertLoanHistoryAndRLoanHistory(orderNo)
		updateLoanStageToInvestor(orderNo)
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	}

}

func UpdateSuccess(orderNo string) {
	query := `update loan_order set remark = 'SUCCESS' where "orderNo" = ?`
	err := services.DBCPsql.Exec(query, orderNo).Error
	if err != nil {
		fmt.Println(err)
	}
}

type LoanId struct {
	ID int64
}

func GetLoans(orderNo string) []int64 {
	query := `select l.id from loan_order lo join r_loan_order rlo on rlo."loanOrderId" = lo.id
	join loan l on l.id = rlo."loanId" where lo."orderNo"=?`

	var L []LoanId
	services.DBCPsql.Raw(query, orderNo).Scan(&L)
	var l []int64
	for _, val := range L {
		l = append(l, val.ID)
	}
	return l
}

func UpdateCredit(loans []int64, accountId uint64) {
	for _, loanId := range loans {

		query := `with ins_1 as (insert into account_transaction_credit ("type","amount","transactionDate","createdAt")
		select 'INVEST', plafond, current_timestamp, current_timestamp from loan l where l.id = 34216 returning id),
		ins_2 as (
			insert into r_account_transaction_credit_loan ("loanId","accountTransactionCreditId","createdAt")
			select ?, ins_1.id,current_timestamp from ins_1 returning "accountTransactionCreditId")
			insert into r_account_transaction_credit ("accountTransactionCreditId","accountId","createdAt")
			select ins_2."accountTransactionCreditId",?, current_timestamp from ins_2`

		services.DBCPsql.Exec(query, loanId, accountId)
	}
}

func UpdateAccount(orderNo string, accountId uint64) {
	query := `with ins as (select SUM(plafond) "total"
	from loan l join r_loan_order rlo on l.id = rlo."loanId"
	join loan_order lo on lo.id = rlo."loanOrderId"
	where lo."orderNo"=?)
	update account set "totalCredit" = "totalCredit"+ins."total", "totalDebit" = "totalDebit"+ins."total"  from ins where account.id = ?`

	services.DBCPsql.Exec(query, orderNo, accountId) // ntar
}

type AccId struct {
	AccountId uint64 `gorm:"column:accountId"`
}

func GetAccountId(orderNo string) uint64 {
	query := `select rai."accountId" from loan_order lo
	join r_loan_order rlo on rlo."loanOrderId" = lo.id
	join r_investor_product_pricing_loan rippl on rippl."loanId" = rlo."loanId"
	join r_account_investor rai on rai."investorId" = rippl."investorId"
	where lo."orderNo"=?`

	var accId AccId
	services.DBCPsql.Raw(query, orderNo).Scan(&accId) // ntar
	return accId.AccountId
}

func UpdateAccount2(orderNo string, accountId uint64) {
	query := `select SUM(plafond) "total"
from loan l join r_loan_order rlo on l.id = rlo."loanId"
join loan_order lo on lo.id = rlo."loanOrderId"
where lo."orderNo"=?`

	r := struct{ Total int64 }{}
	services.DBCPsql.Raw(query, orderNo).Scan(&r)

	query = `update account set "totalCredit" = "totalCredit"+?, "totalBalance" = "totalBalance"-? where account.id = ?`
	services.DBCPsql.Exec(query, r.Total, r.Total, accountId)
}

func insertLoanHistoryAndRLoanHistory(orderNo string) {
	query := `with ins as (INSERT INTO loan_history("stageFrom","stageTo","remark","createdAt","updatedAt")
	select  upper('CART'),upper('INVESTOR'),concat('loan id = ' ,l.id,' updated stage to INVESTOR ', ' orderNo=` + orderNo + `'),current_timestamp,current_timestamp from loan_order lo join r_loan_order rlo on rlo."loanOrderId" = lo.id join loan l on l.id = rlo."loanId" where lo."orderNo"='` + orderNo + `' returning id, (string_to_array(remark,' '))[4]::int as loanId)
	INSERT INTO r_loan_history("loanId","loanHistoryId","createdAt","updatedAt") select  ins.loanId,ins.id ,current_timestamp,current_timestamp from ins`
	services.DBCPsql.Exec(query)
}

func updateLoanStageToInvestor(orderNo string) {
	query := `update loan set stage ='INVESTOR' where id  IN (select l.id from loan_order lo join r_loan_order rlo on rlo."loanOrderId" = lo.id join loan l on l.id = rlo."loanId" where lo."orderNo"='` + orderNo + `')`
	services.DBCPsql.Exec(query)
}

func FetchAllPendingWaiting(ctx *iris.Context) {
	loansOrderPendingWaiting := []LoanOrderInvestorPendingWaiting{}

	query := "select lo.id , c.name, a.\"totalBalance\",\"orderNo\",sum(plafond) as totalPlafond, "
	query += "case when rlov.id is not null then TRUE else FALSE end \"usingVoucher\", "
	query += "case when rlov.id is not null then v.amount else 0 end \"voucherAmount\" "
	query += "from loan l join r_loan_order rlo on l.id = rlo.\"loanId\" "
	query += "join loan_order lo on lo.id = rlo.\"loanOrderId\" "
	query += "join r_investor_product_pricing_loan rippl on rippl.\"loanId\" = l.id "
	query += "join investor i on i.id = rippl.\"investorId\" "
	query += "join r_cif_investor rci on rci.\"investorId\" = i.id "
	query += "join cif c on c.id = rci.\"cifId\" "
	query += "join r_account_investor rai on rai.\"investorId\" = i.id "
	query += "join account a on a.id = rai.\"accountId\" "
	query += "left join r_loan_order_voucher rlov on rlov.\"loanOrderId\" = lo.id "
	query += "left join voucher v on v.id = rlov.\"voucherId\" "
	query += "where lo.remark = 'PENDING' "
	query += "and a.\"deletedAt\" isnull and l.\"deletedAt\" isnull "
	query += "and lo.\"deletedAt\" isnull and c.\"deletedAt\" isnull "
	query += "and i.\"deletedAt\" isnull "
	query += "group by c.name,a.\"totalBalance\",\"orderNo\",lo.id, rlov.id, v.amount order by lo.id desc "

	services.DBCPsql.Raw(query).Find(&loansOrderPendingWaiting)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loansOrderPendingWaiting,
	})
}

func RejectLoanOrder(ctx *iris.Context) {
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

func CheckVoucherAndInsertToDebit(accountID uint64, orderNo string) {
	voucher_data := voucher.ChekVoucherByOrderNo(orderNo)
	if voucher_data != (voucher.Voucher{}) {
		accountTRDebit := accountTransactionDebit.AccountTransactionDebit{Type: "VOUCHER", Amount: voucher_data.Amount, TransactionDate: time.Now()}
		services.DBCPsql.Table("account_transaction_debit").Create(&accountTRDebit)

		r_accountTRDebit := r.RAccountTransactionDebit{AccountId: accountID, AccountTransactionDebitId: accountTRDebit.ID}
		services.DBCPsql.Table("r_account_transaction_debit").Create(&r_accountTRDebit)

	}

}
