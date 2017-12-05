package systemParameter

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

const (
	keyBackdate = "VT-BACKDATE"
)

func Init() {
	services.DBCPsql.AutoMigrate(&SystemParameter{})
	services.BaseCrudInit(SystemParameter{}, []SystemParameter{})
}

// FindByKey - Find system parameter by key
func FindByKey(key string) (SystemParameter, error) {
	var systemParameter SystemParameter
	if len(strings.Trim(key, " ")) == 0 {
		log.Println("#ERROR: Key param is empty")
		return systemParameter, errors.New("Key param can not be empty")
	}
	query := `select system_parameter.id,
				system_parameter."key",
				system_parameter.value
			from system_parameter
			where UPPER(system_parameter."key") = UPPER(?)`
	err := services.DBCPsql.Raw(query, key).Scan(&systemParameter).Error
	if err != nil {
		log.Println("#ERROR: Unable to retrive System Parameter by Key", err.Error())
		return systemParameter, errors.New("Unable to retrive System Parameter by Key")
	}
	return systemParameter, nil
}

// IsAllowedBackdate - Check if whether allow backdate or not
func IsAllowedBackdate() bool {
	systemParams, _ := FindByKey(keyBackdate)
	/*
		date, _ := misUtility.StringToDate(dateParam)
		if misUtility.IsBeforeToday(date) && (systemParams.ID == 0 || systemParams.Value == "FALSE") {
			return false
		}
	*/
	if strings.ToUpper(systemParams.Value) == "FALSE" {
		fmt.Println("Allow backdate: false")
		return false
	}
	fmt.Println("Allow backdate: true")
	return true
}

func GetSwift(ctx *iris.Context) {
	dataSwift := []Swift{}
	if err := services.DBCPsql.Table(`swift`).Find(&dataSwift).Error; err != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   dataSwift,
	})
}
