package emergency_loan

import (
	"gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/services"
)

func FetchAllAvailableBorrower(ctx *iris.Context) {
	emergencyLoanBorrowers := []EmergencyLoanBorrower{}
	branchID := ctx.Param("branch_id")
	query := "SELECT * "
	query += "FROM emergency_loan_borrower "
	query += "WHERE \"status\"=false and \"branchId\"=?"
	services.DBCPsql.Raw(query, branchID).Find(&emergencyLoanBorrowers)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   emergencyLoanBorrowers,
	})
}
