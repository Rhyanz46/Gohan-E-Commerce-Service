package main

import (
	server2 "main/api"
	"main/settings"
)

var server = server2.Server{}

func main() {
	server.Initialize()
	server.Run(settings.DataSettings.Port)
}
