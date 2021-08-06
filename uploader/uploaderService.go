package uploader

import (
	"bytes"
	"fmt"

	"log"

	cloud "github.com/aicdev/s3_streaming_with_golang/aws"
	"github.com/aicdev/s3_streaming_with_golang/env"
	pb "github.com/aicdev/s3_streaming_with_golang/proto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"google.golang.org/protobuf/proto"
)

type UploaderServiceInterface interface {
	UploadTestData(*pb.User)
}

type uploaderService struct {
	collectionChunkSize int
}

func NewUploaderService() UploaderServiceInterface {
	return &uploaderService{
		collectionChunkSize: 500,
	}
}

func (ups *uploaderService) UploadTestData(user *pb.User) {
	for i := 0; i < cap(user.GetTransactions()); i += ups.collectionChunkSize {
		if (i + int(ups.collectionChunkSize)) > cap(user.GetTransactions())-1 {
			fragments := user.GetTransactions()[i : cap(user.GetTransactions())-1]

			tc := &pb.TransactionCollection{
				Transactions: fragments,
			}

			raw, _ := proto.Marshal(tc)
			ups.upload(user.GetId(), fmt.Sprintf("transaction-%d", i), raw)
		} else {
			fragments := user.GetTransactions()[i:cap(user.GetTransactions())]
			tc := &pb.TransactionCollection{
				Transactions: fragments,
			}

			raw, _ := proto.Marshal(tc)
			ups.upload(user.GetId(), fmt.Sprintf("transaction-%d", i), raw)
		}
	}
}

func (ups *uploaderService) upload(userId string, key string, raw []byte) {
	parsedEnv, _ := env.ParseEnvironmentVariables()

	session := cloud.AwsService.GetSession()
	up := s3manager.NewUploader(session)

	_, err := up.Upload(&s3manager.UploadInput{
		Bucket: aws.String(parsedEnv.S3BucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", userId, key)),
		Body:   bytes.NewReader(raw),
	})

	if err != nil {
		log.Println(err.Error())
	}
}
