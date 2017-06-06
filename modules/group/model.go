package group

import (
	"time"

	"github.com/nferruzzi/gormGIS"
)

type Group struct {
	ID           uint64           `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name         string           `gorm:"column:name" json:"name"`
	Lat          float64          `gorm:"column:lat" json:"lat"`
	Lng          float64          `gorm:"column:lng" json:"lng"`
	ScheduleDay  string           `gorm:"column:scheduleDay" json:"scheduleDay"`
	ScheduleTime string           `gorm:"column:scheduleTime" json:"scheduleTime"`
	Geopoint     gormGIS.GeoPoint `gorm:"column:geopoint" sql:"type:geometry(Geometry,4326)"`
	CreatedAt    time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time       `gorm:"column:deletedAt" json:"deletedAt"`
}



type GroupAgentBorrower struct{
	ID           uint64           `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name         string           `gorm:"column:name" json:"name"`
	Lat          float64          `gorm:"column:lat" json:"lat"`
	Lng          float64          `gorm:"column:lng" json:"lng"`
	ScheduleDay  string           `gorm:"column:scheduleDay" json:"scheduleDay"`
	ScheduleTime string           `gorm:"column:scheduleTime" json:"scheduleTime"`
	BorrowerName string 					`gorm:"column:borrowerName" json:"borrowerName"`
	Agent 				string           `gorm:"column:agentName" json:"agentName"`
	AgentId 				string           `gorm:"column:agentId" json:"agentId"`
}

type GroupBranchAreaAgent struct {
	ID        uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Branch    string    `gorm:"column:branch" json:"branch"`
	Area      string    `gorm:"column:area" json:"area"`
	Agent     string    `gorm:"column:agent" json:"agent"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
}

func (g *Group) SetGeopoint() {
	g.Geopoint = gormGIS.GeoPoint{
		Lat: g.Lat,
		Lng: g.Lng,
	}
}
