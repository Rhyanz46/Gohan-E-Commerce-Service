package api

import (
	"main/api/account"
	"main/api/chat"
	"main/api/product"
	"main/settings"
	"net/http"
)

func (server *Server) initializeRoutes() {
	route := server.Router.HandleFunc

	// account router
	accountRoute := account.Routes(&account.Account{DB: server.DB})
	route(AdminEndpoint.Login, accountRoute.HandlerLogin)
	route(AdminEndpoint.Register, accountRoute.HandleRegister)
	route(AdminEndpoint.Detail, accountRoute.HandleDetail)

	// product router
	productRoute := product.Routes(&product.Product{DB: server.DB})
	route(AdminEndpoint.Product, productRoute.HandleProduct)
	route(AdminEndpoint.ProductDetail, productRoute.HandleProductDetail)
	route(AdminEndpoint.ProductDetailPhoto, productRoute.HandleProductDetailPhoto)
	route(AdminEndpoint.ProductDetailPhotoDetail, productRoute.HandleProductDetailPhotoDetail)

	// websocket router
	chatRoute := chat.Routes(&chat.Chat{DB: server.DB})
	go chatRoute.Socket().Run()
	route(AdminEndpoint.Ws, chatRoute.HandleChat)

	server.Router.PathPrefix(settings.ProductPhotosPrefixUrl).Handler(
		http.StripPrefix(settings.ProductPhotosPrefixUrl,
			http.FileServer(http.Dir("./"+settings.DataSettings.ProductPhotosFolder)),
		),
	)
}
