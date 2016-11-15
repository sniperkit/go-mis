package branch

import (
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Branch struct {
	ID        uint             `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Code      int64            `gorm:"column:code" json:"code"`
	Area      int64            `gorm:"column:area" json:"area"`
	Name      string           `gorm:"column:name" json:"name"`
	Location  string           `gorm:"column:location" json:"location"`
	Lat       float64          `gorm:"column:lat" json:"lat"`
	Lng       float64          `gorm:"column:lng" json:"lng"`
	Geopoint  gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}
