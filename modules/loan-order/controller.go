package loanOrder

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
	"fmt"
	"strconv"
)

func Init() {
	services.DBCPsql.AutoMigrate(&LoanOrder{})
	services.BaseCrudInit(LoanOrder{}, []LoanOrder{})
}

type LoanOrderCompact struct {
	ID           uint64  `json:"_id"`
	Username     string  `json:"username"`
	Name         string  `json:"name"`
	OrderNo      string  `gorm:"column:orderNo" json:"orderNo"`
	TotalBalance float64 `gorm:"column:totalBalance" json:"totalBalance"`
	TotalPlafond float64 `gorm:"column:totalPlafond" json:"totalPlafond"`
	Remark       string  `json:"remark"`
}


func FetchAll(ctx *iris.Context) {

	query := `select c.username, c.name, i.id, acc."totalBalance", lo."orderNo", sum(l.plafond) "totalPlafond", lo.remark
from investor i join r_account_investor rai on i.id = rai."investorId" join account acc on acc.id = rai."accountId"
join r_cif_investor rci on i.id=rci."investorId" join cif c on c.id=rci."cifId"
join r_investor_product_pricing_loan rippl on i.id = rippl."investorId" join loan l on l.id=rippl."loanId"
join r_loan_order rlo on l.id = rlo."loanId" join loan_order lo on lo.id = rlo."loanOrderId" 
where lo.remark = 'WAITING PAYMENT'
group by c.username, c.name, i.id, acc."totalBalance", lo."orderNo", lo.remark`

	loanOrderSchema := []LoanOrderCompact{}
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


	query := `select c.username, c.name, i.id, acc."totalBalance", lo."orderNo", l.plafond, lo.remark
from investor i join r_account_investor rai on i.id = rai."investorId" join account acc on acc.id = rai."accountId"
join r_cif_investor rci on i.id=rci."investorId" join cif c on c.id=rci."cifId"
join r_investor_product_pricing_loan rippl on i.id = rippl."investorId" join loan l on l.id=rippl."loanId"
join r_loan_order rlo on l.id = rlo."loanId" join loan_order lo on lo.id = rlo."loanOrderId"
where lo.remark = 'WAITING PAYMENT' and i.id = ?`

	loanOrderSchema := []LoanOrderCompact{}

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
	UpdateSuccess(orderNo)
	UpdateCredit(loans, accountId)
	UpdateAccount2(orderNo, accountId)
	
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

func GetLoans(orderNo string) []int64{
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
	
func UpdateCredit(loans []int64, accountId int64) {
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

func UpdateAccount(orderNo string, accountId int64) {
	query := `with ins as (select SUM(plafond) "total"
	from loan l join r_loan_order rlo on l.id = rlo."loanId" 
	join loan_order lo on lo.id = rlo."loanOrderId"
	where lo."orderNo"=?)
	update account set "totalCredit" = "totalCredit"+ins."total", "totalDebit" = "totalDebit"+ins."total"  from ins where account.id = ?`

	services.DBCPsql.Exec(query, orderNo, accountId) // ntar
}


type AccId struct {
	AccountId int64 `gorm:"column:accountId"`
}
func GetAccountId(orderNo string) int64 {
	query := `select rai."accountId" from loan_order lo 
	join r_loan_order rlo on rlo."loanOrderId" = lo.id
	join r_investor_product_pricing_loan rippl on rippl."loanId" = rlo."loanId"
	join r_account_investor rai on rai."investorId" = rippl."investorId"
	where lo."orderNo"=?`

	var accId AccId 
	services.DBCPsql.Raw(query, orderNo).Scan(&accId) // ntar
	return accId.AccountId
}

func UpdateAccount2(orderNo string, accountId int64) {
	query := `select SUM(plafond) "total"
from loan l join r_loan_order rlo on l.id = rlo."loanId" 
join loan_order lo on lo.id = rlo."loanOrderId"
where lo."orderNo"=?`
	
	r := struct { Total int64 }{}
	services.DBCPsql.Raw(query, orderNo).Scan(&r)
	
	query = `update account set "totalCredit" = "totalCredit"+?, "totalBalance" = "totalBalance"-? where account.id = ?`
	services.DBCPsql.Exec(query, r.Total, r.Total, accountId)
}
