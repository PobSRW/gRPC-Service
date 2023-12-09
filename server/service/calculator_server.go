package service

import (
	context "context"
	"fmt"
	"time"
)

type calculatorServer struct{}

// concept plug and adapter
func NewCalculatorServer() CalculatorServer {
	return calculatorServer{}
}

func (calculatorServer) mustEmbedUnimplementedCalculatorServer() {}

func (calculatorServer) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	result := fmt.Sprintf("Hello %[1]s, at %[2]v", req.Name, req.CreatedDate.AsTime().Local())

	return &HelloResponse{
		Result: result,
	}, nil
}

func (calculatorServer) Fibonacci(req *FibonacciRequest, stream Calculator_FibonacciServer) error {
	for n := uint32(0); n <= req.N; n++ {
		result := fib(n)
		resp := FibonacciResponse{
			Result: result,
		}

		stream.Send(&resp)

		time.Sleep(time.Second)
	}
	return nil
}

func fib(n uint32) uint32 {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return fib(n-1) + fib(n-2)
	}
}
