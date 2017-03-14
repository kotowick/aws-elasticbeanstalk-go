// Package cloudformation provides a client for AWS CloudFormation.
//
package cloudformation

import (
	"../config"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type CloudFormation struct {
	AwsConfig *config.Config
	Service   *cloudformation.CloudFormation
}

// New creates a new instance of the Cloudformation client.
//
// Example:
//     // Create a Cloudformation client.
//     svc := cloudformation.New(awsConfig config.Config)
//
func New(awsConfig config.Config) *CloudFormation {
	return newClient(awsConfig)
}

// newClient creates, initializes and returns a new Cloudformation client instance.
//
func newClient(awsConfig config.Config) *CloudFormation {
	svc := &CloudFormation{
		AwsConfig: &awsConfig,
		Service:   cloudformation.New(awsConfig.Session, awsConfig.AwsConfig),
	}

	return svc
}
