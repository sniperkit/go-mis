package mitramanagement

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"

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

	// time
	t := time.Now().Format("2006-01-02 15:04:05")

	db := services.DBCPsql.Begin()
	// update installment
	q := `update Installment set "statusId" = ?, "reasonId" = ?, "updatedAt"=? where id=?`
	err = db.Exec(q, payload.StatusID, payload.ReasonID, t, payload.InstallmentID).Error
	if err != nil {
		ProcessErrorAndRollback(ctx, db, "Error Update Installment: "+err.Error())
		return
	}

	if payload.StatusID == 1 {
		q = `update borrower set "doDate" = ? where id = ?`
		err = services.DBCPsql.Exec(q, t, payload.BorrowerID).Error
		if err != nil {
			ProcessErrorAndRollback(ctx, db, "Error Update Borrower: "+err.Error())
			return
		}
	}

	db.Commit()

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   "installment data has ben updated",
	})
}

func ProcessErrorAndRollback(ctx *iris.Context, db *gorm.DB, message string) {
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status":  "error",
		"message": message,
	})
}
