package service

import (
	context "context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type calculatorServer struct{}

// concept plug and adapter
func NewCalculatorServer() CalculatorServer {
	return calculatorServer{}
}

func (calculatorServer) mustEmbedUnimplementedCalculatorServer() {}

func (calculatorServer) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {

	// handle error
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

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

func (calculatorServer) Average(stream Calculator_AverageServer) error {
	var sum, count float64

	// วนแบบ infinity เพราะไม่รู้ว่า client จะ stream มาเท่าไร
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sum += req.Number
		count++
	}

	result := sum / count
	resp := AverageResponse{
		Result: result,
	}

	// SendAndClose = ส่งข้อความและปิด connection
	return stream.SendAndClose(&resp)
}

func (calculatorServer) Sum(stream Calculator_SumServer) error {
	var sum int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// คำนวณแล้วตรวจ resp ทันที เพราะเป็นแบบ bi directional
		sum += req.Number
		resp := SumResponse{
			Result: sum,
		}

		err = stream.Send(&resp)
		if err != nil {
			return err
		}
	}

	return nil
}
