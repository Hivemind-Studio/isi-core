package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/url"
	"os"
	"strings"
)

func s3Client() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("DO_SPACE_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("DO_SPACE_ACCESS_KEY"), os.Getenv("DO_SPACE_SECRET_KEY"), ""),
		Endpoint:    aws.String(os.Getenv("DO_SPACE_ENDPOINT")),
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func UploadFile(filePath string, fileName string, username string) (string, error) {
	s, err := s3Client()
	if err != nil {
		return "", fmt.Errorf("failed to create S3 client: %v", err)
	}

	uploader := s3manager.NewUploader(s)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	fileName = strings.ToLower(username) + "/" + fileName

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("DO_SPACE_BUCKET")),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	cdnURL := fmt.Sprintf("https://%s.%s/%s", os.Getenv("DO_SPACE_BUCKET"), os.Getenv("DO_SPACE_ENDPOINT"), fileName)

	return cdnURL, nil
}

func DeleteFile(fileURL string) error {
	s, err := s3Client()
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %v", err)
	}

	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("failed to parse file URL: %v", err)
	}

	hostParts := strings.Split(parsedURL.Host, ".")
	if len(hostParts) < 2 {
		return fmt.Errorf("invalid file URL format")
	}

	bucket := hostParts[0]
	key := strings.TrimPrefix(parsedURL.Path, "/")

	svc := s3.New(s)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to confirm file deletion: %v", err)
	}

	return nil
}
