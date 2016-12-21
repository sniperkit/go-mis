package installmentPresence

import "time"

type InstallmentPresence struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type      string     `gorm:"column:type" json:"type"` // type: [ 'A', 'TR1', 'TR2', 'TR3', 'TA' ]
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
