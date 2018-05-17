package dataTransfer

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&DataTransfer{})
	services.BaseCrudInit(DataTransfer{}, []DataTransfer{})
}

func Save(ctx *iris.Context) {
	var payload DataTransfers

	err := ctx.ReadJSON(&payload)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}

	// Find data transfer by validation date
   for _, val := range payload.Items {
		// Create new data transfer
		if err := services.DBCPsql.Create(&val).Error; err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"message":      "Create Error",
				"errorMessage": err.Error(),
			})
			return
		}
    }
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": payload})
}
