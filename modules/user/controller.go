package user

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&UserMis{})
	services.BaseCrudInit(UserMis{}, []UserMis{})
}
