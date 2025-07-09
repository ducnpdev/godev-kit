package models

import (
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
)

// ShipperLocationModel represents the DB model for shipper location
// TableName: shipper_locations
// (You may need to create a migration for this table)
type ShipperLocationModel struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ShipperID string    `json:"shipper_id" gorm:"column:shipper_id;index;not null"`
	Latitude  float64   `json:"latitude" gorm:"column:latitude;not null"`
	Longitude float64   `json:"longitude" gorm:"column:longitude;not null"`
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (ShipperLocationModel) TableName() string {
	return "shipper_locations"
}

func ToShipperLocationEntity(model ShipperLocationModel) entity.ShipperLocation {
	return entity.ShipperLocation{
		ShipperID: model.ShipperID,
		Latitude:  model.Latitude,
		Longitude: model.Longitude,
		Timestamp: model.Timestamp,
	}
}
