package borrower

import (
	"fmt"	
	
	iris "gopkg.in/kataras/iris.v4"
//	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/services"

)

func CreateEmergencyLoan (ctx *iris.Context) {

	type Payload struct {
		BorrowerId uint64 `json:borrowerId`
		GroupId		 uint64 `json:groupId`
		BranchId 	 uint64 `json:branchId`
		OldLoanId	 uint64 `json:loanId`
		Date			 string `json:date`
		SectorId   uint64 `json:sectorId`
		Purpose 	 string `json:loan_purpose`
	}
	
	el := []Payload{}
	if err := ctx.ReadJSON(&el); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// process each data

	db := services.DBCPsql.Begin()
	for idx, _ := range el {

		borrowerID := el[idx].BorrowerId
		groupID 	 := el[idx].GroupId
		branchID	 := el[idx].BranchId
		oldLoanID  := el[idx].OldLoanId
		SubmittedDate := el[idx].Date
		sectorID := el[idx].SectorId
		purpose := el[idx].Purpose

		// get requiredData from oldLoad
		oldLoan := loan.Loan{}
		query := `select * from loan where id = ?`
		services.DBCPsql.Raw(query, oldLoanID).Scan(&oldLoan)
		
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
		
	}
	
	db.Commit()
}



