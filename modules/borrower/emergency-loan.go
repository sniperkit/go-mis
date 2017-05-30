package borrower

import (
	//"fmt"	
	
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

	if db.Table("loan_raw").Create(&loanRaw.LoanRaw{Raw: dataRaw, LoanID: newLoan.ID}).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan Raw")
		return 0
	}

	if db.Table("r_loan_sector").Create(&r.RLoanSector{LoanId: newLoan.ID, SectorId: sectorID}).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan Sector Relation")
		return 0
	}

	rLoanBorrower := r.RLoanBorrower{
		LoanId:     loan.ID,
		BorrowerId: borrowerId,
	}

	if db.Table("r_loan_borrower").Create(&rLoanBorrower).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan Borrower Relation")
		return 0
	}

	if UseProductPricing(0, newLoan.ID, db) != nil {
		processErrorAndRollback(ctx, db, "Error Use Product Pricing")
		return 0
	}

	if CreateRelationLoanToGroup(newLoan.ID, groupID, db) != nil {
		processErrorAndRollback(ctx, db, "Error Create Relation to Group")
		return 0
	}

	if CreateRelationLoanToBranch(newLoan.ID, groupID, db) != nil {
		processErrorAndRollback(ctx, db, "Error Create Relation to Branch")
		return 0
	}

	// define disbursement date
	// --- --
	if CreateDisbursementRecord(loan.ID, payload["disbursementDate"].(string), db) != nil {
		processErrorAndRollback(ctx, db, "Error Create Disbusrement")
		return 0
	}

	
	//if sourceType == "OLD" {
	//	dbSurvey := services.DBCPsqlSurvey.Begin()

	//	idCardNo := payload["client_ktp"].(string)
	//	if setOldSurveyStatus(idCardNo, "APPROVE", dbSurvey) != nil {
	//		processErrorAndRollback(ctx, db, "Error Setting Data in DB Survey")
	//		return 0
	//	}

	//	dbSurvey.Commit()
	//} else {
	//	uuid := payload["uuid"].(string)
	//	if setNewSurveyStatus(uuid, "APPROVE", db) != nil {
	//		processErrorAndRollback(ctx, db, "Error Setting Data in DB Survey")
	//		return 0
	//	}
	//}

	db.Commit()
}



