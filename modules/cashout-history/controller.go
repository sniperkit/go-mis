package cashoutHistory

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(CashoutHistory{}, []CashoutHistory{})
}
