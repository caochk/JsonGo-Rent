package main

import (
	"GetSmsCd/handler"
	pb "GetSmsCd/proto"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/micro/micro/v3/service/logger"
)

const ServerName = "go.micro.srv.GetSmsCd"

func main() {
	reg := consul.NewRegistry()
	// Create service
	srv := micro.NewService(
		micro.Registry(reg),
		micro.Name("getsmscd"),
		micro.Name("latest"),
	)

	// Initialise service
	srv.Init()

	// Register handler
	pb.RegisterGetSmsCdHandler(srv.Server(), new(handler.GetSmsCd))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
