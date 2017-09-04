package mitramanagement

import (
	"log"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Status{})
	services.DBCPsql.AutoMigrate(&Reason{})
}

// GetPortfolioAtRisk - Get loan data where stage is equal to 'PAR'
func GetPortfolioAtRisk(ctx *iris.Context) {

}

func FindPortfolioAtRisk(parList []PortfolioAtRisk) (uint64, error) {
	return 0, nil
}

func GetStatusAll(ctx *iris.Context) {
	s := []Status{}
	err := services.DBCPsql.Where("type = 'mitra_management'").Find(&s).Error
	if err != nil {
		log.Println("[INFO] Params is not valid")
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   s,
	})
}

func SubmitReason(ctx *iris.Context) {
	payload := struct {
		InstallmentID uint64 `json:"installmentId"`
		BorrowerID    uint64 `json:"borrowerId"`
		Date          string `json:"date"`
		StatusID      uint64 `json:"statusId"`
		ReasonID      uint64 `json:"reasonId"`
	}{}

	err := ctx.ReadJSON(&payload)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// update installment
	q := `update Installment set status_id = ?, reason_id = ?, "updatedAt"=? where id=?`
	err = services.DBCPsql.Exec(q, payload.StatusID, payload.ReasonID, payload.Date, payload.InstallmentID).Error
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if payload.StatusID == 1 {
		q = `update borrower set "doDate" = ? where id = ?`
		err = services.DBCPsql.Exec(q, payload.Date, payload.BorrowerID).Error
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   "installment data has ben updated",
	})
}
