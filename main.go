package main

import (
	"points/pkg/router"
)

func main() {
	server := router.Setup()
	//server.Run(os.Getenv("server_host") + ":" + os.Getenv("server_port"))
	server.Run()
}
