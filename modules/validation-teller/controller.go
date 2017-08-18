package validationTeller

import (
	"errors"
	"strings"
	"time"

	"log"

	"gopkg.in/kataras/iris.v4"

	"strconv"

	ins "bitbucket.org/go-mis/modules/installment"
	misUtility "bitbucket.org/go-mis/modules/utility"
	"bitbucket.org/go-mis/services"
)

// GetDataValidationTeller - Get data validation teller
// Route: /api/v2//validation-teller/getdata
func GetDataValidationTeller(ctx *iris.Context) {
	var err error
	var instalmentData ResponseGetData
	branchIDParam, _ := ctx.URLParamInt64("branchId")
	branchID := uint64(branchIDParam)
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
	instalmentData, err = FindInstallmentData(branchID, dateParam, false)
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
	notes, err := services.GetNotes(services.ConstructNotesGroupId(branchID, dateParam))
	if err == nil && len(notes) > 0 {
		instalmentData.BorrowerNotes = services.GetBorrowerNotes(notes)
		instalmentData.MajelisNotes = services.GetMajelisNotes(notes)
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   instalmentData,
	})
}

// SaveValidationTellerNotes: Save notes validation teller
// Route: /api/v2//validation-teller/group-notes/:logType/save
func SaveValidationTellerNotes(ctx *iris.Context) {
	params := struct {
		Date     string  `json:"date"`
		BranchId uint64  `json:"branchId"`
		Notes    []Notes `json:"notes"`
	}{}
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	logGroupId := services.ConstructNotesGroupId(params.BranchId, params.Date)
	log := services.Log{Data: params.Notes, GroupID: logGroupId, ArchiveID: ctx.Param("logType")}
	err = services.PostToLog(log)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   params,
	})
}

// SaveValidationTellerDetail - Save detail and notes validation teller
// Route: - /api/v2//validation-teller/detail/save
//		  - /validation-teller/borrower-notes/save
func SaveValidationTellerDetail(ctx *iris.Context) {
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

// GetValidationTellerDetail - Get detail data validation teller
// Route: /api/v2/validation-teller/detail
func GetValidationTellerDetail(ctx *iris.Context) {
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

// SaveValidationTeller - Controller
// Route: api/v2/validation-teller/save
func SaveValidationTeller(ctx *iris.Context) {
	var err error
	var installment ins.Installment

	params := struct {
		BranchID uint64 `json:"branchId"`
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
	instalmentData, err := FindInstallmentData(params.BranchID, params.Date, false)
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

// GetDataValidationAndTransfer - Get validation teller data by branch, date and stage = 'APPROVE'
// Route: /api/v2//validation-teller/branch/:branchId/date/:date
func GetDataValidationAndTransfer(ctx *iris.Context) {
	var err error
	branchParam := ctx.Param("branchId")
	intBranchID, _ := strconv.Atoi(branchParam)
	branchID := uint64(intBranchID)
	dateParam := ctx.Param("date")
	// Check branchID, if equal to 0 return error message to client
	if branchID == 0 {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":       iris.StatusBadRequest,
			"errorMessage": "Invalid Branch ID",
		})
		return
	}
	if len(strings.Trim(dateParam, " ")) == 0 {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":       iris.StatusBadRequest,
			"errorMessage": "Date can not be empty",
		})
		return
	}

	instalmentData, err := FindInstallmentData(branchID, dateParam, true)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"errorMessage": "System Error",
			"message":      err.Error(),
		})
	}
	dataTransfer, err := FindDataTransfer(time.Now().Local().Format("2006-01-02"))
	if err == nil {
		instalmentData.DataTransfer = dataTransfer
	}
	notes, err := services.GetNotes(services.ConstructNotesGroupId(branchID, dateParam))
	if err == nil && len(notes) > 0 {
		instalmentData.BorrowerNotes = services.GetBorrowerNotes(notes)
		instalmentData.MajelisNotes = services.GetMajelisNotes(notes)
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   instalmentData,
	})
}

