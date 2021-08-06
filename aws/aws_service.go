package cloud

import (
	"log"

	"github.com/aicdev/s3_streaming_with_golang/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsServiceInterface interface {
	GetSession() *session.Session
}

type awsService struct {
	Session *session.Session
}

var (
	AwsService AwsServiceInterface = &awsService{}
)

func (as *awsService) GetSession() *session.Session {
	if as.Session != nil {
		return as.Session
	}

	parsedEnv, _ := env.ParseEnvironmentVariables()

	sess, err := session.NewSession(&aws.Config{
		Region:      &parsedEnv.S3Region,
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		log.Fatal(err)
	}

	as.Session = sess

	return as.Session
}
