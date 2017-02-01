package notification

import (
	"strings"

	"bitbucket.org/go-mis/modules/borrower"
	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Notification{})
	services.BaseCrudInit(Notification{}, []Notification{})
}

func SendPush(ctx *iris.Context) {
	notificationInput := NotificationInput{}

	err := ctx.ReadJSON(&notificationInput)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if strings.EqualFold(notificationInput.SentTo, "investor") {
		investors := []investor.Investor{}
		services.DBCPsql.Find(&investors)

		for _, investor := range investors {
			go insertPushInvestor(investor.ID, notificationInput)
		}

	} else {
		borrowers := []borrower.Borrower{}
		services.DBCPsql.Find(&borrowers)

		for _, borrower := range borrowers {
			go insertPushBorrower(borrower.ID, notificationInput)
		}

	}

	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status": "success",
		"data":   notificationInput,
	})
}

func insertPushInvestor(investorID uint64, notificationInput NotificationInput) {
	notification := Notification{Type: "info", Message: notificationInput.Message, IsRead: false, RedirectUrl: notificationInput.RedirectUrl}
	services.DBCPsql.Create(&notification)

	notificationInvestorData := r.RNotificationInvestor{NotificationId: notification.ID, InvestorId: investorID}
	services.DBCPsql.Create(&notificationInvestorData)

}

func insertPushBorrower(borrowerID uint64, notificationInput NotificationInput) {
	notification := Notification{Type: "info", Message: notificationInput.Message, IsRead: false, RedirectUrl: notificationInput.RedirectUrl}
	services.DBCPsql.Create(&notification)

	notificationBorrowerData := r.RNotificationBorrower{NotificationId: notification.ID, BorrowerId: borrowerID}
	services.DBCPsql.Create(&notificationBorrowerData)

}
