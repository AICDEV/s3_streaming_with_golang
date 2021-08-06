package env

import (
	"github.com/kelseyhightower/envconfig"
)

type Streamingnvironment struct {
	Port         string `required:"true" default:"9999"`
	S3Region     string `required:"true" default:"eu-central-1"`
	S3BucketName string `required:"true" default:"user-transactions-dev"`
}

func ParseEnvironmentVariables() (*Streamingnvironment, error) {
	env := &Streamingnvironment{}
	if err := envconfig.Process("streaming_example", env); err != nil {
		return env, err
	}

	return env, nil
}
