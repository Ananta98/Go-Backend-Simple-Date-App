package models

import "time"

type SwipeUserHistory struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	UserMatchID uint      `json:"user_match_id" gorm:"not null"`
	Like        bool      `json:"like" gorm:"default:false;not null"`
	SwapDate    time.Time `json:"swap_date" gorm:"not null"`

	// Relationship
	User      User `gorm:"foreignKey:UserID" json:"-"`
	UserMatch User `gorm:"foreignKey:UserMatchID" json:"-"`
}
