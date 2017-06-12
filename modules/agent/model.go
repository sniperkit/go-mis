package agent

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Agent struct {
	ID              uint64           `gorm:"primary_key" gorm:"column:id" json:"_id"`
	Username        string           `gorm:"column:username" json:"username"`
	Password        string           `gorm:"column:password" json:"password"`
	Fullname        string           `gorm:"column:fullname" json:"fullname"`
	BankName        string           `gorm:"column:bankName" json:"bankName"`
	BankAccountName string           `gorm:"column:bankAccountName" json:"bankAccountName"`
	BankAccountNo   string           `gorm:"column:bankAccountNo" json:"bankAccountNo"`
	PicUrl          string           `gorm:"column:picUrl" json:"picUrl"`
	PhoneNo         string           `gorm:"column:phoneNo" json:"phoneNo"`
	Address         string           `gorm:"column:address" json:"address"`
	Kelurahan       string           `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan       string           `gorm:"column:kecamatan" json:"kecamatan"`
	City            string           `gorm:"column:city" json:"city"`
	Province        string           `gorm:"column:province" json:"province"`
	Nationality     string           `gorm:"column:nationality" json:"nationality"`
	Zipcode         string           `gorm:"column:zipCode" json:"zipcode"`
	Lat             float64          `gorm:"column:lat" json:"lat"`
	Lng             float64          `gorm:"column:lng" json:"lng"`
	Geopoint        gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt       time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}

type AgentBranch struct {
	ID              uint64           `gorm:"primary_key" gorm:"column:id" json:"_id"`
	Username        string           `gorm:"column:username" json:"username"`
	Password        string           `gorm:"column:password" json:"password"`
	Fullname        string           `gorm:"column:fullname" json:"fullname"`
	BankName        string           `gorm:"column:bankName" json:"bankName"`
	BankAccountName string           `gorm:"column:bankAccountName" json:"bankAccountName"`
	BankAccountNo   string           `gorm:"column:bankAccountNo" json:"bankAccountNo"`
	PicUrl          string           `gorm:"column:picUrl" json:"picUrl"`
	PhoneNo         string           `gorm:"column:phoneNo" json:"phoneNo"`
	Address         string           `gorm:"column:address" json:"address"`
	Kelurahan       string           `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan       string           `gorm:"column:kecamatan" json:"kecamatan"`
	City            string           `gorm:"column:city" json:"city"`
	Province        string           `gorm:"column:province" json:"province"`
	Nationality     string           `gorm:"column:nationality" json:"nationality"`
	Lat             float64          `gorm:"column:lat" json:"lat"`
	Lng             float64          `gorm:"column:lng" json:"lng"`
	Branch          string         	 `gorm:"column:branchName" json:"branchName"`
	Geopoint        gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt       time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}

type FragmentAgent struct {
	ID          uint64           `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Username    string           `gorm:"column:username" json:"username"`
	Fullname    string           `gorm:"column:fullname" json:"fullname"`
	PicUrl      string           `gorm:"column:picUrl" json:"picUrl"`
	PhoneNo     string           `gorm:"column:phoneNo" json:"phoneNo"`
	Address     string           `gorm:"column:address" json:"address"`
	Kelurahan   string           `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan   string           `gorm:"column:kecamatan" json:"kecamatan"`
	City        string           `gorm:"column:city" json:"city"`
	Province    string           `gorm:"column:province" json:"province"`
	Nationality string           `gorm:"column:nationality" json:"nationality"`
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

func (a *Agent) BeforeCreate() (err error) {
	if a.Password != "" {
		bytePassword := []byte(a.Password)
		sha256Bytes := sha256.Sum256(bytePassword)
		a.Password = hex.EncodeToString(sha256Bytes[:])
	}

	return
}

func (a *Agent) BeforeUpdate() (err error) {
	if a.Password != "" {
		bytePassword := []byte(a.Password)
		sha256Bytes := sha256.Sum256(bytePassword)
		a.Password = hex.EncodeToString(sha256Bytes[:])
	}

	return
}
