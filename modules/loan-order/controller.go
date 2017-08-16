package loanOrder

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/go-mis/modules/account-transaction-credit"
	"bitbucket.org/go-mis/modules/account-transaction-debit"
	"bitbucket.org/go-mis/modules/campaign"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/voucher"
	"bitbucket.org/go-mis/services"
	"github.com/jinzhu/gorm"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&LoanOrder{})
	services.BaseCrudInit(LoanOrder{}, []LoanOrder{})
}

func FetchAll(ctx *iris.Context) {

	query := `select a.threshold as "threshold", lo."createdAt", camp.amount as "campaignAmount", lo.id , c.name, c.username, a."totalBalance","orderNo",sum(l.plafond) as "totalPlafond",
	case when rlov.id is not null then TRUE else FALSE end "usingVoucher",
	case when rlov.id is not null then v.amount else 0 end "voucherAmount",
	case when rloc.id is not null then TRUE else FALSE end "participateCampaign",
	case when rloc.id is not null then rloc.quantity else 0 end "quantityOfCampaignItem",
	case when lo.remark = 'PENDING-REFERRAL' then TRUE else FALSE end "usingRefreal"
	from loan l
	join r_loan_order rlo on l.id = rlo."loanId"
	join loan_order lo on lo.id = rlo."loanOrderId"
	join r_investor_product_pricing_loan rippl on rippl."loanId" = l.id
	join investor i on i.id = rippl."investorId" join r_cif_investor rci on rci."investorId" = i.id
	join cif c on c.id = rci."cifId" join r_account_investor rai on rai."investorId" = i.id
	join account a on a.id = rai."accountId"
	left join r_loan_order_voucher rlov on rlov."loanOrderId" = lo.id
	left join voucher v on v.id = rlov."voucherId"
	left join r_loan_order_campaign rloc on rloc."loanOrderId" = lo.id
	left join campaign camp on rloc."campaignId" = camp.id
	where (lo.remark = 'PENDING' or lo.remark = 'PENDING-REFERRAL' or lo.remark = 'PENDING-FIRST') and a."deletedAt" isnull and l."deletedAt" isnull and lo."deletedAt" isnull and c."deletedAt" isnull and i."deletedAt" isnull
	group by  a.threshold,camp.amount, c.name, c.username, a."totalBalance","orderNo",lo.id, rlov.id, rloc.id, rloc.quantity, v.amount order by lo.id desc`

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

	query := `select acc.threshold as "threshold",i.id, camp.amount as "campaignAmount", c.username, c.name, lo."orderNo", l.id "loanId", l.purpose, acc."totalBalance", l.plafond, lo.remark,
	case when rlov.id is not null then TRUE else FALSE end "usingVoucher", 
	case when rlov.id is not null then v.amount else 0 end "voucherAmount",
	case when rloc.id is not null then TRUE else FALSE end "participateCampaign",
	case when rloc.id is not null then rloc.quantity else 0 end "quantityOfCampaignItem",
	case when lo.remark = 'PENDING-REFERRAL' then TRUE else FALSE end "usingRefreal"
	from investor i
	join r_account_investor rai on i.id = rai."investorId"
	join account acc on acc.id = rai."accountId"
	join r_cif_investor rci on i.id=rci."investorId"
	join cif c on c.id=rci."cifId"
	join r_investor_product_pricing_loan rippl on i.id = rippl."investorId"
	join loan l on l.id=rippl."loanId"
	join r_loan_order rlo on l.id = rlo."loanId"
	join loan_order lo on lo.id = rlo."loanOrderId"
	left join r_loan_order_campaign rloc on rloc."loanOrderId" = lo.id
	left join campaign camp on rloc."campaignId" = camp.id
	left join r_loan_order_voucher rlov on rlov."loanOrderId" = lo.id
	left join voucher v on v.id = rlov."voucherId"
	where (lo.remark = 'PENDING' or lo.remark = 'PENDING-REFERRAL' or lo.remark = 'PENDING-FIRST') and lo.id = ?`

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

func calculateTotalPayment(orderNo string, db *gorm.DB) (float64, error) {
	query := `select SUM(plafond) "total"
	from loan l join r_loan_order rlo on l.id = rlo."loanId"
	join loan_order lo on lo.id = rlo."loanOrderId"
	where lo."orderNo"=?`

	r := struct{ Total float64 }{}
	if err := db.Raw(query, orderNo).Scan(&r).Error; err != nil {
		return 100000000000, err
	}
	return r.Total, nil
}

// fungsi-fungsi dewa
func AcceptLoanOrder(ctx *iris.Context) {
	// seting order no
	orderNo := ctx.Param("orderNo")
	// get loanid
	loans := GetLoans(orderNo)
	// account
	accId := GetAccountId(orderNo)
	// update success
	var voucherAmount float64 = 0.0
	voucherData := voucher.CheckVoucherByOrderNo(orderNo)
	if voucherData != (voucher.Voucher{}) {
		voucherAmount = voucherData.Amount
	}

	db := services.DBCPsql.Begin()

	totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(accId.AccountId)
	totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(accId.AccountId)
	totalOrder, _ := calculateTotalPayment(orderNo, db)
	totalBalance := (totalDebit + voucherAmount) - totalCredit - totalOrder

	if totalBalance < 0 {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":  "error",
			"message": "totalBalance not enought",
			"data":    iris.Map{},
		})
		return
	}

	if err := CheckVoucherAndInsertToDebit(accId.AccountId, orderNo, db); err != nil {
		processErrorAndRollback(ctx, orderNo, db, err, "Check Voucher and Insert into Debit")
		return
	}

	if err := CheckingCampaignAndProgressIntoAccountTransaction(accId.AccountId, orderNo, db); err != nil {
		processErrorAndRollback(ctx, orderNo, db, err, "Check Campaign and Insert into Credit")
		return
	}

	accountTRCredit, errUpdateCredit := UpdateCredit(loans, totalOrder, accId.AccountId, db)
	if errUpdateCredit != nil {
		processErrorAndRollback(ctx, orderNo, db, errUpdateCredit, "Update Credit")
		return
	}

	if err := UpdateAccountCredit(totalOrder, accId.AccountId, db); err != nil {
		processErrorAndRollback(ctx, orderNo, db, err, "Update Account")
		return
	}

	if err := CheckReferalAndEmptytreshold(accId.InvestorId, accountTRCredit.ID, orderNo, db); err != nil {
		processErrorAndRollback(ctx, orderNo, db, err, "Check Referal and Empty Treshold")
		return
	}

	if err := UpdateSuccess(orderNo, db); err != nil {
		processErrorAndRollback(ctx, orderNo, db, err, "Update Success")
		return
	}

	if err := insertLoanHistoryAndRLoanHistory(orderNo, db); err != nil {
		processErrorAndRollback(ctx, orderNo, db, err, "Insert Loan History")
		return
	}
	if err := updateLoanStageToInvestor(orderNo, db); err != nil {
		processErrorAndRollback(ctx, orderNo, db, err, "Update Loan Stage")
		return
	}

	db.Commit()

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   "",
	})

}

