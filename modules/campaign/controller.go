package campaign

import (
	"fmt"

	"bitbucket.org/go-mis/services"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Campaign{})
	services.BaseCrudInit(Campaign{}, []Campaign{})
}

//CampaignItem - this will store quantity of campaign item
type CampaignItem struct {
	QuantityOfCampaignItem uint64 `gorm:"column:quantityOfCampaignItem" json:"quantityOfCampaignItem"`
}

//GetActiveCampaign - checking an active campaign
func GetActiveCampaignByOrderNo(orderNo string) (uint64, Campaign) {
	campaign := Campaign{}
	item := CampaignItem{}
	query := `select c.* , rloc.quantity as "quantityOfCampaignItem" from r_loan_order_campaign as rloc
	join campaign as c on rloc."campaignId" = c."id" 
	join loan_order as lo on rloc."loanOrderId" = lo.id 
	where rloc."deletedAt" isnull and c."isActive" = true and lo."orderNo"='` + orderNo + `'`
	if err := services.DBCPsql.Raw(query).Scan(&campaign).Error; err != nil {
		fmt.Println(err)
	}

	services.DBCPsql.Raw(query).Scan(&item)

	return item.QuantityOfCampaignItem, campaign
}
