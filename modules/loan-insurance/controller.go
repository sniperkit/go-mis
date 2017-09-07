package loanInsurance

import (
  "gopkg.in/kataras/iris.v4"
  "bitbucket.org/go-mis/services"
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

func ApplyRefund (ctx *iris.Context) {
  id := ctx.Param("loan_id")
  
  // TODO:
  // Calculate total remaining principle and refund it to investor
  
  services.DBCPsql.Table("loan").Where("id = ?", id).Update("isInsuranceRefund", "TRUE")
  ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{},
	})
}