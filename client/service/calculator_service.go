package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalculatorService interface {
	Hello(name string) error
	Fibonacci(n uint32) error
	Average(numbers ...float64) error
	Sum(numbers ...int32) error
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

func (c calculatorService) Sum(numbers ...int32) error {
	stream, err := c.calculatorClient.Sum(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Service Sum\n")

	// ประกอบด้วย 2 ส่วนคือ 1.send req 2.recv resp
	// แนะนำให้แยกออกจากกัน

	// 1.send req => ใช้ stream ส่งไปให้ server
	go func() {
		for _, number := range numbers {
			req := SumRequest{
				Number: number,
			}
			stream.Send(&req)
			fmt.Printf("Request: %v\n", req.Number)

			// หน่วงเวลาไว้ดูหน่อย
			time.Sleep(time.Second)
		}

		// พอหมด numbers ต้องบอก server ด้วยว่าจะไม่ส่งอะไรไปแล้วนะจ้ะ
		stream.CloseSend()
	}()

	done := make(chan bool)
	errs := make(chan error)

	// 2.recv res => ไว้รับค่า stream จาก server
	go func() {
		// ฝั่ง res ไม่รู้ว่าจบเมื่อไร จึงต้องวนลูปไม่รู้จบ
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errs <- err
			}
			fmt.Printf("Response %v\n", resp.Result)
			fmt.Println("--------------------------")
		}
		done <- true
	}()

	// ใช้ channel เพื่อให้รอให้ func ทำงานให้จบ
	select {
	case <-done:
		log.Println("Done")
		return nil
	case err := <-errs:
		return err
	}
}
