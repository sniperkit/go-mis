package loanRaw

import "time"
import "encoding/json"
import "database/sql/driver"

type JSONB map[string]interface{}

type LoanRaw struct {
	ID        uint64      `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanID    uint64      `gorm:"column:loanId" json:"loanId"`
	Version   string      `gorm:"column:_version" json:"_version"`
	Raw       JSONB `gorm:"column:_raw" json:"_raw" sql:"type:jsonb"`
	CreatedAt time.Time   `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time   `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time  `gorm:"column:deletedAt" json:"deletedAt"`
}

type LoanRawNew struct {
	ID        uint64      `json:"id"`
	LoanID    uint64      `json:"loanId"`
	Version   string      `json:"_version"`
	Raw       string      `json:"_raw" sql:"type:jsonb"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt *time.Time  `json:"deletedAt"`
}

func (j JSONB) Value() (driver.Value, error) {
  	valueString, err := json.Marshal(j)
  	return string(valueString), err
  }
  
func (j *JSONB) Scan(value interface{}) error {
  	if err := json.Unmarshal(value.([]byte), &j); err != nil {
  		return err
  	}
 	return nil
 }
