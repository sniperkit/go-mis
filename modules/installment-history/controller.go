package installmentHistory

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(InstallmentHistory{}, []InstallmentHistory{})
}
