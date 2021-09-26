package product

import "time"

type ProductData struct {
	ID          uint
	UserID      uint
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`

	EditData map[string]interface{}
}

type ProductPhotoData struct {
	ProductID uint
	UserID    uint
	PhotoID   uint
}

type ProductPhotoResponse struct {
	Id         uint      `json:"id"`
	Address    string    `json:"address"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
