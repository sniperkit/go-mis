package area

import (
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Area struct {
	ID        uint64           `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name      string           `gorm:"column:name" json:"name"`
	City      string           `gorm:"column:city" json:"city"`
	Province  string           `gorm:"column:province" json:"province"`
	Lat       float64          `gorm:"column:lat" json:"lat"`
	Lng       float64          `gorm:"column:lng" json:"lng"`
	Geopoint  gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}

func (a *Area) SetGeopoint() {
	a.Geopoint = gormGIS.GeoPoint{
		Lat: a.Lat,
		Lng: a.Lng,
	}
}
