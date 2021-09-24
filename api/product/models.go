package product

type ProductData struct {
	ID          uint
	UserID      uint
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`

	EditData map[string]interface{}
}

type ProductPhotoData struct {
	ID     uint
	UserID uint
}
