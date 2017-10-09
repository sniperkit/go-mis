package validationTeller

import (
	"errors"
	"strings"

	"log"

	"gopkg.in/kataras/iris.v4"

	"strconv"

	MISInstallment "bitbucket.org/go-mis/modules/installment"
	SystemParameter "bitbucket.org/go-mis/modules/system-parameter"
	MISUtility "bitbucket.org/go-mis/modules/utility"
	"bitbucket.org/go-mis/services"
)

const (
	agentStatus  = "AGENT"
	tellerStatus = "TELLER"
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

	if !SystemParameter.IsAllowedBackdate(dateParam) {
		log.Println("#ERROR: Not Allowed back date")
		ctx.JSON(405, iris.Map{
			"message":      "Not Allowed",
			"errorMessage": "View back date is not allowed",
		})
		return
	}
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

// SaveValidationTellerNotes - save validation teller additional information / notes
// Routes: api/v2/validation-teller/group-notes/:logType/save
func SaveValidationTellerNotes(ctx *iris.Context) {
	params := struct {
		Date     string  `json:"date"`
		BranchID uint64  `json:"branchId"`
		Notes    []Notes `json:"notes"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	logGroupID := services.ConstructNotesGroupId(params.BranchID, params.Date)
	log := services.Log{Data: params.Notes, GroupID: logGroupID, ArchiveID: ctx.Param("logType")}
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

// SaveValidationTellerDetail - save data validation teller and additional information/note
// Routes: api/v2/validation-teller/detail/save
func SaveValidationTellerDetail(ctx *iris.Context) {
	params := []struct {
		CashOnReserve     float64 `json:"cashOnReserve"`
		CashOnHand        float64 `json:"cashOnHand"`
		ID                uint64  `json:"id"`
		GroupID           uint64  `json:"groupId"`
		Note              string  `json:"note,omitempty"`
		CashOnHandNote    string  `json:"cashOnHandNote"`
		CashOnReserveNote string  `json:"cashOnReserveNote"`
	}{}
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	db := services.DBCPsql.Begin()
	for _, param := range params {
		cashOnHandUpdate := map[string]interface{}{"cash_on_hand": param.CashOnHand, "cash_on_hand_note": param.CashOnHandNote, "cash_on_reserve_note": param.CashOnReserveNote}
		if err := db.Table("installment").Where("\"id\" = ?", param.ID).Updates(cashOnHandUpdate).Error; err != nil {
			db.Rollback()
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": "Error Update cashOnHand",
			})
			return
		}
		cashOnReserveUpdate := map[string]interface{}{"cash_on_reserve": param.CashOnReserve, "cash_on_hand_note": param.CashOnHandNote, "cash_on_reserve_note": param.CashOnReserveNote}
		if err := db.Table("installment").Where("\"id\" = ?", param.ID).Updates(cashOnReserveUpdate).Error; err != nil {
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
	var installmentDetails []RawInstallmentDetail
	params := struct {
		Date    string `json:"date"`
		GroupID int64  `json:"groupId"`
	}{}
	params.GroupID, _ = ctx.URLParamInt64("groupId")
	params.Date = ctx.URLParam("date")
	installmentDetails, err := FindVTDetailByGroupAndDate(uint64(params.GroupID), params.Date)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   installmentDetails,
	})
}

// SaveValidationTeller - Controller
// Route: api/v2/validation-teller/save
func SaveValidationTeller(ctx *iris.Context) {
	var err error
	var installment MISInstallment.Installment

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
	installments, err := MISInstallment.FindByBranchAndDate(params.BranchID, params.Date)
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
		err := MISInstallment.UpdateStageInstallmentApproveOrReject(db, installment.ID, installment.Stage, "PENDING")
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
	dataTransfer, err := FindDataTransfer(branchID, dateParam)
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

// SaveRejectNotes - save note rejection of Validation Teller
// Routes: api/v2/reject-notes/:status/:stage/save
func SaveRejectNotes(ctx *iris.Context) {
	params := struct {
		Date    string `json:"date"`
		GroupID uint64 `json:"groupId"`
		Notes   string `json:"notes"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	logGroupID := services.ConstructRejectsNotesGroupId(params.GroupID, params.Date, ctx.Param("status"), ctx.Param("stage"))
	dataLog := services.Log{Data: params.Notes, GroupID: logGroupID, ArchiveID: ctx.Param("stage")}
	err = services.PostToLog(dataLog)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   params,
	})
}

// GetRejectNotes - Get rejection notes of Validation teller
// Routes: api/v2/reject-notes/:status/:stage/get/:groupId/:date
func GetRejectNotes(ctx *iris.Context) {
	l, err := services.GetRejectNotesData(ctx.Param("status"), ctx.Param("groupId"), ctx.Param("date"), ctx.Param("stage"))
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	response := struct {
		GroupID string      `json:"groupId"`
		Date    string      `json:"date"`
		Stage   string      `json:"archiveId"`
		Notes   interface{} `json:"notes"`
	}{
		GroupID: ctx.Param("groupId"),
		Date:    ctx.Param("date"),
		Stage:   ctx.Param("stage"),
		Notes:   l.Data,
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   response,
	})
}

