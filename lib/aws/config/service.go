// Package config provides a client for AWS Config.
//
package config

import (
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
)

// AWS Elastic Beanstalk makes it easy for you to create, deploy, and manage
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
type Config struct {
  Profile         string
  Region          string
  Session         *session.Session
  AwsConfig       *aws.Config
}

// New creates a new instance of the ElasticBeanstalk client.
//
// Example:
//     // Create a ElasticBeanstalk client.
//     svc := elasticbeanstalk.New(applicationName, environmentName,
//                                 versionLabel, description, s3Bucket, s3Key)
//
func New(region string, accessKeyId string, secretAccessKey string, cred_path string, profile string) *Config {
	return newClient(region, accessKeyId, secretAccessKey, cred_path, profile)
}

// newClient creates, initializes and returns a new AWS Config service.
//
func newClient(region string, accessKeyId string, secretAccessKey string, cred_path string, profile string) *Config {
  conf := &aws.Config{
    Credentials: GetCredentials(accessKeyId, secretAccessKey, cred_path, profile),
    Region: aws.String(region),
  }

  sess := session.Must(session.NewSession(conf))

  svc := &Config{
    Profile: profile,
    Region: region,
    Session: sess,
    AwsConfig: conf,
	}

  svc.PrintInfo()
	return svc
}
