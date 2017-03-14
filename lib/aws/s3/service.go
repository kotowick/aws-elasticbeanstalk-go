// Package s3 provides a client for AWS S3.
//
package s3

import (
	"../config"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 is a client for Amazon S3.
// The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01
type S3 struct {
	AwsConfig *config.Config
	Service   *s3.S3
}

// New creates a new instance of the S3 client.
//
// Example:
//     // Create a S3 client.
//     svc := s3.New(awsConfig config.Config)
//
func New(awsConfig config.Config) *S3 {
	return newClient(awsConfig)
}

// newClient creates, initializes and returns a new S3 client instance.
//
func newClient(awsConfig config.Config) *S3 {
	svc := &S3{
		AwsConfig: &awsConfig,
		Service:   s3.New(awsConfig.Session, awsConfig.AwsConfig),
	}

	return svc
}
