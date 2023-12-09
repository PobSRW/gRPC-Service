package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalculatorService interface {
	Hello(name string) error
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
