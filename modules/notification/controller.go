package notification

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/notification"

	services.DomainStructSingle = Notification{}
	services.DomainStructArray = []Notification{}

	services.DBCPsql.AutoMigrate(&Notification{})

	notification := iris.Party(BASE_URL)
	{
		notification.Get("", services.Get)
		notification.Get("/get/:id", services.GetById)
		notification.Get("/q", services.GetByQuery)
		notification.Post("", services.Post)
		notification.Put("/set/:id", services.Put)
		notification.Delete("/delete/:id", services.Delete)
	}
}
