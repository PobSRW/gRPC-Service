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

	// err = calculatorService.Hello("Pob")
	// err = calculatorService.Fibonacci(7)
	err = calculatorService.Average(1, 2, 3, 4, 6)

	if err != nil {
		log.Fatal(err)
	}
}
