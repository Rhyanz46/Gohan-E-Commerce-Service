package api

import (
	"fmt"
	"main/database"
	"main/settings"
)

func (server *Server) initializePrimaryDB() {
	server.DB = settings.DataSettings.DB.CreateConnection()

	//database migration
	err := server.DB.AutoMigrate(
		&database.User{},
		&database.ProductPhoto{},
		&database.Product{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
