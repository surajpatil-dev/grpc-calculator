package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/surajpatil-dev/grpc-calculator/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func doSum(cs calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{Input: &calculatorpb.Input{Number1: 10, Number2: 11}}

	res, err := cs.Sum(context.Background(), req)
	if err != nil {
		log.Fatalln("Sum Failed,", err)
	}

	log.Println("The result is - ", res.GetResult())
}

func getFactors(cs calculatorpb.CalculatorServiceClient) {
	number := int64(24)
	req := &calculatorpb.GetFactorRequest{Number: number}

	stream, err := cs.GetFactor(context.Background(), req)
	if err != nil {
		log.Fatalln("GetFactor Server call failed,", err)
	}
	log.Printf("Factors of %d are : ", number)

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		}
		if err != nil {
			log.Fatalln("Probelm with Stream,", err)
		}
		log.Print(res.GetNumber())

	}
}

func doAverage(cs calculatorpb.CalculatorServiceClient) {
	numbers := []int32{121, 119, 121, 123, 177}

	stream, err := cs.Average(context.Background())

	if err != nil {
		log.Fatalln("Error while calling Average,", err)
	}

	for _, number := range numbers {
		log.Println("Sending", number)
		stream.Send(&calculatorpb.AverageRequest{Number: number})
		time.Sleep(time.Millisecond * 1000)
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalln("Failed to close,", err)
	}

	log.Printf("The average of %v is %v.\n", numbers, res.GetResult())
}

func main() {
	log.Println("Starting Client...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalln("Failed to create client,", err)
	}

	defer cc.Close()

	cs := calculatorpb.NewCalculatorServiceClient(cc)

	// unary
	// doSum(cs)

	// server streaming
	// getFactors(cs)

	//client streaming
	doAverage(cs)
}
