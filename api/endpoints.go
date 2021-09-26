package api

var AdminEndpoint = struct {
	Login    string
	Register string
	Detail   string

	Product                  string
	ProductDetail            string
	ProductDetailPhoto       string
	ProductDetailPhotoDetail string

	Ws string
}{
	Login:    "/admin/login",
	Register: "/admin/register",
	Detail:   "/admin/detail",

	Product:                  "/admin/product",
	ProductDetail:            "/admin/product/{product_id:[0-9]+}",
	ProductDetailPhoto:       "/admin/product/{product_id:[0-9]+}/photo",
	ProductDetailPhotoDetail: "/admin/product/{product_id:[0-9]+}/photo/{photo_id:[0-9]+}",

	Ws: "/admin/ws",
}
