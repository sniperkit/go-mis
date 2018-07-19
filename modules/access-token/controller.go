package accessToken

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(AccessToken{}, []AccessToken{})
}
