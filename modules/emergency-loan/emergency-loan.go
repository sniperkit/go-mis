package emergency_loan

import (
	
	iris "gopkg.in/kataras/iris.v4"
  "bitbucket.org/go-mis/modules/loan"
  "bitbucket.org/go-mis/modules/borrower"
	"bitbucket.org/go-mis/services"
	//loanRaw "bitbucket.org/go-mis/modules/loan-raw"
	"bitbucket.org/go-mis/modules/r"

)

func SubmitEmergencyLoan (ctx *iris.Context) {

	type Payload struct {
		EmergencyLoanBorrower
		Date			 string `json:"date"`
		SectorId   uint64 `json:"sectorId"`
		Purpose 	 string `json:"loan_purpose"`
		DisbusrsementDate string `json:"disbursement_date"`
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

	for idx, _ := range el {
		
		db := services.DBCPsql.Begin()

		borrowerID := el[idx].BorrowerId
		groupID 	 := el[idx].GroupId
		branchID	 := el[idx].BranchId
		oldLoanID  := el[idx].OldLoanId
		submittedLoanDate := el[idx].Date
		sectorID := el[idx].SectorId
		purpose := el[idx].Purpose
		disbusementDate := el[idx].DisbusrsementDate

		// get requiredData from oldLoad
		oldLoan := loan.Loan{}
		query := `select * from loan where id = ?`
		services.DBCPsql.Raw(query, oldLoanID).Scan(&oldLoan)
		
		newLoan := loan.Loan{}
		newLoan.URLPic1 = oldLoan.URLPic1 
		newLoan.URLPic2 = oldLoan.URLPic2 
		newLoan.Purpose = purpose 
		newLoan.LoanPeriod = 2
		newLoan.Tenor = 25 
		newLoan.Rate = 0.3 // hc 
		newLoan.Installment = 40000 // hc 
		newLoan.Plafond = 1000000 
		newLoan.SubmittedLoanDate = submittedLoanDate 
	
		newLoan.CreditScoreGrade = oldLoan.CreditScoreGrade 
		newLoan.CreditScoreValue  = oldLoan.CreditScoreValue 

		newLoan.Stage = "PRIVATE"
		newLoan.IsLWK = true
		newLoan.IsUPK = true

		if db.Table("loan").Create(&newLoan).Error != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Loan")
			break
		}

		// TODO: insert loan raw	

		if db.Table("r_loan_sector").Create(&r.RLoanSector{LoanId: newLoan.ID, SectorId: sectorID}).Error != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Loan Sector Relation")
			break
		}

		rLoanBorrower := r.RLoanBorrower{
			LoanId:     newLoan.ID,
			BorrowerId: borrowerID,
		}

		if db.Table("r_loan_borrower").Create(&rLoanBorrower).Error != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Loan Borrower Relation")
			break
		}

		if borrower.UseProductPricing(0, newLoan.ID, db) != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Use Product Pricing")
			break
		}

		if borrower.CreateRelationLoanToGroup(newLoan.ID, groupID, db) != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Relation to Group")
			break
		}

		if borrower.CreateRelationLoanToBranch(newLoan.ID, branchId, db) != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Relation to Branch")
			break
		}
		
		if borrower.CreateDisbursementRecord(newLoan.ID, disbusementDate, db) != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Disbusrement")
			break
		}
				
	  db.Commit()
		// update table emergency loan set newLoanId = newLoanId
		elb := EmergencyLoanBorrower{}
		services.DBCPsql.Model(elb).Where("\"deletedAt\" IS NULL AND id = ?", el[idx].EmergencyLoanBorrower.ID).UpdateColumn("newLoanId", newLoan.ID)

	}
}

