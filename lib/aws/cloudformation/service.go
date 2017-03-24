// Package cloudformation provides a client for AWS CloudFormation.
//
package cloudformation

import (
	"../config"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// CloudFormation allows you to create and manage AWS infrastructure deployments
// predictably and repeatedly. You can use AWS CloudFormation to leverage AWS
// products, such as Amazon Elastic Compute Cloud, Amazon Elastic Block Store,
// Amazon Simple Notification Service, Elastic Load Balancing, and Auto Scaling
// to build highly-reliable, highly scalable, cost-effective applications without
// creating or configuring the underlying AWS infrastructure.
//
// With AWS CloudFormation, you declare all of your resources and dependencies
// in a template file. The template defines a collection of resources as a single
// unit called a stack. AWS CloudFormation creates and deletes all member resources
// of the stack together and manages all dependencies between the resources
// for you.
//
// For more information about AWS CloudFormation, see the AWS CloudFormation
// Product Page (http://aws.amazon.com/cloudformation/).
//
// Amazon CloudFormation makes use of other AWS products. For additional technical
// information about a specific AWS product, see its technical documentation
// (http://docs.aws.amazon.com/).
// The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cloudformation-2010-05-15
//
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
