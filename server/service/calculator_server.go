package service

import (
	context "context"
	"fmt"
)

type calculatorServer struct{}

// concept plug and adapter
func NewCalculatorServer() CalculatorServer {
	return calculatorServer{}
}

func (calculatorServer) mustEmbedUnimplementedCalculatorServer() {}

func (calculatorServer) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	result := fmt.Sprintf("Hello %[1]s", req.Name)

	return &HelloResponse{
		Result: result,
	}, nil
}
