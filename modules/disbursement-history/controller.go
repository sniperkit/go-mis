package disbursementHistory

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(DisbursementHistory{}, []DisbursementHistory{})
}
