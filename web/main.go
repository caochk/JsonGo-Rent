package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
)

const ServerName = "go.micro.web.rent"

// 待测试
func main() {
	// Create web server
	srv := http.NewServer(
		server.Name(ServerName),
		server.Address(":8888"),
	)
	// Create service
	service := micro.NewService(micro.Server(srv))
	// Initiate service
	service.Init()
	// Run service
	if err := service.Run(); err != nil {
		fmt.Println("11", err)
	}
}
