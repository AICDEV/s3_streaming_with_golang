package app

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aicdev/s3_streaming_with_golang/data"
	"github.com/aicdev/s3_streaming_with_golang/env"
	services "github.com/aicdev/s3_streaming_with_golang/streaming"
	"github.com/aicdev/s3_streaming_with_golang/uploader"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	pb "github.com/aicdev/s3_streaming_with_golang/proto"
)

var (
	mode = kingpin.Flag("mode", "mode to start application").Default("stream").String()
)

type server struct {
	pb.UnimplementedTransactionServiceServer
}

func init() {
	_, err := env.ParseEnvironmentVariables()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (s server) GetUserTransactions(req *pb.User, stream pb.TransactionService_GetUserTransactionsServer) error {
	c := make(chan *pb.TransactionCollection)
	go services.StreamingService.Stream(req.GetId(), c)
	for transactions := range c {
		for _, transaction := range transactions.GetTransactions() {
			stream.Send(transaction)
		}
	}

	return nil
}

func StartApplication() {
	kingpin.Parse()
	log.Printf("starting streaming service in mode: %s", *mode)
	switch *mode {
	case "test":
		runTestDataCreation()

	case "stream":
		bootRPCService()

	default:
		log.Fatalf("unknown mode: %s", *mode)
	}
}

func runTestDataCreation() {

	uploader := uploader.NewUploaderService()

	testDataCreator := data.NewTestDataService()
	testDataCreator.CreateTestData(1)

	for _, v := range testDataCreator.GetTestData() {
		uploader.UploadTestData(v)
	}
	os.Exit(0)
}

func bootRPCService() {
	parsedEnv, _ := env.ParseEnvironmentVariables()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", parsedEnv.Port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterTransactionServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
