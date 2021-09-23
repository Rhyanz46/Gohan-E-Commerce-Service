package database

import "time"

type Product struct {
	Id          uint      `gorm:"primaryKey;" json:"id"`
	Name        string    `gorm:"size:200;" json:"name"`
	Price       float64   `gorm:"column:price;" json:"price"`
	Description string    `gorm:"column:description;" json:"description"`
	CreateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"update_time"`

	ProductPhotos []ProductPhoto `gorm:"foreignKey:ProductId;references:Id" json:"-"`
	UserId        uint           `gorm:"foreignKey:Id;" json:"-"`
	//ProductPhotos []ProductPhoto `gorm:"foreignKey:Id"`
}

type ProductPhoto struct {
	Id         uint      `gorm:"primaryKey;" json:"id"`
	Address    string    `gorm:"size:500;" json:"address"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`

	ProductId uint `gorm:"foreignKey:Id;"`
}
