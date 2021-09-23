package api

var Endpoint = struct {
	UserLogin    string
	UserRegister string
	UserDetail   string

	Product                     string
	ProductList                 string
	ProductIdProduct            string
	ProductIdProductUploadProof string
}{
	UserLogin:    "/user/login",
	UserRegister: "/user/register",
	UserDetail:   "/user/detail",

	Product:                     "/product",
	ProductList:                 "/product/list",
	ProductIdProduct:            "/product/{id_product}",
	ProductIdProductUploadProof: "/product/{id_product}/upload-proof",
}
