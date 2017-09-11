package loanInsurance

import (
  "gopkg.in/kataras/iris.v4"
  "bitbucket.org/go-mis/modules/account"
  "bitbucket.org/go-mis/modules/loan"
  accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
  accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
  "bitbucket.org/go-mis/modules/r"
  "bitbucket.org/go-mis/services"
  "time"
)

func GetLoanWithInsurance (ctx *iris.Context) {
  query := `SELECT 
    l.id as "loanId",
    cif.name as "borrowerName",
    SUM(CASE WHEN i.type = 'PAR' THEN 1 ELSE 0 END) as "totalPar",
    SUM(CASE WHEN i.type != 'PAR' THEN 1 ELSE 0 END) as "totalOtherType",
    l."isInsurance",
    l."isInsuranceRequested",
    l."isInsuranceRefund"
  FROM loan l
  JOIN r_loan_borrower rlb ON rlb."loanId" = l.id
  JOIN r_cif_borrower rcb ON rcb."borrowerId" = rlb."borrowerId"
  JOIN cif ON cif.id = rcb."cifId"
  JOIN r_loan_installment rli ON rli."loanId" = l.id
  JOIN installment i ON i.id = rli."installmentId"
  WHERE 
  l.stage = 'INSTALLMENT'
  AND l."isInsurance" = TRUE
  GROUP BY l.id, cif.name limit 10`
  
  var loanInsuranceSchema []LoanInsuranceSchema
  
  services.DBCPsql.Raw(query).Find(&loanInsuranceSchema)
  
  ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"totalRows": 0,
		"data": loanInsuranceSchema,
	})
}

func RequestRefund (ctx *iris.Context) {
  id := ctx.Param("loan_id")
  services.DBCPsql.Table("loan").Where("id = ?", id).Update("isInsuranceRequested", "TRUE")
  ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{},
	})
}

type Total struct {
  TotalFrequency uint64 `gorm:"column:totalFrequency"`
}

func ApplyRefund (ctx *iris.Context) {
  loanID := ctx.Param("loan_id")
  
   loanSchema := loan.Loan{}
  services.DBCPsql.Table("loan").Where("id = ?", loanID).Scan(&loanSchema)
  
  tenor := loanSchema.Tenor
  installment := loanSchema.Plafond / float64(loanSchema.Tenor)
  
  query := `select sum(frequency) as "totalFrequency" from installment join r_loan_installment rli on rli."installmentId" = installment.id where rli."loanId"= ?`
  totalSchema := Total{}
  services.DBCPsql.Raw(query, loanID).Scan(&totalSchema)
  
  if totalSchema.TotalFrequency < tenor {
    totalFrequencyRefund := tenor - totalSchema.TotalFrequency
    totalRefund := installment * float64(totalFrequencyRefund)
    
    accountSchema := account.Account{}
    queryGetAccountInvestor := `SELECT * FROM account JOIN r_account_investor rai ON rai."accountId" = account.id JOIN r_investor_product_pricing_loan rippl on rippl."investorId" = rai."investorId" WHERE rippl."loanId" = ?`
    services.DBCPsql.Raw(queryGetAccountInvestor, loanID).Scan(&accountSchema)
    
    accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "REFUND", TransactionDate: time.Now(), Amount: totalRefund, Remark: "Refund using ..."}
    services.DBCPsql.Table("account_transaction_debit").Create(accountTransactionDebitSchema)
    
    rAccountTransactionDebitSchema := &r.RAccountTransactionDebit{AccountId: accountSchema.ID, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
    services.DBCPsql.Table("r_account_transaction_debit").Create(rAccountTransactionDebitSchema)
    
    totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(accountSchema.ID)
    totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(accountSchema.ID)
    totalBalance := totalDebit - totalCredit
    
    services.DBCPsql.Table("account").Where("id = ?", accountSchema.ID).Updates(account.Account{TotalDebit: totalDebit, TotalCredit: totalCredit, TotalBalance: totalBalance})
    
    services.DBCPsql.Table("loan").Where("id = ?", loanID).Update("isInsuranceRefund", "TRUE")
  
    ctx.JSON(iris.StatusOK, iris.Map{
  		"status": "success",
  		"data": iris.Map{},
  	}) 
  } else {
    ctx.JSON(iris.StatusInternalServerError, iris.Map{
  		"status": "error",
  		"message": "Loan has been paid. Refund request failed.",
  		"data": iris.Map{},
  	})
  }
}