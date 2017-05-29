package borrower

import (
	"fmt"	
	
	iris "gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/services"

)

func CreateEmergencyLoan (ctx *iris.Context) {
	// from payload:
	// - get borrower id
	borrower_id := 27251
	// - oldLoanId
	oldLoanId := 37139
	// submittedDate
	submittedDate = "2013-02-02"
	// sector 
	sector := 10

	// get requiredData from oldLoad
	oldLoan := loan.Loan{}
	query := `select * from loan where id = ?`
	services.DBCPsql.Raw(query, oldLoanId).Scan(&oldLoan)

	newLoan = loan.Loan{}
	newLoan.URLPic1 = oldLoan.URLPic1 
	newLoan.URLPic2 = oldLoan.URLPic2 
	newLoan.Purpose = oldLoan.Purpose 

	newLoan.LoanPeriod = 2

	newLoan.Tenor = 25 
	newLoan.Rate = 0.3 // hc 
	newLoan.Installment = 40000 // hc 
	newLoan.Plafond = 1000000 

	newLoan.CreditScoreGrade = oldLoan.CreditScoreGrade; 
	newLoan.CreditScoreValue, _ = oldLoan.CreditScoreValue 

	newLoan.Stage = "PRIVATE"
	newLoan.IsLWK = true
	newLoan.IsUPK = true

	db := services.DBCPsql.Begin()
	if db.Table("loan").Create(&newLoan).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan")
		return 0
	}

	// loan raw

	// loan sector

	// r loan borrower	

	// product pricing

	// loan to group 

	// loan to branch

	// disbursment

	db.Commit()

}



