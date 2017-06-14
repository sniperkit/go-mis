package disbursementReport

import "time"

type DisbursementReport struct {
	ID               	uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	DisbursementDateFrom 	string     `gorm:"column:disbursementDateFrom" json:"disbursementDateFrom"`
	DisbursementDateTo 	string     `gorm:"column:disbursementDateTo" json:"disbursementDateTo"`
	Filename            	string     `gorm:"column:filename" json:"filename"`
	IsActive            	bool        `gorm:"column:isActive" json:"isActive"`
	CreatedAt        	time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        	time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        	*time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type DisbursementReportDetail struct {
	Dates	[]string     `json:"dates"`
	Details []DisbursementArea     `json:"details"`
}

type DisbursementArea struct {
	Name 	string	`json:"name"`
	Branchs []DisbursementBranch	`json:"branchs"`
	Prices 		[]float64	`json:"prices"`
	Total 		float64	`json:"total"`
}

type DisbursementBranch struct {
	Name 		string	`json:"name"`
	Prices 		[]float64	`json:"prices"`
	Total 		float64	`json:"total"`
}
