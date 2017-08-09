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
	"bitbucket.org/go-mis/services"
)

type Log struct {
	GroupID   string      `json:"groupId"`
	ArchiveID string      `json:"archiveId"`
	Data      interface{} `json:"data"`
}

var (
	logAPIPath = config.GoLogPath
)

func GetLog(branchID int64, data interface{}, groupId string) Log {
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
	_, err := http.Post(logAPIPath+"archive", "application/json; charset=utf-8", logBytes)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// FindNotes - Find data validation teller notes
// This method will call API GO-LOG APP
// Path: http://localhost:5500/api/v1/archive-list/group/{groupId}
func FindNotes(groupID string) ([]Log, error) {
	var err error
	var notes []Log
	if len(strings.Trim(groupID, " ")) == 0 {
		return nil, errors.New("Group ID can not be empty")
	}
	resp, err := http.Get(logAPIPath + "archive-list/group/" + groupID)
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
	err = json.Unmarshal([]byte(body), &notes)
	if err != nil {
		log.Println("#ERROR: Unable to unmarshall body")
		log.Println("#ERROR: ", err)
		return nil, errors.New("Unable to unmarshall body")
	}
	return notes, nil
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

func GenerateArchiveID(branchID int64) string {
	if branchID == 0 {
		return ""
	}
	branchIDStr := strconv.FormatInt(branchID, 10)
	return branchIDStr + "-" + time.Now().Local().Format("2006-01-02")
}

// GetLogAPIPath base path of GO-LOG APP API
func GetLogAPIPath() string {
	return logAPIPath
}
