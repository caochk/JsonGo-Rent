package main

import (
	"GetArea/handler"
	pb "GetArea/proto"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"

	"github.com/asim/go-micro/v3/logger"
)

const ServerName = "go.micro.srv.GetArea" // 微服务名

func main() {
	reg := consul.NewRegistry()
	// Create service
	service := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
		micro.Registry(reg),
	)

	// Initialise service
	service.Init()

	// Register handler
	if err := pb.RegisterGetAreaHandler(service.Server(), new(handler.GetArea)); err != nil {
		logger.Fatal(err)
	}
}
