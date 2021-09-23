package api

import (
	"main/api/product"
	"main/api/user"
)

func (server *Server) initializeRoutes() {
	route := server.Router.HandleFunc

	// admin router
	userRoute := user.Routes(&user.User{DB: server.DB})
	route(Endpoint.UserLogin, userRoute.HandleUserLogin)
	route(Endpoint.UserRegister, userRoute.HandleUserRegister)
	route(Endpoint.UserDetail, userRoute.HandleUserDetail)

	// product router
	chatRoute := product.Routes(&product.Product{DB: server.DB})
	route(Endpoint.Product, chatRoute.HandleProduct)
	route(Endpoint.ProductList, chatRoute.HandleProductList)
	route(Endpoint.ProductIdProduct, chatRoute.HandleUserDetail)
	route(Endpoint.ProductIdProductUploadProof, chatRoute.HandleUserDetail)

}
