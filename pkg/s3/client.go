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
		Region:           aws.String(os.Getenv("CDN_REGION")),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("CDN_ACCESS_KEY"), os.Getenv("CDN_SECRET_KEY"), ""),
		Endpoint:         aws.String(os.Getenv("CDN_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(true),
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
		Bucket: aws.String(os.Getenv("CDN_BUCKET")),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	cdnUrl := fmt.Sprintf("%s/%s", os.Getenv("CDN_URL"), fileName)

	return cdnUrl, nil
}

func DeleteFile(fileURL string) error {
	s, err := s3Client()
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %v", err)
	}

	_, err = url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("failed to parse file URL: %v", err)
	}

	cdnEndpoint := os.Getenv("CDN_ENDPOINT")
	if cdnEndpoint == "" {
		return fmt.Errorf("CDN_ENDPOINT environment variable is not set")
	}

	path := strings.TrimPrefix(fileURL, cdnEndpoint)

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid file URL format")
	}

	bucket := os.Getenv("CDN_BUCKET")
	key := strings.Join(parts, "/")

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
