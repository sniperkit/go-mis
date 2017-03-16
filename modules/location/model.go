package location

type Location struct {
	LocationCode string `gorm:"column:locationCode" json:"locationCode"`
	Name         string `gorm:"column:name" json:"name"`
	Province     string `gorm:"column:province" json:"province"`
	City         string `gorm:"column:city" json:"city"`
	Kecamatan    string `gorm:"column:kecamatan" json:"kecamatan"`
	Kelurahan    string `gorm:"column:kelurahan" json:"kelurahan"`
}
