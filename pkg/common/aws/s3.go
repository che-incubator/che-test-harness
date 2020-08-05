package aws

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// WriteToS3 writes the given byte array to S3.
func WriteToS3(outputKey string, data []byte) error {
	bucket, key, err := ParseS3URL(outputKey)

	if err != nil {
		return fmt.Errorf("error trying to parse S3 URL: %v", err)
	}

	session, err := AWSSession.getSession()

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(session)

	reader := bytes.NewReader(data)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	fmt.Println(err)

	return err
}

// CreateS3URL creates an S3 URL from a bucket and a key string.
func CreateS3URL(bucket string, keys ...string) string {
	strippedBucket := strings.Trim(bucket, "/")

	strippedKeys := make([]string, len(keys))
	for i, key := range keys {
		strippedKeys[i] = strings.Trim(key, "/")
	}

	s3JoinArray := []string{"s3:", strippedBucket}
	s3JoinArray = append(s3JoinArray, strippedKeys...)

	return strings.Join(s3JoinArray, "/")
}

// ParseS3URL parses an S3 url into a bucket and key.
func ParseS3URL(s3URL string) (string, string, error) {
	parsedURL, err := url.Parse(s3URL)

	if err != nil {
		return "", "", err
	}

	return parsedURL.Host, parsedURL.Path, nil
}
