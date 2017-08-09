package validationTeller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	"log"
	"net/http"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"

	ins "bitbucket.org/go-mis/modules/installment"
	misUtility "bitbucket.org/go-mis/modules/utility"
)

var STAGE map[int]string = map[int]string{
	1: "AGENT",
	2: "TELLER",
	3: "PENDING",
	4: "REVIEW",
	5: "APPROVE",
}

// GetData - Get data validation teller
func GetData(ctx *iris.Context) {
	var err error
	branchID, _ := ctx.URLParamInt64("branchId")
	dateParam := ctx.URLParam("date")
	// Check data whehter valid or not
	err = GetDataValidation(branchID, dateParam)
	if err != nil {
		log.Println("[INFO] Params is not valid")
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}
	log.Println("[INFO] Validation get data installment pass")
	instalmentData, err := FindInstallmentData(branchID, dateParam)
	if err != nil {
		log.Println("[INFO] Can not retrive installment data")
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"message":      "Internal Server Error",
				"errorMessage": "Unable to retrive data from LOG",
			})
			return
		}
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   instalmentData,
	})
	// Change Requirement
	// Always send data form DB
	// dataDate, _ := misUtility.StringToDate(dateParam)
	// if misUtility.IsBeforeToday(dataDate) {
	// 	log.Println(dataDate)
	// 	dataLog, err := GetDataFromLog(branchID)
	// 	if err != nil {
	// 		ctx.JSON(iris.StatusInternalServerError, iris.Map{
	// 			"message":      "Internal Server Error",
	// 			"errorMessage": "Unable to retrive data from LOG",
	// 		})
	// 		return
	// 	}
	// 	ctx.JSON(iris.StatusOK, iris.Map{
	// 		"status": "success",
	// 		"data":   dataLog.Data,
	// 	})
	// 	return
	// }
	// if misUtility.IsToday(dataDate) {
	// 	ctx.JSON(iris.StatusOK, iris.Map{
	// 		"status": "success",
	// 		"data":   instalmentData,
	// 	})
	// }
}

func SaveDetail(ctx *iris.Context) {
	params := []struct {
		CashOnReserve float64 `json:"cashOnReserve"`
		CashOnHand    float64 `json:"cashOnHand"`
		Id            uint64  `json:"id"`
		GroupId       uint64  `json:"groupId"`
		Note          string  `json:"note"`
	}{}
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	db := services.DBCPsql.Begin()
	for _, param := range params {
		if err := db.Table("installment").Where("\"id\" = ?", param.Id).UpdateColumn("cash_on_hand", param.CashOnHand).Error; err != nil {
			db.Rollback()
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": "Error Update cashOnHand",
			})
			return
		}
		if err := db.Table("installment").Where("\"id\" = ?", param.Id).UpdateColumn("cash_on_reserve", param.CashOnReserve).Error; err != nil {
			db.Rollback()
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": "Error Update cashOnReserve",
			})
			return
		}
	}
	db.Commit()

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   params,
	})
}

func GetDetail(ctx *iris.Context) {
	params := struct {
		Date    string `json:"date"`
		GroupId int64  `json:"groupId"`
	}{}
	params.GroupId, _ = ctx.URLParamInt64("groupId")
	params.Date = ctx.URLParam("date")

	query := `select i.id,rlbo."borrowerId" as "borrowerId",cif."name", 
					i."paidInstallment" as "repayment",i.reserve as "tabungan",
					(i."paidInstallment"+i.reserve) as "total",
					i.cash_on_hand as "cashOnHand",
					i.cash_on_reserve as "cashOnReserve"
				from loan l join r_loan_group rlg on l.id = rlg."loanId"
					join "group" g on g.id = rlg."groupId"
					join r_group_agent rga on g.id = rga."groupId"
					join agent a on a.id = rga."agentId"
					join r_loan_branch rlb on rlb."loanId" = l.id
					join r_loan_borrower rlbo on rlbo."loanId" = l.id
					join r_cif_borrower on r_cif_borrower."borrowerId"=rlbo."borrowerId"
					join cif on cif.id=r_cif_borrower."cifId"
					join branch b on b.id = rlb."branchId"
					join r_loan_installment rli on rli."loanId" = l.id
					join installment i on i.id = rli."installmentId"
					join r_loan_disbursement rld on rld."loanId" = l.id
					join disbursement d on d.id = rld."disbursementId"
				where l."deletedAt" isnull and coalesce(i."transactionDate",i."createdAt")::date = ? and 
				l.stage = 'INSTALLMENT' and g.id=?`

	queryResult := []RawInstallmentDetail{}
	services.DBCPsql.Raw(query, params.Date, params.GroupId).Scan(&queryResult)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   queryResult,
	})
}

