package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func S3Client() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("DO_SPACE_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("DO_SPACE_ACCESS_KEY"), os.Getenv("DO_SPACE_SECRET_KEY"), ""),
		Endpoint:    aws.String(os.Getenv("DO_SPACE_ENDPOINT")),
	})
	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}
