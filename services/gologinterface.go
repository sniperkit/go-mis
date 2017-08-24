package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/go-mis/config"
)

type Log struct {
	GroupID   string      `json:"groupId"`
	ArchiveID string      `json:"archiveId"`
	Data      interface{} `json:"data"`
}

type LogResponse struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Data    []Log  `json:"data"`
}

var (
	logAPIPath = config.GoLogPath
)

func GetLog(branchID uint64, data interface{}, groupId string) Log {
	var logger Log
	if branchID <= 0 {
		return logger
	}
	logger = Log{
		GroupID:   groupId,
		ArchiveID: GenerateArchiveID(branchID),
		Data:      data,
	}
	return logger
}

// PostToLog - POST data archive to GO-LOG API
func PostToLog(l Log) error {
	if l.Data == nil {
		return errors.New("Can not send empty data")
	}
	logBytes := new(bytes.Buffer)
	json.NewEncoder(logBytes).Encode(l)
	log.Println(logAPIPath)
	log.Println(l)
	_, err := http.Post(logAPIPath+"/api/v1/archive", "application/json; charset=utf-8", logBytes)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// GetNotes - Find data validation teller notes
// This method will call API GO-LOG APP
// Path: http://localhost:5500/api/v1/archive-list/group/{groupId}
func GetNotes(groupID string) ([]Log, error) {
	var err error
	var logResp LogResponse
	if len(strings.Trim(groupID, " ")) == 0 {
		return nil, errors.New("Group ID can not be empty")
	}
	apiPath := logAPIPath + "/api/v1/archive-list/group/" + groupID
	log.Println(apiPath)
	resp, err := http.Get(apiPath)

	// In case of GO-LOG App not running
	if resp == nil {
		log.Println("Unable to get response form GO-LOG App")
		return nil, errors.New("Unable to get response form GO-LOG App")
	}

	if resp.StatusCode != 200 || err != nil {
		return nil, errors.New("Unable to retrieve data notes from GO-LOG APP")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("#ERROR: Unable to read body response when calling API Notes")
		log.Println("#ERROR: ", err)
		return nil, errors.New("Unable to retrieve data notes from GO-LOG APP")
	}
	err = json.Unmarshal([]byte(body), &logResp)
	if err != nil {
		log.Println("#ERROR: Unable to unmarshall body")
		log.Println("#ERROR: ", err)
		return nil, errors.New("Unable to unmarshall body")
	}
	logNotes := logResp.Data
	return logNotes, nil
}

// GetDataFromLog - Retrive data from GO-LOG App
func GetDataFromLog(branchID uint64) (Log, error) {
	var logger Log
	archiveID := GenerateArchiveID(branchID)
	groupID := "VALIDATION TELLER"
	apiPath := GetLogAPIPath() + "/api/v1/archive/" + archiveID + "/group/" + groupID
	log.Println("[INFO]", apiPath)
	resp, err := http.Get(apiPath)

	// In case of GO-LOG App not running
	if resp == nil {
		log.Println("Unable to get response from GO-LOG App")
		return logger, errors.New("Unable to get response from GO-LOG App")
	}
	if err != nil {
		log.Println("#ERROR: Unable to retrive data from GO-LOG App")
		log.Println("#ERROR: ", err)
		return logger, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("#ERROR: When read body reponse from GO-LOG App")
		return logger, err
	}
	err = json.Unmarshal([]byte(body), &logger)
	if err != nil {
		log.Println("#ERROR: When unmarshall resp body GO-LOG App to Log struct")
		return logger, err
	}
	return logger, nil
}

func GetRejectNotesData(status string, groupId string, date string, stage string) (Log, error) {
	logger := Log{}

	// 1-2017-08-18VTRejectIN-REVIEW
	logGroupId := groupId + "-" + date + "VT" + strings.Title(status) + strings.ToUpper(stage)
	archiveId := stage

	apiPath := GetLogAPIPath() + "/api/v1/archive/" + archiveId + "/group/" + logGroupId
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
		log.Println("#ERROR: When read body reponse from GO-LOG App")
		return logger, err
	}
	err = json.Unmarshal([]byte(body), &logger)
	if err != nil {
		log.Println("#ERROR: When unmarshall resp body GO-LOG App to Log struct")
		return logger, err
	}
	log.Println("[INFO]", logger)
	return logger, nil
}

func GenerateArchiveID(branchID uint64) string {
	if branchID == 0 {
		return ""
	}
	branchIDStr := string(branchID)
	return branchIDStr + "-" + time.Now().Local().Format("2006-01-02")
}

// GetLogAPIPath base path of GO-LOG APP API
func GetLogAPIPath() string {
	return logAPIPath
}

func ConstructNotesGroupId(branchId uint64, date string) string {
	log.Println("[INFO]", branchId)
	log.Println("[INFO]", date)
	groupID := strconv.FormatUint(branchId, 10) + "-" + date + "-VTNotes"
	return groupID
}

func ConstructRejectsNotesGroupId(groupId uint64, date string, status string, stage string) string {
	status = strings.Title(status)
	groupID := strconv.FormatUint(groupId, 10) + "-" + date + "VT" + status + strings.ToUpper(stage)
	return groupID
}

func GetBorrowerNotes(logNotes []Log) (borrowerNotes interface{}) {
	if len(logNotes) == 0 {
		return
	}
	for _, note := range logNotes {
		if strings.ToLower(note.ArchiveID) == "borrower" {
			borrowerNotes = note.Data
		}
	}
	return
}

func GetMajelisNotes(logNotes []Log) (majelisNotes interface{}) {
	if len(logNotes) == 0 {
		return
	}
	for _, note := range logNotes {
		if strings.ToLower(note.ArchiveID) == "majelis" {
			majelisNotes = note.Data
		}
	}
	return
}
