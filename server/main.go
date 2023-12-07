package main

import (
	"fmt"
	"log"
	"net"
	"obp-gRPC/service"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	// regiser calculator service
	service.RegisterCalculatorServer(s, service.NewCalculatorServer())

	fmt.Println("gRPC server is listening on port 50051")

	// run server
	err = s.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
