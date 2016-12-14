package agent

import (
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Agent struct {
	ID          uint             `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Username    string           `gorm:"column:username" json:"username"`
	Password    string           `gorm:"column:password" json:"password"`
	Fullname    string           `gorm:"column:fullname" json:"fullname"`
	PicUrl      string           `gorm:"column:picUrl" json:"picUrl"`
	PhoneNo     string           `gorm:"column:phoneNo" json:"phoneNo"`
	Address     string           `gorm:"column:address" json:"address"`
	Kelurahan   string           `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan   string           `gorm:"column:kecamatan" json:"kecamatan"`
	City        string           `gorm:"column:city" json:"city"`
	Province    string           `gorm:"column:province" json:"province"`
	Nationality string           `gorm:"column:nationality" json:"nationality"`
	Zipcode     string           `gorm:"column:zipCode" json:"zipcode"`
	Lat         float64          `gorm:"column:lat" json:"lat"`
	Lng         float64          `gorm:"column:lng" json:"lng"`
	Geopoint    gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt   time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}

func (a *Agent) SetGeopoint() {
	a.Geopoint = gormGIS.GeoPoint{
		Lat: a.Lat,
		Lng: a.Lng,
	}
}
