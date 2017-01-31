package branch

import (
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Branch struct {
	ID        uint64           `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Code      int64            `gorm:"column:code" json:"code"`
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

type BranchManager struct {
	ID           uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name         string `gorm:"column:name" json:"name"`
	City         string `gorm:"column:city" json:"city"`
	Province     string `gorm:"column:province" json:"province"`
	Address      string `gorm:"column:address" json:"address"`
	ManagaerName string `gorm:"column:managerName" gorm:"managerName"`
	AreaName     string `gorm:"column:areaName" gorm:"areaName"`
}

func (b *Branch) SetGeopoint() {
	b.Geopoint = gormGIS.GeoPoint{
		Lat: b.Lat,
		Lng: b.Lng,
	}
}
