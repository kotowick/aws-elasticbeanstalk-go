// Package elasticbeanstalk provides a client for AWS Elastic Beanstalk.
//
package elasticbeanstalk

import (
	"../config"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

// ElasticBeanstalk makes it easy for you to create, deploy, and manage
// scalable, fault-tolerant applications running on the Amazon Web Services
// cloud.
//
// For more information about this product, go to the AWS Elastic Beanstalk
// (http://aws.amazon.com/elasticbeanstalk/) details page. The location of the
// latest AWS Elastic Beanstalk WSDL is http://elasticbeanstalk.s3.amazonaws.com/doc/2010-12-01/AWSElasticBeanstalk.wsdl
// (http://elasticbeanstalk.s3.amazonaws.com/doc/2010-12-01/AWSElasticBeanstalk.wsdl).
// To install the Software Development Kits (SDKs), Integrated Development Environment
// (IDE) Toolkits, and command line tools that enable you to access the API,
// go to Tools for Amazon Web Services (http://aws.amazon.com/tools/).
//
type ElasticBeanstalk struct {
	ApplicationName    string
	EnvironmentName    string
	VersionLabel       string
	UniqueVersionLabel string
	Description        string
	S3Bucket           string
	S3Key              string
	Tier               string
	AwsConfig          *config.Config
	Service            *elasticbeanstalk.ElasticBeanstalk
}

// New creates a new instance of the ElasticBeanstalk client.
//
// Example:
//     // Create a ElasticBeanstalk client.
//     svc := elasticbeanstalk.New(applicationName, environmentName,
//                                 versionLabel, description, s3Bucket, s3Key)
//
func New(applicationName string, environmentName string, versionLabel string, description string, s3Bucket [2]string, s3Key string, tier string, awsConfig config.Config) *ElasticBeanstalk {
	return newClient(applicationName, environmentName, versionLabel, description, s3Bucket, s3Key, tier, awsConfig)
}

// newClient creates, initializes and returns a new elasticbeanstalk client instance.
//
func newClient(applicationName string, environmentName string, versionLabel string, description string, s3Bucket [2]string, s3Key string, tier string, awsConfig config.Config) *ElasticBeanstalk {
	svc := &ElasticBeanstalk{
		ApplicationName:    applicationName,
		EnvironmentName:    environmentName,
		UniqueVersionLabel: environmentName + "-" + tier + "-" + versionLabel,
		VersionLabel:       versionLabel,
		Description:        description,
		S3Bucket:           s3Bucket[0],
		S3Key:              s3Bucket[1] + "/" + versionLabel + ".zip",
		AwsConfig:          &awsConfig,
		Service:            elasticbeanstalk.New(awsConfig.Session, awsConfig.AwsConfig),
	}

	return svc
}