func processErrorAndRollback(ctx *iris.Context, orderNo string, db *gorm.DB, err error, process string) {
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": "Error on " + process + " " + err.Error()})
}

func UpdateSuccess(orderNo string, db *gorm.DB) error {
	query := `update loan_order set remark = 'SUCCESS' where "orderNo" = ?`
	return db.Exec(query, orderNo).Error
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

func UpdateCredit(loans []int64, totalPayment float64, accountId uint64, db *gorm.DB) (*accountTransactionCredit.AccountTransactionCredit, error) {
	accountTRCredit := &accountTransactionCredit.AccountTransactionCredit{Type: "INVEST", Amount: totalPayment, TransactionDate: time.Now()}
	db.Table("account_transaction_credit").Create(accountTRCredit)

	//Insert into r_account_transaction_credit
	r_accountTRCredit := &r.RAccountTransactionCredit{AccountId: accountId, AccountTransactionCreditId: accountTRCredit.ID}
	db.Table("r_account_transaction_credit").Create(r_accountTRCredit)

	for _, loanId := range loans {
		insertRatclQuery := "INSERT INTO r_account_transaction_credit_loan(\"loanId\",\"accountTransactionCreditId\", \"createdAt\", \"updatedAt\") VALUES(?,?,current_timestamp,current_timestamp)"
		if err := db.Exec(insertRatclQuery, loanId, accountTRCredit.ID).Error; err != nil {
			return accountTRCredit, err
		}
	}
	return accountTRCredit, nil
}

type AccId struct {
	AccountId  uint64 `gorm:"column:accountId"`
	InvestorId uint64 `gorm:"column:investorId"`
}

func GetAccountId(orderNo string) AccId {
	query := `select rai."investorId",rai."accountId" from loan_order lo
	join r_loan_order rlo on rlo."loanOrderId" = lo.id
	join r_investor_product_pricing_loan rippl on rippl."loanId" = rlo."loanId"
	join r_account_investor rai on rai."investorId" = rippl."investorId"
	where lo."orderNo"=?`

	var accId AccId
	services.DBCPsql.Raw(query, orderNo).Scan(&accId) // ntar
	return accId
}

func UpdateAccountCredit(totalOrder float64, accountId uint64, db *gorm.DB) error {
	query := `update account set "totalCredit" = "totalCredit"+?, "totalBalance" = "totalBalance"-? where account.id = ?`
	return db.Exec(query, totalOrder, totalOrder, accountId).Error
}

type InvestorStruct struct {
	ID uint64 `gorm:"column:investorId"`
}

func insertLoanHistoryAndRLoanHistory(orderNo string, db *gorm.DB) error {

	getInvestorIdQuery := `select "investorId" from r_investor_product_pricing_loan rippl
												join r_loan_order rlo on rlo."loanId" = rippl."loanId"
												join loan_order lo on lo.id = rlo."loanOrderId" where lo."orderNo"='` + orderNo + `'`

	investorStruct := InvestorStruct{}

	db.Raw(getInvestorIdQuery).Scan(&investorStruct)

	query := `with ins as (INSERT INTO loan_history("stageFrom","stageTo","remark","createdAt","updatedAt")
					select  upper('ORDERED'),upper('INVESTOR'),concat('loan id = ' ,l.id,' updated stage to INVESTOR ', ' orderNo = %d investorId %d '),current_timestamp,current_timestamp from loan_order lo join r_loan_order rlo on rlo."loanOrderId" = lo.id join loan l on l.id = rlo."loanId" where lo."orderNo"='` + orderNo + `' returning id, (string_to_array(remark,' '))[4]::int as loanId)
					INSERT INTO r_loan_history("loanId","loanHistoryId","createdAt","updatedAt") select  ins.loanId,ins.id ,current_timestamp,current_timestamp from ins`

	query = fmt.Sprintf(query, orderNo, investorStruct.ID)
	return db.Exec(query).Error
}

func updateLoanStageToInvestor(orderNo string, db *gorm.DB) error {
	query := `update loan set stage ='INVESTOR' where id  IN (select l.id from loan_order lo join r_loan_order rlo on rlo."loanOrderId" = lo.id join loan l on l.id = rlo."loanId" where lo."orderNo"='` + orderNo + `')`
	return db.Exec(query).Error
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

	queryUpdateLoanStage := "update loan set stage = 'PRIVATE' where id in (select l.id from loan l join r_loan_order rlo on l.id = rlo.\"loanId\" join loan_order lo on lo.id = rlo.\"loanOrderId\" where lo.\"orderNo\"='" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateLoanStage)

	queryInsertLoanHistory := `with ins as (INSERT INTO loan_history("stageFrom","stageTo","remark","createdAt","updatedAt")
	select  upper('ORDERED'),upper('PRIVATE'),concat('loan id = ' ,l.id,' updated stage to PRIVATE ', ' orderNo=` + orderNo + `'),current_timestamp,current_timestamp from loan_order lo join r_loan_order rlo on rlo."loanOrderId" = lo.id join loan l on l.id = rlo."loanId" where lo."orderNo"='` + orderNo + `' returning id, (string_to_array(remark,' '))[4]::int as loanId)
	INSERT INTO r_loan_history("loanId","loanHistoryId","createdAt","updatedAt") select  ins.loanId,ins.id ,current_timestamp,current_timestamp from ins`
	services.DBCPsql.Exec(queryInsertLoanHistory)

	queryUpdateRLoanOrderVouher := "update r_loan_order_voucher set \"deletedAt\" = current_timestamp where id in( select rlov.id from r_loan_order_voucher rlov join loan_order lo on lo.id = rlov.\"loanOrderId\" where \"orderNo\" = '" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateRLoanOrderVouher)

	queryUpdateRLoanOrderCampaign := "update r_loan_order_campaign set \"deletedAt\" = current_timestamp where id in( select rloc.id from r_loan_order_campaign rloc join loan_order lo on lo.id = rloc.\"loanOrderId\" where \"orderNo\" = '" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateRLoanOrderCampaign)

	queryUpdateLoanOrderRemark := "update loan_order set remark = 'FAILED' where \"orderNo\" = '" + orderNo + "'"
	services.DBCPsql.Exec(queryUpdateLoanOrderRemark)

	queryUpdateRipplInvestorID := "update r_investor_product_pricing_loan set \"investorId\" = null where \"loanId\" in (select l.id from loan l join r_loan_order rlo on l.id = rlo.\"loanId\" join loan_order lo on lo.id = rlo.\"loanOrderId\" where lo.\"orderNo\"='" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateRipplInvestorID)

	queryUpdateLoanOrderDeleted := "update loan_order set \"deletedAt\" = current_timestamp where \"orderNo\" = '" + orderNo + "';"
	services.DBCPsql.Exec(queryUpdateLoanOrderDeleted)

	queryUpdateRLoanOrderDeleted := "update r_loan_order set \"deletedAt\" = current_timestamp where id in (select rlo.id from r_loan_order rlo join loan_order lo on lo.id = rlo.\"loanOrderId\" where \"orderNo\" = '" + orderNo + "');"
	services.DBCPsql.Exec(queryUpdateRLoanOrderDeleted)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}

