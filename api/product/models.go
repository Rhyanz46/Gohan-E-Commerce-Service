package product

type ProductData struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	UserID      uint
}

type RegisterData struct {
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}
