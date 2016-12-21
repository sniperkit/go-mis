package loan

import "time"

type Loan struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	SubmittedPlafond float64    `gorm:"column:submittedPlafond" json:"submittedPlafond"`
	CreditScoreText  string     `gorm:"column:creditScoreText" json:"creditScoreText"`
	CreditScoreValue float64    `gorm:"column:creditScoreValue" json:"creditScoreValue"`
	Tenor            float64    `gorm:"column:tenor" json:"tenor"`
	Rate             float64    `gorm:"column:rate" json:"rate"`
	Installment      float64    `gorm:"column:installment" json:"installment"`
	Plafond          float64    `gorm:"column:plafond" json:"plafond"`
	IsPublish        bool       `gorm:"column:isPublish" json:"isPublish"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