func calculateTotalPayment(orderNo string, db *gorm.DB) (float64, error) {
	query := `select SUM(plafond) "total"
	from loan l join r_loan_order rlo on l.id = rlo."loanId"
	join loan_order lo on lo.id = rlo."loanOrderId"
	where lo."orderNo"=?`

	r := struct{ Total float64 }{}
	if err := db.Raw(query, orderNo).Scan(&r).Error; err != nil {
		return 100000000000, err
	}
	return r.Total, nil
}

func CheckReferalAndEmptytreshold(investorId uint64, atcId uint64, orderNo string, db *gorm.DB) error {
	type LoanOrder struct {
		Remark string `json:"remark" gorm:"column:remark"`
	}
	lo := &LoanOrder{}
	queryLo := `select * from loan_order where "orderNo"=?`
	db.Raw(queryLo, orderNo).Scan(lo)

	if lo.Remark == "PENDING-REFERRAL" {
		fmt.Println("Pending REFERAL", lo)
		investorID := strconv.FormatUint(investorId, 10)
		atcID := strconv.FormatUint(atcId, 10)

		queryRefferal := `with A as (
			select id, "inviterInvestorId" from referral where "inviteeInvestorId" = ` + investorID + ` and "inviterGetTimestamp" isnull and "deletedAt" isnull
			),
			B as (
				update referral set "inviterGetTimestamp" = current_timestamp, "inviteeUseTimestamp" = current_timestamp
				from A where A.id = referral.id
				returning referral.id,concat('SUCCESSFUL REFERRAL INVITEE = ',"inviteeInvestorId",' INVITER = ',referral."inviterInvestorId",' REFERALL ID = ',referral.id) "remark",
				referral."inviterInvestorId","inviteeInvestorId","inviterGetTimestamp","inviteeGetTimestamp","createdAt","updatedAt","deletedAt"
			),
			C as (
				insert into account_transaction_debit ("type", amount, remark, "transactionDate","createdAt", "updatedAt")
				select 'REFERRAL-INVITER',100000,"remark",current_timestamp,current_timestamp,current_timestamp from B
				returning id, remark
			),
			D as (
				insert into r_account_transaction_debit_referral ("accountTransactionDebitId","referralId","createdAt")
				select C.id, split_part(remark,' ',12)::int, current_timestamp from C),
			E as (
				insert into r_account_transaction_credit_referral ("accountTransactionCreditId","referralId","createdAt")
				select ` + atcID + `, A.id, current_timestamp from A
			),
			F as(
				update referral set "inviterUseTimestamp" = current_timestamp
				where "inviterInvestorId" = ` + investorID + ` and "inviterUseTimestamp" is null and "inviterGetTimestamp" is not null and "deletedAt" isnull
				returning id
			),
			G as (
				insert into r_account_transaction_credit_referral ("accountTransactionCreditId","referralId","createdAt")
				select ` + atcID + `, F.id, current_timestamp from F
			),
			H as (
				update account set threshold = coalesce(threshold,0) - (bar.amount+foobar.amount), "totalDebit" = "totalDebit" - (bar.amount+foobar.amount), "totalBalance" = "totalBalance" - (bar.amount+foobar.amount)
				from (select id from r_account_investor where "investorId" = ` + investorID + `) foo ,
				(select count(1)*100000 "amount" from B) bar,
				(select count(1)*100000 "amount" from F) foobar
				where account.id = foo.id
				returning account.id),
			I as (
				insert into r_account_transaction_debit ("accountTransactionDebitId","accountId","createdAt")
				select C.id, rai."accountId", current_timestamp from C
				join r_account_investor rai on split_part(C.remark,' ',8)::int = rai."investorId"
				),
			J as(
				update account set threshold = coalesce(threshold,0) + total * 100000, "totalDebit" = coalesce("totalDebit",0) + total * 100000, "totalBalance" = coalesce("totalBalance",0) + total * 100000
				from (select rai."accountId" id, count(1) "total" from r_account_investor rai join B on rai."investorId" = B."inviterInvestorId" group by "accountId" ) foo where account.id = foo.id
				returning account.id
			)
			select * from B`
		return db.Exec(queryRefferal).Error
	}
	if lo.Remark == "PENDING-FIRST" {
		fmt.Println("Pending First", lo)
		investorID := strconv.FormatUint(investorId, 10)

		queryRefferal := `with A as (
        select id, "inviterInvestorId" from referral where "inviteeInvestorId" = ` + investorID + ` and "inviterGetTimestamp" isnull and "deletedAt" isnull
    ),
    B as (
        update referral set "inviterGetTimestamp" = current_timestamp
        from A where A.id = referral.id
        returning referral.id,concat('SUCCESSFUL REFERRAL INVITEE = ',"inviteeInvestorId",' INVITER = ',referral."inviterInvestorId",' REFERALL ID = ',referral.id) "remark",
        referral."inviterInvestorId","inviteeInvestorId","inviterGetTimestamp","inviteeGetTimestamp","createdAt","updatedAt","deletedAt"
    ),
    C as (
        insert into account_transaction_debit ("type", amount, remark, "transactionDate","createdAt", "updatedAt")
        select 'REFERRAL-INVITER',100000,"remark",current_timestamp,current_timestamp,current_timestamp from B
        returning id, remark
    ),
    D as (
        insert into r_account_transaction_debit_referral ("accountTransactionDebitId","referralId","createdAt")
        select C.id, split_part(remark,' ',12)::int, current_timestamp from C),
    E as (
        insert into r_account_transaction_debit ("accountTransactionDebitId","accountId","createdAt")
        select C.id, rai."accountId", current_timestamp from C
        join r_account_investor rai on split_part(C.remark,' ',8)::int = rai."investorId"
        ),
    F as(
        update account set threshold = coalesce(threshold,0) + total * 100000, "totalDebit" = coalesce("totalDebit",0) + total * 100000, "totalBalance" = coalesce("totalBalance",0) + total * 100000
        from (select rai."accountId" id, count(1) "total" from r_account_investor rai join B on rai."investorId" = B."inviterInvestorId" group by "accountId" ) foo where account.id = foo.id
        returning account.id
    )
    select * from B`
		return db.Exec(queryRefferal).Error
	}
	fmt.Println("Not Pending REFERAL", lo)
	return nil

}

