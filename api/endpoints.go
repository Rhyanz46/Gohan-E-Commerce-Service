package api

var Endpoint = struct {
	UserLogin    string
	UserRegister string
	UserDetail   string

	Product             string
	ProductDetail       string
	ProductDetailPhotos string
}{
	UserLogin:    "/user/login",
	UserRegister: "/user/register",
	UserDetail:   "/user/detail",

	Product:             "/product",
	ProductDetail:       "/product/{product_id:[0-9]+}",
	ProductDetailPhotos: "/product/{product_id:[0-9]+}/photos",
}
