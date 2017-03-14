// Package s3 provides a client for AWS S3.
//
package s3

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadSingleFile uploads a file to S3.
//
// Example:
//     // Upload fileName to S3.
//     svc := s3.UploadSingleFile(bucketPath, fileName, versionLabel)
//
func (c *S3) UploadSingleFile(bucketPath string, fileName string, versionLabel string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("err opening file: %s", err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	bucketInfo := c.ParseS3Bucket(bucketPath)

	fmt.Printf("Uploading [ %s ] to S3 bucket : %s", file.Name(), bucketPath)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketInfo[0]),
		Key:           aws.String(bucketInfo[1] + "/" + versionLabel + ".zip"),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	resp, err := c.Service.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}
	fmt.Printf("response %s", awsutil.StringValue(resp))
}

// ParseS3Bucket splits up a s3 path (bucket/directory) string.
//
// Example:
//     // Split s3 path
//     svc := s3.ParseS3Bucket(bucketPath)
//
func (c *S3) ParseS3Bucket(bucketPath string) [2]string {
	bucketSplit := strings.Split(bucketPath, "/")
	var rtn [2]string

	rtn[0] = bucketSplit[0]

	rtn[1] = strings.Join(bucketSplit[1:], "/")

	return rtn
}
