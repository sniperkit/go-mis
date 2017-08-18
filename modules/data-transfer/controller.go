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
	m := DataTransfer{}
	err := ctx.ReadJSON(&m)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}

	dt := DataTransfer{}
	services.DBCPsql.Where("validation_date = ?", m.ValidationDate).Find(&dt)

	if (DataTransfer{}) != dt {
		err := services.DBCPsql.Model(&m).Updates(&m).Error
		if err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"message":      "Update Error",
				"errorMessage": err.Error(),
			})
			return
		}
	} else {
		err = services.DBCPsql.Create(&m).Error
		if err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"message":      "Create Error",
				"errorMessage": err.Error(),
			})
			return
		}
	}

	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})

}
