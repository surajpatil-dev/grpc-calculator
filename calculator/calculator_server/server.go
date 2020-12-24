package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/grpc-calculator/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(cxt context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum invoked...")
	number1 := req.GetInput().GetNumber1()
	number2 := req.GetInput().GetNumber2()
	return &calculatorpb.SumResponse{Result: number1 + number2}, nil
}

func (*server) GetFactor(req *calculatorpb.GetFactorRequest, stream calculatorpb.CalculatorService_GetFactorServer) error {
	log.Println("Get Fator was invoked.")
	N := int(req.GetNumber())
	k := 2
	for {
		if N%k == 0 {
			stream.Send(&calculatorpb.GetFactorResponse{Number: int64(k)})
			N = N / k
			log.Println(k, N)
			time.Sleep(time.Millisecond * 1000)
		} else {
			k = k + 1
		}

		if N < k {
			break
		}
	}

	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	numbers := []int32{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// end of stream
			sum := int32(0)

			for _, number := range numbers {
				sum += number
			}
			average := float64(sum) / float64(len(numbers))
			return stream.SendAndClose(&calculatorpb.AverageResponse{Result: average})
		}

		if err != nil {
			log.Println(numbers)
			log.Fatalln("error while reciving stream", err)
		}
		log.Println("recieved ", req.GetNumber())
		numbers = append(numbers, req.GetNumber())

	}

}

func main() {
	log.Print("Server listening on 50051...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to bind listener,", err)
	}

}
