package group

import (
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Group struct {
	ID        uint             `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name      string           `gorm:"column:name" json:"name"`
	Lat       float64          `gorm:"column:lat" json:"lat"`
	Lng       float64          `gorm:"column:lng" json:"lng"`
	Geopoint  gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}

func (g *Group) SetGeopoint() {
	g.Geopoint = gormGIS.GeoPoint{
		Lat: g.Lat,
		Lng: g.Lng,
	}
}
