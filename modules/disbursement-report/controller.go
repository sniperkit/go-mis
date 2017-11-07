package disbursementReport

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
	"fmt"
	"bitbucket.org/go-mis/config"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
)

func Init() {
	services.DBCPsql.AutoMigrate(&DisbursementReport{})
	services.BaseCrudInit(DisbursementReport{}, []DisbursementReport{})
}

func FetchAllActive(ctx *iris.Context) {
	query := "SELECT * "
	query += "FROM disbursement_report "
	query += "WHERE \"isActive\"=true and \"deletedAt\" IS NULL"
	fmt.Println(query)
	disbursementReports := []DisbursementReport{}
	services.DBCPsql.Raw(query).Scan(&disbursementReports)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   disbursementReports,
	})
}

func GetDetail(ctx *iris.Context) {
	// select report from DB
	query := "SELECT * "
	query += "FROM disbursement_report "
	query += "WHERE \"id\"=?"
	id := ctx.Param("id")
	disbursementReport:=DisbursementReport{}
	services.DBCPsql.Raw(query,id).Scan(&disbursementReport)
	if (disbursementReport == DisbursementReport{}) {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status": "error",
			"message":   "Report not found",
		})
		return
	}
	//get json from node uploader
	urlStr := config.UploaderApiPath + "/disbursement"+disbursementReport.Filename+"?secretKey=n0de-U>lo4d3r"
	request := gorequest.New()
	_, resReportStr, errs := request.Get(urlStr).End()
	disbursementReportDetail := DisbursementReportDetail{}
	json.Unmarshal([]byte(resReportStr), &disbursementReportDetail)
	if len(errs) > 0 {
		fmt.Println(errs)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"message":   "Error integration to uploader",
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   disbursementReportDetail,
	})
}