func UpdateInstallmentStage(ctx *iris.Context) {
	params := struct {
		InstallmentID int `json:"installmentId"`
		Stage         int `json:"stage"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}

	if err := services.DBCPsql.Table("installment").Where("id = ?", params.InstallmentID).UpdateColumn("stage", STAGE[params.Stage]).Error; err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   "Installment id:" + strconv.Itoa(params.InstallmentID) + " updated. Stage:" + STAGE[params.Stage],
	})
}

// SubmitValidationTeller - Controller
func SubmitValidationTeller(ctx *iris.Context) {
	var err error
	var installment ins.Installment

	params := struct {
		BranchID int64  `json:"branchId"`
		Date     string `json:"date"`
	}{}

	err = ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"errorMessage": "Bad Request",
			"message":      "Can not Unmarshall JSON Body",
		})
		return
	}
	installments, err := ins.FindByBranchAndDate(params.BranchID, params.Date)
	if err != nil {
		log.Println("#ERROR: ", err)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"errorMessage": "System Error",
			"message":      err.Error(),
		})
		return
	}
	// db.begin
	db := services.DBCPsql.Begin()
	for _, installment = range installments {
		err := ins.UpdateStageInstallmentApproveOrReject(db, installment.ID, installment.Stage, "PENDING")
		if err != nil {
			log.Println("#ERROR: ", err)
			db.Rollback()
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"errorMessage": "System Error",
				"message":      err.Error(),
			})
			return
		}
	}
	db.Commit()
	instalmentData, err := FindInstallmentData(params.BranchID, params.Date)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"errorMessage": "System Error",
			"message":      err.Error(),
		})
	}
	// Create go routine to push data to GO-LOG APP
	// in order to no need to wait
	go func() {
		_ = services.PostToLog(services.GetLog(params.BranchID, instalmentData, "VALIDATION TELLER"))
	}()
	ctx.JSON(iris.StatusOK, iris.Map{
		"message": "Success",
	})

}

// FindInstallmentData - function to get installment data by branch ID and date
func FindInstallmentData(branchID int64, date string) (ResponseGetData, error) {
	var err error
	queryResult := []RawInstallmentData{}
	response := ResponseGetData{}
	res := []InstallmentData{}
	agents := map[string]bool{"": false}
	if branchID <= 0 {
		return response, errors.New("Invalid Branch ID")
	}
	if len(strings.Trim(date, " ")) == 0 {
		return response, errors.New("Invalid Date")
	}
	query := `select g.id as "groupId", a.fullname,g.name, 
					sum(i."paidInstallment") "repayment",sum(i.reserve) "tabungan",
					sum(i."paidInstallment"+i.reserve) "total", 
					sum(i.cash_on_hand) "cashOnHand",
					sum(i.cash_on_reserve) "cashOnReserve",
					coalesce(sum(case
					when d.stage = 'SUCCESS' then plafond end
					),0) "totalCair",
					coalesce(sum(case
					when d.stage = 'FAILED' then plafond end
					),0) "totalGagalDropping",
					split_part(string_agg(i.stage,'| '),'|',1) "status"
				from loan l join r_loan_group rlg on l.id = rlg."loanId"
					join "group" g on g.id = rlg."groupId"
					join r_group_agent rga on g.id = rga."groupId"
					join agent a on a.id = rga."agentId"
					join r_loan_branch rlb on rlb."loanId" = l.id
					join branch b on b.id = rlb."branchId"
					join r_loan_installment rli on rli."loanId" = l.id
					join installment i on i.id = rli."installmentId"
					join r_loan_disbursement rld on rld."loanId" = l.id
					join disbursement d on d.id = rld."disbursementId"
				where l."deletedAt" isnull and b.id= ? and coalesce(i."transactionDate",i."createdAt")::date = ? and 
				l.stage = 'INSTALLMENT'
				group by g.name, a.fullname, g.id
				order by a.fullname`
	err = services.DBCPsql.Raw(query, branchID, date).Scan(&queryResult).Error
	if err != nil {
		log.Println("#ERROR: Unable to retrieve Installment data")
		log.Println("#ERROR: ", err)
		return response, errors.New("Unable to retrieve Installment data")
	}
	for _, val := range queryResult {
		if agents[val.Fullname] == false {
			agents[val.Fullname] = true
			res = append(res, InstallmentData{Agent: val.Fullname})
		}
	}

	for idx, rval := range res {
		m := []Majelis{}
		var totalRepayment float64
		var totalCashOnHand float64
		var totalTabungan float64
		var totalCashOnReserve float64
		var totalCair float64
		var totalGagalDroping float64
		for _, qrval := range queryResult {
			if rval.Agent == qrval.Fullname {
				m = append(m, Majelis{
					GroupId:            qrval.GroupId,
					Name:               qrval.Name,
					Repayment:          qrval.Repayment,
					Tabungan:           qrval.Tabungan,
					Total:              qrval.Total,
					TotalCair:          qrval.TotalCair,
					TotalGagalDropping: qrval.TotalGagalDropping,
					Status:             qrval.Status,
					CashOnHand:         qrval.CashOnHand,
					CashOnReserve:      qrval.CashOnReserve,
				})
				totalRepayment += qrval.Repayment
				totalCashOnHand += qrval.CashOnHand
				totalCashOnReserve += qrval.CashOnReserve
				totalTabungan += qrval.Tabungan
				totalCair += qrval.TotalCair
				totalGagalDroping += qrval.TotalGagalDropping
			}
		}
		response.TotalActualRepayment+=totalRepayment
		response.TotalCair+=totalCair
		response.TotalTabungan+=totalTabungan
		response.TotalGagalDroping+=totalGagalDroping
		response.TotalCashOnReserve+=totalCashOnReserve
		response.TotalCashOnHand+=totalCashOnHand

		res[idx].Majelis = m
		res[idx].TotalActualRepayment = totalRepayment
		res[idx].TotalCair = totalCair
		res[idx].TotalCashOnHand = totalCashOnHand
		res[idx].TotalCashOnReserve = totalCashOnReserve
		res[idx].TotalGagalDroping = totalGagalDroping
		res[idx].TotalTabungan = totalTabungan
	}
	response.InstallmentData = res
	return response, nil
}

// GetDataFromLog - Retrive data from GO-LOG App
func GetDataFromLog(branchID int64) (Log, error) {
	var logger Log
	archiveID := services.GenerateArchiveID(branchID)
	groupID := "VALIDATION TELLER"
	apiPath := services.GetLogAPIPath() + "archive/" + archiveID + "/group/" + groupID
	log.Println("[INFO]", apiPath)
	resp, err := http.Get(apiPath)
	if err != nil {
		log.Println("#ERROR: Unable to retrive data from GO-LOG App")
		log.Println("#ERROR: ", err)
		return logger, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("#ERROR: When read body reponse from GO-LOG")
		return logger, err
	}
	log.Println("[INFO]", body)
	err = json.Unmarshal([]byte(body), &logger)
	if err != nil {
		log.Println("#ERROR: When unmarshall resp body GO-LOG to Log struct")
		return logger, err
	}
	log.Println("[INFO]", logger)
	return logger, nil
}

// GetDataValidation - Data validation before Get Validation Teller
func GetDataValidation(branchID int64, date string) error {
	if branchID <= 0 {
		return errors.New("Invalid Branch ID")
	}
	_, err := misUtility.StringToDate(date)
	if err != nil || len(strings.Trim(date, " ")) == 0 {
		return errors.New("Invalid Date")
	}
	// if misUtility.IsAfterToday(dataDate) {
	// 	return errors.New("Date must be today or less than today")
	// }
	return nil
}
