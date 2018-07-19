package profitAndLoss

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(ProfitAndLoss{}, []ProfitAndLoss{})
}
