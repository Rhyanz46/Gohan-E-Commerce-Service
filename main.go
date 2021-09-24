package main

import (
	serverApi "main/api"
	"main/settings"
)

var server = serverApi.Server{}

func main() {
	server.Initialize()
	server.Run(settings.DataSettings.Port)
}
