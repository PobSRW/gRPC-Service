package main

import (
	"log"
	"obp-gRPC-client/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	creds := insecure.NewCredentials()

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	calculatorClient := service.NewCalculatorClient(cc)
	calculatorService := service.NewCalculatorService(calculatorClient)

	err = calculatorService.Hello("Pob")
	if err != nil {
		log.Fatal(err)
	}
}
