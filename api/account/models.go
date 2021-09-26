package account

type LoginData struct {
	Id       uint
	Username string
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterData struct {
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}
