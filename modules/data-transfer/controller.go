package dataTransfer

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&DataTransfer{})
	services.BaseCrudInit(DataTransfer{}, []DataTransfer{})
}

// get open validation teller date
func AvailableValidationTellerDate (ctx *iris.Context) {

    query := `with vd as (select to_char(installment."createdAt"::date, 'YYYY-MM-DD') as "createdAt"
    from installment 
    where installment."createdAt" >  CURRENT_DATE - INTERVAL '1 month'
    group by 1)
    select * from vd
    where vd."createdAt" not in (select "validationDate" from data_transfer)`

    vd := []ValidationDate{}

    services.DBCPsql.Raw(query).Find(&vd)

    ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"data":      vd,
	})
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
