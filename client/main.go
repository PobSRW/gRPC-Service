package main

import (
	"flag"
	"log"
	"obp-gRPC-client/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {

	var cc *grpc.ClientConn
	var err error
	var creds credentials.TransportCredentials

	host := flag.String("host", "localhost:50051", "gRPC Server host")
	tls := flag.Bool("tls", false, "use a secure TLS connection")
	flag.Parse()

	if *tls {
		cretFile := "../tls/ca.crt"
		creds, err = credentials.NewClientTLSFromFile(cretFile, "")
		if err != nil {
			log.Fatal(err)
		}

	} else {
		creds = insecure.NewCredentials()
	}

	cc, err = grpc.Dial(*host, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	calculatorClient := service.NewCalculatorClient(cc)
	calculatorService := service.NewCalculatorService(calculatorClient)

	err = calculatorService.Hello("Pob")
	// err = calculatorService.Fibonacci(7)
	// err = calculatorService.Average(1, 2, 3, 4, 6)
	// err = calculatorService.Sum(1, 2, 3, 4, 5)

	if err != nil {
		// แยกแยะ err ว่ามาจาก grpc หรือที่อื่นๆ จาก status.FromError
		if grpcErr, ok := status.FromError(err); ok {
			log.Printf("[%[1]v] %[2]v", grpcErr.Code(), grpcErr.Message())
		} else {
			log.Fatal(err)
		}
	}
}
