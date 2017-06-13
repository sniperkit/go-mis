package disbursementReport

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&DisbursementReport{})
	services.BaseCrudInit(DisbursementReport{}, []DisbursementReport{})
}

func FetchAllActive(ctx *iris.Context) {
	//TODO select from database where isActive and deletedAt is null
	disbursementReports := []DisbursementReport{
		{ID:1,Filename:"A",IsActive:true,DisbursementDateFrom:"2017-06-07",DisbursementDateTo:"2017-06-10"},
		{ID:2,Filename:"A",IsActive:true,DisbursementDateFrom:"2017-06-17",DisbursementDateTo:"2017-06-20"}}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   disbursementReports,
	})
}

func GetDetail(ctx *iris.Context) {
	//TODO select from database where isActive and deletedAt is null
	disbursementReportDetail := DisbursementReportDetail{
		Dates:[]string{"2017-06-07", "2017-06-08","2017-06-09","2017-06-10"},
		Details:[]DisbursementArea{
			DisbursementArea{
				Name:"Bogor",
				Branchs:[]DisbursementBranch{
					DisbursementBranch{Name:"Baranang Siang",Prices:[]float64{10000,10000,10000,10000}},
					DisbursementBranch{Name:"Padjajaran",Prices:[]float64{20000,20000,20000,20000}},
				}},
			DisbursementArea{
				Name:"Bandung",
				Branchs:[]DisbursementBranch{
					DisbursementBranch{Name:"Cimahi",Prices:[]float64{10000,10000,10000,10000}},
					DisbursementBranch{Name:"Ciampelas",Prices:[]float64{15000,15000,15000,15000}},
				}},
		},
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   disbursementReportDetail,
	})
}