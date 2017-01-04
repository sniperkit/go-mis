package borrowerProspective

import "time"

type BorrowerProspective struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Raw       string     `gorm:"column:raw" json:"raw" sql:"json"`
	IsApprove bool       `gorm:"column:isApprove" json:"isApprove"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