func CheckVoucherAndInsertToDebit(accountID uint64, orderNo string, db *gorm.DB) error {
	voucher_data := voucher.CheckVoucherByOrderNo(orderNo)
	if voucher_data == (voucher.Voucher{}) {
		return nil
	}

	accountTRDebit := accountTransactionDebit.AccountTransactionDebit{Type: "VOUCHER", Amount: voucher_data.Amount, Remark: voucher_data.VoucherNo, TransactionDate: time.Now()}
	if err := db.Table("account_transaction_debit").Create(&accountTRDebit).Error; err != nil {
		return err
	}

	r_accountTRDebit := r.RAccountTransactionDebit{AccountId: accountID, AccountTransactionDebitId: accountTRDebit.ID}
	if err := db.Table("r_account_transaction_debit").Create(&r_accountTRDebit).Error; err != nil {
		return err
	}

	query := `update account set "totalDebit" = "totalDebit"+?, "totalBalance" = "totalBalance"+? where account.id = ?`
	return db.Exec(query, voucher_data.Amount, voucher_data.Amount, accountID).Error

}

func CheckingCampaignAndProgressIntoAccountTransaction(accountID uint64, orderNo string, db *gorm.DB) error {

	quantityOfCampaignItem, campaignData := campaign.GetActiveCampaignByOrderNo(orderNo)
	if campaignData == (campaign.Campaign{}) {
		return nil
	}
	var campaignAmount float64 = float64(campaignData.Amount * quantityOfCampaignItem)

	atc := accountTransactionCredit.AccountTransactionCredit{Type: "CAMPAIGN", Amount: campaignAmount, TransactionDate: time.Now(), Remark: "1KMSAJADAH"}
	if err := db.Table("account_transaction_credit").Create(&atc).Error; err != nil {
		return err
	}

	rAccountTransactionCredit := r.RAccountTransactionCredit{AccountId: accountID, AccountTransactionCreditId: atc.ID}
	if err := db.Table("r_account_transaction_credit").Create(&rAccountTransactionCredit).Error; err != nil {
		return err
	}

	query := `update account set "totalCredit" = "totalCredit"+?, "totalBalance" = "totalBalance"-? where account.id = ?`
	return db.Exec(query, campaignAmount, campaignAmount, accountID).Error

}
