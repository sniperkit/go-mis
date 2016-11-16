package notification

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Notification{})
	services.BaseCrudInit(Notification{}, []Notification{})
}
