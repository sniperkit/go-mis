package virtualAccount

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(VirtualAccount{}, []VirtualAccount{})
}
