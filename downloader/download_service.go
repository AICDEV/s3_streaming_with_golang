package downloader

import (
	cloud "github.com/aicdev/s3_streaming_with_golang/aws"
	"github.com/aicdev/s3_streaming_with_golang/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type DownloaderServiceInterface interface {
	DownloadFromS3(string) ([]byte, error)
}

type downloadService struct {
	env *env.Streamingnvironment
}

func NewDownloaderService() DownloaderServiceInterface {
	parsedEnv, _ := env.ParseEnvironmentVariables()
	return &downloadService{
		env: parsedEnv,
	}
}

func (dw *downloadService) DownloadFromS3(key string) ([]byte, error) {

	awsSession := cloud.AwsService.GetSession()
	dwm := s3manager.NewDownloader(awsSession, func(d *s3manager.Downloader) {
		d.Concurrency = 10
	})

	raw := &aws.WriteAtBuffer{}
	_, err := dwm.Download(raw, &s3.GetObjectInput{
		Bucket: aws.String(dw.env.S3BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return raw.Bytes(), nil
}
