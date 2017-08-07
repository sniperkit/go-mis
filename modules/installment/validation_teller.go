package installment

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	misConfig "bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

var (
	logAPIPath = misConfig.GoLogPath + "archive"
)

// Coh - Cash on hand struct
type Coh struct {
	installmentId uint64
	cash          float64
}

// TellerValidation struct
type TellerValidation struct {
	ID         string `json:"id"`
	CashOnHand []Coh
}

// Log struct
type Log struct {
	GroupID   string      `json:"groupId"`
	ArchiveID string      `json:"archiveId"`
	Data      interface{} `json:"data"`
}

func ValidationTeller(ctx *iris.Context) {
	var err error
	db := services.DBCPsql.Begin()
	validationTellerModel := TellerValidation{}
	err = ctx.ReadJSON(&validationTellerModel)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"errorMessage": "Bad Request",
			"message":      "Can not Unmarshall JSON Body",
		})
		return
	}
	var installment Installment
	date := ctx.Param("date")
	branchID := ctx.Param("branchID")
	installments, err := installment.FindByBranchAndDate(branchID, date)
	if err != nil {
		log.Println("#ERROR: ", err)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"errorMessage": "System Error",
			"message":      err.Error(),
		})
		return
	}
	// db.begin
	for _, installment = range installments {
		coh := getCOH(installment.ID, validationTellerModel.CashOnHand)
		if coh != -1 {
			err = UpdateStageAndCashOnHand(db, installment.ID, installment.Stage, "PENDING", coh)
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
	}
	db.Commit()
	go postToLog(getLog(branchID, validationTellerModel))
	ctx.JSON(iris.StatusOK, iris.Map{
		"message": "Success",
	})

}

// getCOH - filter cash on hand from client based on installment ID
func getCOH(instalmentID uint64, coh []Coh) float64 {
	for _, c := range coh {
		if c.installmentId == instalmentID {
			return c.cash
		}
	}
	return -1
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
