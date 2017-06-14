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
	loanRawID := ctx.Param("id")

	query := "select * from loan_raw where id = ?"

	m := LoanRaw{}
	if e := services.DBCPsql.Raw(query, loanRawID).Scan(&m).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
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
