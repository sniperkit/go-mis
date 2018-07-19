package voucher

import (
	"fmt"

	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
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

func CheckVoucherByOrderNo(orderNo string) Voucher {
	voucher := Voucher{}
	query := `select v.* from r_loan_order_voucher as rlov join voucher as v on rlov."voucherId" = v."id" join loan_order as lo on rlov."loanOrderId" = lo.id where rlov."deletedAt" isnull and lo."orderNo"='` + orderNo + `'`
	if err := services.DBCPsql.Raw(query).Scan(&voucher).Error; err != nil {
		fmt.Println(err)
	}
	return voucher
}
