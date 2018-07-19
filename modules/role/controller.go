package role

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(Role{}, []Role{})
}
