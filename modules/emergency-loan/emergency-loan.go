package emergency_loan

import (
	
	//"time"
	iris "gopkg.in/kataras/iris.v4"
  "bitbucket.org/go-mis/modules/loan"
  "bitbucket.org/go-mis/modules/borrower"
	"bitbucket.org/go-mis/services"
	"bitbucket.org/go-mis/modules/r"
)


func SubmitEmergencyLoan (ctx *iris.Context) {

	type Payload struct {
		EmergencyLoanBorrower
		Date			 string `json:"date"`
		SectorId   uint64 `json:"sectorId"`
		Purpose 	 string `json:"loan_purpose"`
		DisbursementDate string `json:"disbursement_date"`
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
		oldLoanID  := el[idx].OldLoanId
		sectorID := el[idx].SectorId
		purpose := el[idx].Purpose

		/*
		disbursementDate, err := time.Parse(DATE_LAYOUT, el[idx].DisbursementDate + DATE_TAIL)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		dd := disbursementDate.String()
		dd = dd[:len(dd)-10] // -_-"
		
		submittedLoanDate, err := time.Parse(DATE_LAYOUT, el[idx].Date + DATE_TAIL)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		sld := submittedLoanDate.String()
		sld = sld[:len(sld)-10] // get rid of "0000 UTC" thingy -_-"
		*/

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
		newLoan.Installment = 43000 // hc
		newLoan.Plafond = 1000000 
		
	
		newLoan.SubmittedLoanDate = el[idx].Date 
	
		newLoan.CreditScoreGrade = oldLoan.CreditScoreGrade 
		newLoan.CreditScoreValue  = oldLoan.CreditScoreValue 

		newLoan.Stage = "PRIVATE"
		newLoan.IsLWK = true
		newLoan.IsUPK = true
		newLoan.Subgroup= oldLoan.Subgroup

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

		if borrower.UseProductPricing(0,newLoan.ID,services.DBCPsql)!=nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Relation to Product pricing")
			break
		}
		if borrower.CreateRelationLoanToGroup(newLoan.ID, groupID, db) != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Relation to Group")
			break
		}

		if borrower.CreateRelationLoanToBranch(newLoan.ID, groupID, db) != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Relation to Branch")
			break
		}
		
		if borrower.CreateDisbursementRecord(newLoan.ID, el[idx].DisbursementDate, db) != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Disbusrement")
			break
		}

		// update table emergency loan set newLoanId = newLoanId
		// only do this after all process above has completed
		elb := EmergencyLoanBorrower{}
	  err := services.DBCPsql.Model(elb).Where("\"deletedAt\" IS NULL AND id = ?", el[idx].EmergencyLoanBorrower.ID).UpdateColumns(EmergencyLoanBorrower{NewLoanId:newLoan.ID, Status:true}).Error
		if err != nil {
			borrower.ProcessErrorAndRollback(ctx, db, "Error Create Emergency Loan")
			break	
		}

	  db.Commit()
	}
}

