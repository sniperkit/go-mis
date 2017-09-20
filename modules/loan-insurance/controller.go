package loanInsurance

import (
  "gopkg.in/kataras/iris.v4"
  "bitbucket.org/go-mis/modules/account"
  "bitbucket.org/go-mis/modules/loan"
  "bitbucket.org/go-mis/modules/loan-history"
  accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
  accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
  "bitbucket.org/go-mis/modules/r"
  "bitbucket.org/go-mis/services"
  "time"
  "strconv"
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
  LEFT JOIN r_loan_installment rli ON rli."loanId" = l.id
  LEFT JOIN installment i ON i.id = rli."installmentId"
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
  services.DBCPsql.Table("loan").Where("id = ?", id).Update(map[string]interface{}{"isInsuranceRequested": true, "stage": "END"})
  
  loanHistorySchema := &loanHistory.LoanHistory{StageFrom: "INSTALLMENT", StageTo: "END", Remark: "Investor requested for REFUND.", CreatedAt: time.Now(), UpdatedAt: time.Now()}
  services.DBCPsql.Table("loan_history").Create(loanHistorySchema);
  
  loanID, _ := strconv.ParseUint(id, 10, 64)
  
  rLoanHistorySchema := &r.RLoanHistory{LoanId: loanID, LoanHistoryId: loanHistorySchema.ID}
  services.DBCPsql.Table("r_loan_history").Create(rLoanHistorySchema);
  
  ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{},
	})
}

type Total struct {
  TotalFrequency uint64 `gorm:"column:totalFrequency"`
}

const refundRate float64 = 0.75

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
    totalRefund := (installment * float64(totalFrequencyRefund)) * refundRate
    
    accountSchema := account.Account{}
    queryGetAccountInvestor := `SELECT * FROM account JOIN r_account_investor rai ON rai."accountId" = account.id JOIN r_investor_product_pricing_loan rippl on rippl."investorId" = rai."investorId" WHERE rippl."loanId" = ?`
    services.DBCPsql.Raw(queryGetAccountInvestor, loanID).Scan(&accountSchema)
    
    accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "REFUND-INSURANCE", TransactionDate: time.Now(), Amount: totalRefund, Remark: "Refund via JAMKRINDO"}
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

type InsuranceFinanceReport struct {
  TransactionDate *time.Time `json:"transactionDate" gorm:"column:transactionDate"`
  TransactionID uint64 `json:"transactionId" gorm:"column:transactionId"`
  TransactionType string `json:"type" gorm:"column:type"`
  Amount float64 `json:"amount" gorm:"column:amount"`
  LoanID string `json:"loanId" gorm:"column:loanId"`
}

func GetFinanceReport (ctx *iris.Context) {
  query := `SELECT atc."transactionDate", atc.id as "transactionId", atc.type, atc.amount, 
    array_to_string(array_agg(l.id), ',') as "loanId"
    FROM loan l
    JOIN r_account_transaction_credit_loan ratcl ON ratcl."loanId" = l.id
    JOIN account_transaction_credit atc ON atc.id = ratcl."accountTransactionCreditId"
    WHERE 
    l."isInsurance" = TRUE 
    AND atc."type" = 'INSURANCE'
    GROUP BY atc.id`
  
  insuranceFinanceReportSchema := []InsuranceFinanceReport{}
  services.DBCPsql.Raw(query).Scan(&insuranceFinanceReportSchema)
  
  ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
		  "transactions": insuranceFinanceReportSchema,
		},
	})
}