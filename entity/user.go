// yourproject/entity/user.go
package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique;not null"`
	Email          string `gorm:"unique;not null"`
	Number         string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	Role           string `gorm:"default:'user'"`
	Status         string `gorm:"default:'active'"`
	ReferalDetails ReferalDetails
	CreatedAt      time.Time
}
type ReferalDetails struct {
	gorm.Model
	UserID      uint   `json:"userID" gorm:"index;foreignKey:UserID"`
	Count       uint   `json:"count"`
	ReferalCode string `json:"referalCode" gorm:"unique_index"`
}
