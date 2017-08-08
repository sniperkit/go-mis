package validationTeller

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
	"bytes"
	"strings"
	"net/http"
	"time"
	"log"
	ins "bitbucket.org/go-mis/modules/installment"
	"errors"
	"encoding/json"
	misConfig "bitbucket.org/go-mis/config"
)

var (
	logAPIPath = misConfig.GoLogPath + "archive"
)

var STAGE map[int]string = map[int]string{
	1: "AGENT",
	2: "TELLER",
	3: "PENDING",
	4: "REVIEW",
	5: "APPROVE",
}

func GetData(ctx *iris.Context) {
	params := struct {
		BranchId int64 `json:"branchId"`
		Date     string `json:"date"`
	}{}
	params.BranchId,_=ctx.URLParamInt64("branchId")
	params.Date=ctx.URLParam("date")

	query := `select g.id as "groupId", a.fullname,g.name, sum(i."paidInstallment") "repayment",sum(i.reserve) "tabungan",sum(i."paidInstallment"+i.reserve) "total", 
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
				where l."deletedAt" isnull and b.id= ? and coalesce(i."transactionDate",i."createdAt")::date = ? and l.stage = 'INSTALLMENT'
				group by g.name, a.fullname, g.id
				order by a.fullname`

	queryResult := []RawInstallmentData{}
	services.DBCPsql.Raw(query, params.BranchId, params.Date).Scan(&queryResult)

	res := []InstallmentData{}
	agents := map[string]bool{"": false}
	for _, val := range queryResult {
		if agents[val.Fullname] == false {
			agents[val.Fullname] = true
			res = append(res, InstallmentData{Agent: val.Fullname})
		}
	}

	for idx, rval := range res {
		m := []Majelis{}
		var totalRepayment float64
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
					Status: qrval.Status,
					CashOnHand:qrval.CashOnHand,
					CashOnReserve:qrval.CashOnReserve,
				})
				totalRepayment += qrval.Repayment
			}
		}
		res[idx].Majelis = m
		res[idx].TotalActualRepayment = totalRepayment
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   res,
	})
}

func SaveDetail(ctx *iris.Context) {
	params := [] struct {
		CashOnReserve  float64 `json:"cashOnReserve"`
		CashOnHand     float64 `json:"cashOnHand"`
		Id  uint64 `json:"id"`
		GroupId  	   uint64 `json:"groupId"`
		Note           string `json:"note"`
	}{}
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	db :=services.DBCPsql.Begin()
	for _,param := range params {
		if err:=db.Table("installment").Where("\"id\" = ?", param.Id).UpdateColumn("cash_on_hand", param.CashOnHand).Error;err!=nil{
			db.Rollback()
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": "Error Update cashOnHand",
			})
			return
		}
		if err:=db.Table("installment").Where("\"id\" = ?", param.Id).UpdateColumn("cash_on_reserve", param.CashOnReserve).Error;err!=nil{
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
		Date     string `json:"date"`
		GroupId  int64 `json:"groupId"`
	}{}
	params.GroupId,_=ctx.URLParamInt64("groupId")
	params.Date=ctx.URLParam("date")

	query := `select i.id,rlbo."borrowerId" as "borrowerId",cif."name", i."paidInstallment" as "repayment",i.reserve as "tabungan",(i."paidInstallment"+i.reserve) as "total",
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
				where l."deletedAt" isnull and coalesce(i."transactionDate",i."createdAt")::date = ? and l.stage = 'INSTALLMENT' and g.id=?`

	queryResult := []RawInstallmentDetail{}
	services.DBCPsql.Raw(query, params.Date, params.GroupId).Scan(&queryResult)


	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   queryResult,
	})
}

// ValidationTeller - Controller
func SubmitValidationTeller(ctx *iris.Context) {
	var err error
	var installment ins.Installment

	params := struct {
		BranchId string `json:"branchId"`
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
	installments, err := ins.FindByBranchAndDate(params.BranchId, params.Date)
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

		err:=ins.UpdateStageInstallmentApproveOrReject(db, installment.ID, installment.Stage, "PENDING")
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
	go postToLog(getLog(params.BranchId, params))
	ctx.JSON(iris.StatusOK, iris.Map{
		"message": "Success",
	})

}

func getLog(branchID string, data interface{}) Log {
	var logger Log
	if len(strings.Trim(branchID, " ")) == 0 || len(strings.Trim(branchID, " ")) == 0 {
		return logger
	}
	logger = Log{
		GroupID:   "Validasi Teller",
		ArchiveID: generateArchiveID(branchID),
		Data:      data,
	}
	return logger
}

func postToLog(l Log) error {
	if l.Data == nil {
		return errors.New("Can not send empty data")
	}
	logBytes := new(bytes.Buffer)
	json.NewEncoder(logBytes).Encode(l)
	log.Println(logAPIPath)
	log.Println(l)
	resp, err := http.Post(logAPIPath, "application/json; charset=utf-8", logBytes)
	log.Println(resp.Status)
	resp.Body.Close()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func generateArchiveID(branchID string) string {
	if len(strings.Trim(branchID, " ")) == 0 {
		return ""
	}
	return branchID + "-" + time.Now().Local().Format("2006-01-02")
}
