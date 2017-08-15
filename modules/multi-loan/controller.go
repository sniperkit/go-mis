package multiloan

import (
  "gopkg.in/kataras/iris.v4"
  "bitbucket.org/go-mis/services"
)

type MultiLoanSchema struct {
  LoanID uint32 `gorm:"column:loanId" json:"loanId"`
  Borrower string `gorm:"column:borrower" json:"borrower"`
  BorrowerNo string `gorm:"column:borrowerNo" json:"borrowerNo"`
  Group string `gorm:"column:group" json:"group"`
  SubmittedLoanDate string `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
  DisbursementDate string `gorm:"column:disbursementDate" json:"disbursementDate"`
  Tenor uint64 `gorm:"column:tenor" json:"tenor"`
	Rate float64 `gorm:"column:rate" json:"rate"`
	Installment float64 `gorm:"column:installment" json:"installment"`
	Plafond float64 `gorm:"column:plafond" json:"plafond"`
	Stage string `gorm:"column:stage" json:"stage"`
}

func GetAllUndisbursedMultiLoan(ctx *iris.Context) {
  multiLoanSchema := []MultiLoanSchema{}
  query := `SELECT loan.id as "loanId", cif.name AS "borrower", borrower."borrowerNo", "group"."name" AS "group", loan."submittedLoanDate", disbursement."disbursementDate", loan.plafond, loan.tenor, loan.rate, loan.stage, loan.installment FROM loan 
  JOIN r_loan_group ON r_loan_group."loanId" = loan.id 
  JOIN r_loan_branch ON r_loan_branch."loanId" = loan.id 
  JOIN r_loan_borrower ON r_loan_borrower."loanId" = loan.id 
  JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan."loanId" = loan.id 
  JOIN borrower ON r_loan_borrower."borrowerId" = borrower.id 
  JOIN r_loan_disbursement ON r_loan_disbursement."loanId" = loan.id 
  JOIN r_cif_borrower ON r_cif_borrower."borrowerId" = r_loan_borrower."borrowerId" 
  JOIN cif ON cif.id = r_cif_borrower."cifId" 
  JOIN "group" ON "group".id = r_loan_group."groupId" 
  JOIN disbursement ON disbursement.id = r_loan_disbursement."disbursementId" 
  JOIN loan_raw ON loan_raw."loanId" = loan.id
  WHERE loan_raw._raw ? 'to_sisa' 
  AND loan.stage IN ('PRIVATE', 'MARKETPLACE', 'INVESTOR')`
  
  services.DBCPsql.Raw(query).Find(&multiLoanSchema)
  ctx.JSON(iris.StatusOK, iris.Map{
    "status": "success",
    "data": multiLoanSchema,
  })
}

