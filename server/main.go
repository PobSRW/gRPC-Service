package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"obp-gRPC/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var s *grpc.Server
	tls := flag.Bool("tls", false, "use a secure TLS connection")

	// หลังใช้ flag แล้วต้องมี .Parse() ด้วยไม่งั้นจะใช้งานไม่ได้
	flag.Parse()

	if *tls {
		certFile := "../tls/server.crt"
		keyFile := "../tls/server.pem"

		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}
		// start server แบบ tls (secure โดยใช้ cert)
		s = grpc.NewServer(grpc.Creds(creds))

	} else {
		// start server แบบ insecure (ส่งข้อมูลแบบไม่ปลอดภัย)
		s = grpc.NewServer()
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	// regiser calculator service
	service.RegisterCalculatorServer(s, service.NewCalculatorServer())

	fmt.Print("gRPC server is listening on port 50051")
	if *tls {
		fmt.Println(" with TLS")
	}

	// run server
	err = s.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
