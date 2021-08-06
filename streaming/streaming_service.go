package services

import (
	"log"

	cloud "github.com/aicdev/s3_streaming_with_golang/aws"
	"github.com/aicdev/s3_streaming_with_golang/downloader"
	"github.com/aicdev/s3_streaming_with_golang/env"
	pb "github.com/aicdev/s3_streaming_with_golang/proto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"google.golang.org/protobuf/proto"
)

type StreamingServiceInterface interface {
	Stream(string, chan *pb.TransactionCollection)
}

type streamingService struct{}

var (
	StreamingService StreamingServiceInterface = &streamingService{}
)

func (ss *streamingService) Stream(id string, c chan *pb.TransactionCollection) {
	parsedEnv, _ := env.ParseEnvironmentVariables()
	sc := cloud.AwsService.GetSession()

	bucket_client := s3.New(sc)

	var continueToken *string

	for {
		resp, _ := bucket_client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(parsedEnv.S3BucketName),
			Prefix:            aws.String(id),
			ContinuationToken: continueToken,
		})

		downloader := downloader.NewDownloaderService()
		for _, key := range resp.Contents {

			rawBytes, err := downloader.DownloadFromS3(*key.Key)

			if err != nil {
				log.Fatal(err)
			}

			transactionCollection := &pb.TransactionCollection{}
			err = proto.Unmarshal(rawBytes, transactionCollection)

			if err != nil {
				log.Fatal(err)
			}

			if len(transactionCollection.GetTransactions()) > 0 {
				c <- transactionCollection
			}
		}

		if !aws.BoolValue(resp.IsTruncated) {
			break
		}

		continueToken = resp.NextContinuationToken
	}

	close(c)
}
