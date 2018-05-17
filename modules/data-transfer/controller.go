package dataTransfer

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"

    "fmt"
)

func Init() {
	services.DBCPsql.AutoMigrate(&DataTransfer{})
	services.BaseCrudInit(DataTransfer{}, []DataTransfer{})
}

func Save(ctx *iris.Context) {
	var payload DataTransfers
	// var dt DataTransfer

	err := ctx.ReadJSON(&payload)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}

    fmt.Println(payload)
    /*
	// Find data transfer by validation date
	services.DBCPsql.Where("validation_date = ?", payload.ValidationDate).Find(&dt)
	if dt.ID > 0 {
		// Update existing data transfer
		payload.ID = dt.ID
		err := services.DBCPsql.Save(&payload).Error
		if err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"message":      "Update Error",
				"errorMessage": err.Error(),
			})
			return
		}
	} else {
		// Create new data transfer
		if err := services.DBCPsql.Create(&payload).Error; err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"message":      "Create Error",
				"errorMessage": err.Error(),
			})
			return
		}
	}
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": payload})
    */
}
