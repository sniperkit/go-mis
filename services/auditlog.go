package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"bitbucket.org/go-mis/config"
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
		ArchiveID: generateArchiveID(branchID),
		Data:      data,
	}
	return logger
}

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

func GenerateArchiveID(branchID int64) string {
	if branchID == 0 {
		return ""
	}
	return string(branchID) + "-" + time.Now().Local().Format("2006-01-02")
}
