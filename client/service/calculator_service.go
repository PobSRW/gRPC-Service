package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalculatorService interface {
	Hello(name string) error
	Fibonacci(n uint32) error
	Average(number ...float64) error
}

// service ตัวนี้ทำงานเองไม่ได้
// จำเป็นต้องใช้ calculatorClient เข้ามาช่วย
type calculatorService struct {
	calculatorClient CalculatorClient
}

func NewCalculatorService(calculatorClient CalculatorClient) CalculatorService {
	return calculatorService{calculatorClient}
}

func (c calculatorService) Hello(name string) error {
	req := HelloRequest{
		Name:        name,
		CreatedDate: timestamppb.Now(),
	}

	// ใช้ calculatorClient เพื่อส่ง msg ไปให้กับ server และเพื่อรอ resp กลับมา
	res, err := c.calculatorClient.Hello(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service Hello\n")
	fmt.Printf("Request: %v\n", req.Name)
	fmt.Printf("Response: %v\n", res.Result)

	return nil
}

func (c calculatorService) Fibonacci(n uint32) error {
	req := FibonacciRequest{
		N: n,
	}

	// ถ้า stream ที่วิ่งเข้ามาหาเรามันนานเกินไป สามารถตัด connection ทิ้งได้
	// โดยกำหนด timeout ในการ call service

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := c.calculatorClient.Fibonacci(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service Fibonacci\n")
	fmt.Printf("Request: %v\n", req.N)

	// วน loop ไม่รู้จบ เพื่อรับ response
	for {
		resp, err := stream.Recv()

		// ถ้าฝั่ง server stream เสร็จแล้ว จะ retrun err มาเป็น io.EOF
		// EOF = end of file
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Response %v\n", resp.Result)
	}

	return nil
}

func (c calculatorService) Average(numbers ...float64) error {

	stream, err := c.calculatorClient.Average(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Service Average\n")
	for _, number := range numbers {
		req := AverageRequest{
			Number: number,
		}
		stream.Send(&req)
		fmt.Printf("Request: %v\n", req.Number)
		time.Sleep(time.Second)
	}

	//ตอนจบแล้วต้องรับ resp กลับมาด้วย
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("Response %v\n", resp.Result)

	return nil
}
