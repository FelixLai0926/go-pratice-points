package main

import (
	"points/pkg/routers"
)

func main() {
	server := routers.Setup()
	//server.Run(os.Getenv("server_host") + ":" + os.Getenv("server_port"))
	server.Run()
}
