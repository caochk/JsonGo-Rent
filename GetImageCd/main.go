package main

import (
	"GetImageCd/handler"
	pb "GetImageCd/proto"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
)

const ServerName = "go.micro.srv.GetImageCd" // server name

func main() {
	reg := consul.NewRegistry()
	srv := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
		micro.Registry(reg),
	)
	// Initialise service
	srv.Init()

	// Register handler
	if err := pb.RegisterGetImageCdHandler(srv.Server(), new(handler.GetImageCd)); err != nil {
		logger.Fatal(err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
