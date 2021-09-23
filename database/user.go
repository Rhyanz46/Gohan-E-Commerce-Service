package database

import "time"

type User struct {
	Id         uint      `gorm:"primaryKey;"`
	Email      string    `gorm:"unique;size:100;index:idx_user_email;"`
	FullName   string    `gorm:"size:100;"`
	Username   string    `gorm:"unique;size:200;index:idx_user_username;"`
	Password   string    `gorm:"column:password;"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	//Products   []Product `gorm:"foreignKey:Id"`
	Products []Product `gorm:"foreignKey:UserId;references:Id"`
}
