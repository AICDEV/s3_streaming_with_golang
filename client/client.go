package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/aicdev/s3_streaming_with_golang/proto"
	"google.golang.org/grpc"
)

func main() {
	con, err := grpc.Dial("localhost:9999", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("connection error: %v", err)
	}

	defer con.Close()

	c := pb.NewTransactionServiceClient(con)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := c.GetUserTransactions(ctx, &pb.User{
		Id: "acb03ff98e576ec4b5437d791d85d98c523821ecb61457668928f6b4c204ab944f5be2ca3605359d3a41aceb4946151a56fa16fdc09161e8a65129c19b4da83b",
	})

	if err != nil {
		log.Fatalf("unable to get data: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}
		log.Printf("Resp received: %s", resp)
	}
}
