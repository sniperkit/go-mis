package branch

import (
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Branch struct {
	ID            uint64           `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Code          int64            `gorm:"column:code" json:"code"`
	Name          string           `gorm:"column:name" json:"name"`
	City          string           `gorm:"column:city" json:"city"`
	AddressDetail string           `gorm:"column:addressDetail" json:"addressDetail"`
	Province      string           `gorm:"column:province" json:"province"`
	Lat           float64          `gorm:"column:lat" json:"lat"`
	Lng           float64          `gorm:"column:lng" json:"lng"`
	BranchNewCode string           `gorm:"column:branchNewCode" json:"branchNewCode"`
	Geopoint      gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt     time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}

type BranchByArea struct {
	ID     uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Area   string `gorm:"column:name" json:"area"`
	Branch string `gorm:"column:name" json:"branch"`
}

type BranchManagerArea struct {
	ID            uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name          string `gorm:"column:name" json:"name"`
	City          string `gorm:"column:city" json:"city"`
	Province      string `gorm:"column:province" json:"province"`
	Address       string `gorm:"column:address" json:"address"`
	Manager       string `gorm:"column:manager" json:"manager"`
	Area          string `gorm:"column:area" json:"area"`
	Role          string `gorm:"column:role" json:"role"`
	AddressDetail string `gorm:"column:addressDetail" json:"addressDetail"`
	BranchNewCode string `gorm:"column:branchNewCode" json:"branchNewCode"`
}

type BranchAreaManager struct {
	ID      uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name    string `gorm:"column:name" json:"name"`
	Manager string `gorm:"column:manager" json:"manager"`
	Area    string `gorm:"column:area" json:"area"`
	Role    string `gorm:"column:role" json:"role"`
}

type BranchManager struct {
	BranchId uint64 `gorm:"column:branchId"`
	Fullname string `gorm:"column:fullname"`
}

func (b *Branch) SetGeopoint() {
	b.Geopoint = gormGIS.GeoPoint{
		Lat: b.Lat,
		Lng: b.Lng,
	}
}