// FindInstallmentData - function to get installment data by branch ID and date
func FindInstallmentData(branchID uint64, date string, isApprove bool) (ResponseGetData, error) {
	// Declare and set variables with zero values
	var err error
	var rawInstallmentData []MISInstallment.RawInstallmentData
	var responseData ResponseGetData
	var installmentData []InstallmentData
	agents := make(map[string]bool)

	// branchID and date must be valid
	if branchID <= 0 {
		return responseData, errors.New("Invalid Branch ID")
	}
	if len(strings.Trim(date, " ")) == 0 {
		return responseData, errors.New("Date can not be empty")
	}
	_, err = MISUtility.StringToDate(date)
	if err != nil {
		return responseData, errors.New("Invalid Date parameter")
	}
	rawInstallmentData = MISInstallment.GetRawPendingInstallmentData("teller",branchID, date, isApprove)
	for _, val := range rawInstallmentData {
		if agents[val.Fullname] == false {
			agents[val.Fullname] = true
			installmentData = append(installmentData, InstallmentData{Agent: val.Fullname})
		}
	}
	totalCabang, majelists, isEnableSubmit := GetTotalCabang(rawInstallmentData, installmentData)
	log.Println("Majelists: ", majelists)
	responseData.InstallmentData = installmentData
	responseData.ListMajelis = majelists
	responseData.IsEnableSubmit = isEnableSubmit
	AssignTotalResponseData(&responseData, totalCabang)
	return responseData, nil
}

// GetTotalCabang - Get summary of data Validation Teller
func GetTotalCabang(rawInstallmentData []MISInstallment.RawInstallmentData, installmentData []InstallmentData) (*TotalCabang, []MajelisId, bool) {
	majelists := make([]MajelisId, 0, len(rawInstallmentData))
	var majelis Majelis
	isEnableSubmit := true
	var tellerCounter int
	totalCabang := new(TotalCabang)
	for idx, rval := range installmentData {
		m := make([]Majelis, 0, len(rawInstallmentData))
		totalRepayment := new(TotalRepayment)
		for _, qrval := range rawInstallmentData {
			if rval.Agent == qrval.Fullname {
				m = append(m, majelis.InitializedByRawInstallmentData(qrval))
				if strings.ToUpper(qrval.Status) == agentStatus {
					isEnableSubmit = false
				}
				if strings.ToUpper(qrval.Status) == tellerStatus {
					tellerCounter++
				}
				if qrval.GroupId > 0 && len(strings.Trim(qrval.Name, " ")) > 0 {
					majelists = append(majelists, MajelisId{GroupId: qrval.GroupId, Name: qrval.Name})
				}
				totalRepayment.AddTotal(qrval)
			}
		}
		installmentData[idx].Majelis = m
		installmentData[idx].AddTotal(totalRepayment)
		totalCabang.AddTotal(totalRepayment)
	}
	if tellerCounter == 0 {
		isEnableSubmit = false
	}
	return totalCabang, majelists, isEnableSubmit
}

// GetDataValidation - Data validation before Get Validation Teller
func GetDataValidation(branchID uint64, date string) error {
	if branchID <= 0 {
		return errors.New("Invalid Branch ID")
	}
	_, err := MISUtility.StringToDate(date)
	if err != nil || len(strings.Trim(date, " ")) == 0 {
		return errors.New("Invalid Date")
	}
	return nil
}

// FindDataTransfer - Get data transfer information based on transfer date
func FindDataTransfer(branchID uint64, date string) (DataTransfer, error) {
	var dataTransfer DataTransfer
	if len(strings.Trim(date, " ")) == 0 {
		log.Println("#ERROR: Date is empty")
		return dataTransfer, errors.New("Date can not be empty")
	}
	_, err := MISUtility.StringToDate(date)
	if err != nil {
		log.Println("#ERROR: Invalid date parameter", date)
		return dataTransfer, errors.New("Invalid date parameter")
	}
	query := `select data_transfer.id,
				data_transfer.validation_date,
				data_transfer.transfer_date,
				data_transfer.transfer_date,
				data_transfer.repayment_id,
				data_transfer.repayment_nominal,
				data_transfer.tabungan_id,
				data_transfer.tabungan_nominal,
				data_transfer.gagal_dropping_id,
				data_transfer.gagal_dropping_nominal,
				data_transfer.gagal_dropping_note,
				data_transfer.branch_id
			from data_transfer
			where data_transfer.validation_date::date = ?
			AND branch_id = ? 
			order by data_transfer.id DESC`
	err = services.DBCPsql.Raw(query, date, branchID).Scan(&dataTransfer).Error
	if err != nil {
		log.Println("#ERROR: Unable to retrieve data transfer", err)
		return dataTransfer, errors.New("Unable to retrive data transfer")
	}
	return dataTransfer, nil
}

// FindVTDetailByGroupAndDate - Get detail VT by branch ID, transaction and created date
func FindVTDetailByGroupAndDate(groupID uint64, date string) ([]RawInstallmentDetail, error) {
	var installmentDetails []RawInstallmentDetail
	if groupID <= 0 {
		return installmentDetails, errors.New("Group ID can not be empty")
	}
	if len(strings.Trim(date, " ")) == 0 {
		return installmentDetails, errors.New("Date can not be empty")
	}
	_, err := MISUtility.StringToDate(date)
	if err != nil {
		return installmentDetails, errors.New("Invalid Date")
	}
	query := `select i.id,rlbo."borrowerId" as "borrowerId",cif."name", 
					i."paidInstallment" as "repayment",i.reserve as "tabungan",
					(i."paidInstallment"+i.reserve) as "total",
					i.cash_on_hand as "cashOnHand",
					i.cash_on_reserve as "cashOnReserve",
					i.cash_on_reserve_note as "cashOnReserveNote",
					i.cash_on_hand_note as "cashOnHandNote"
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
				where l."deletedAt" is null and i."deletedAt" is null 
				and coalesce(i."transactionDate",i."createdAt")::date = ? and
				l.stage = 'INSTALLMENT' and g.id=? and i."deletedAt" is not null`
	err = services.DBCPsql.Raw(query, date, groupID).Scan(&installmentDetails).Error
	if err != nil {
		log.Println("#ERROR: ", err)
		return installmentDetails, errors.New("Unable to retrive installment details")
	}
	return installmentDetails, nil
}
