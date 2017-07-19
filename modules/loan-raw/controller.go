package loanRaw

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&LoanRaw{})
	services.BaseCrudInit(LoanRaw{}, []LoanRaw{})
}

func GetLoanRawById(ctx *iris.Context){
	loanID := ctx.Param("id")

	query := "select * from loan_raw where \"loanId\" = ?"

	m := LoanRaw{}
	if e := services.DBCPsql.Raw(query, loanID).Scan(&m).Error; e != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status": "Error",
			"message":"Data not available",
			"data":   e,
		})
		return
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   m,
	})
	}
}