// FindInstallmentData - function to get installment data by branch ID and date
func FindInstallmentData(branchID uint64, date string, isApprove bool) (ResponseGetData, error) {
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
					sum(i."paidInstallment") "repayment",
					sum(i.reserve) "tabungan",
					sum(i."paidInstallment"+i.reserve) "total",
					sum(i.cash_on_hand) "cashOnHand",
					sum(i.cash_on_reserve) "cashOnReserve",
					coalesce(sum(
					case
					when frequency >= 3 then l.installment+((plafond/tenor)*(frequency-1))
					when frequency >0 then l.installment*frequency
					when frequency = 0 then 0
					end
					),0) "projectionRepayment",
					coalesce(sum(
					case
					when plafond < 0 then 0
					when plafond <= 3000000 then 3000
					when plafond > 3000000 and plafond <= 5000000 then 4000
					when plafond > 5000000 and plafond <= 7000000 then 5000
					when plafond > 7000000 and plafond <= 9000000 then 6000
					when plafond > 9000000 and plafond <= 11000000 then 7000
					else 8000
					end
					),0) "projectionTabungan",
					coalesce(sum(case
					when d."disbursementDate"::date = ? then plafond end
					),0) "totalCairProj",
					coalesce(sum(case
					when d.stage = 'SUCCESS' and d."disbursementDate"::date = ? then plafond end
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
				where l."deletedAt" isnull and b.id= ? and coalesce(i."transactionDate",i."createdAt")::date = ? 
					and l.stage = 'INSTALLMENT' `
	if isApprove {
		query += ` and i."stage" = 'APPROVE' `
	}
	query += ` group by g.name, a.fullname, g.id order by a.fullname `
	err = services.DBCPsql.Raw(query, date, date, branchID, date).Scan(&queryResult).Error
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
	majelists := []MajelisId{}
	isEnableSubmit := true
	for idx, rval := range res {
		var totalRepaymentAct float64
		var totalRepaymentProj float64
		var totalRepaymentCoh float64
		var totalTabunganAct float64
		var totalTabunganProj float64
		var totalTabunganCoh float64
		var totalActualAgent float64
		var totalProjectionAgent float64
		var totalCohAgent float64
		var totalPencairanAgent float64
		var totalPencairanProjAgent float64
		var totalGagalDroppingAgent float64
		m := []Majelis{}
		for _, qrval := range queryResult {
			if rval.Agent == qrval.Fullname {
				m = append(m, Majelis{
					GroupId:             qrval.GroupId,
					Name:                qrval.Name,
					Repayment:           qrval.Repayment,
					Tabungan:            qrval.Tabungan,
					TotalActual:         qrval.Total,
					TotalProyeksi:       qrval.ProjectionRepayment + qrval.ProjectionTabungan,
					TotalCoh:            qrval.CashOnHand + qrval.CashOnReserve,
					TotalCair:           qrval.TotalCair,
					TotalCairProj:       qrval.TotalCairProj,
					TotalGagalDropping:  qrval.TotalGagalDropping,
					Status:              qrval.Status,
					CashOnHand:          qrval.CashOnHand,
					CashOnReserve:       qrval.CashOnReserve,
					ProjectionRepayment: qrval.ProjectionRepayment,
					ProjectionTabungan:  qrval.ProjectionTabungan,
				})
				if qrval.Status == "AGENT" {
					isEnableSubmit = false
				}
				majelists = append(majelists, MajelisId{GroupId: qrval.GroupId, Name: qrval.Name})
				totalRepaymentAct += qrval.Repayment
				totalRepaymentProj += qrval.ProjectionRepayment
				totalRepaymentCoh += qrval.CashOnHand
				totalTabunganAct += qrval.Tabungan
				totalTabunganProj += qrval.ProjectionTabungan
				totalTabunganCoh += qrval.CashOnReserve
				totalActualAgent += qrval.Total
				totalProjectionAgent += qrval.ProjectionRepayment + qrval.ProjectionTabungan
				totalCohAgent += qrval.CashOnHand + qrval.CashOnReserve
				totalPencairanAgent += qrval.TotalCair
				totalPencairanProjAgent += qrval.TotalCairProj
				totalGagalDroppingAgent += qrval.TotalGagalDropping
			}
		}
		res[idx].Majelis = m
		res[idx].TotalActualRepayment = totalRepaymentAct
		res[idx].TotalProjectionRepayment = totalRepaymentProj
		res[idx].TotalCohRepayment = totalRepaymentCoh
		res[idx].TotalActualTabungan = totalTabunganAct
		res[idx].TotalProjectionTabungan = totalTabunganProj
		res[idx].TotalCohTabungan = totalTabunganCoh
		res[idx].TotalActualAgent = totalActualAgent
		res[idx].TotalProjectionAgent = totalProjectionAgent
		res[idx].TotalCohAgent = totalCohAgent
		res[idx].TotalPencairanAgent = totalPencairanAgent
		res[idx].TotalPencairanProjAgent = totalPencairanProjAgent
		res[idx].TotalGagalDroppingAgent = totalGagalDroppingAgent
	}
	response.InstallmentData = res
	response.ListMajelis = majelists
	response.IsEnableSubmit = isEnableSubmit
	return response, nil
}

// GetDataValidation - Data validation before Get Validation Teller
func GetDataValidation(branchID uint64, date string) error {
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

// FindDataTransfer - Get data transfer information based on transfer date
func FindDataTransfer(date string) (DataTransfer, error) {
	// This function is not fixed yet
	// Just for temporary data
	var dataTransfer DataTransfer
	dataTransfer = DataTransfer{
		ID:                   1,
		ValidationDate:       time.Now().Local().Format("2006-01-02"),
		TransferDate:         time.Now().Local().Format("2006-01-02"),
		RepaymentID:          "REPAY-001",
		RepaymentNominal:     2000000.00,
		TabunganID:           "TABUNGAN-001",
		TabunganNominal:      5000000.00,
		GagalDroppingID:      "GAGAL-001",
		GagalDroppingNominal: 1000000.00,
	}
	return dataTransfer, nil
}
