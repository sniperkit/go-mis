package voucher

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Voucher{})
	services.BaseCrudInit(Voucher{}, []Voucher{})
}

func FetchAll(ctx *iris.Context) {
	voucher := []Voucher{}

	query := "SELECT \"ID\", \"amount\", \"voucherNo\", \"description\", \"startDate\", \"endDate\", \"isPersonal\" "
	query += "FROM voucher "
	
	if e := services.DBCPsql.Raw(query).Find(&voucher).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
			"data":   e,
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   voucher,
	})
}
