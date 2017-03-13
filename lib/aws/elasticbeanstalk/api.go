// Package elasticbeanstalk provides a client for AWS Elastic Beanstalk.
//
package elasticbeanstalk

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
  "fmt"
)

// UpdateEnvironment API operation for Elastic Beanstalk.
//
// Updates the environment description, deploys a new application version, updates
// the configuration settings to an entirely new configuration template, or
// updates select configuration option values in the running environment.
//
// Attempting to update both the release and configuration is not allowed and
// AWS Elastic Beanstalk returns an InvalidParameterCombination error.
//
// When updating the configuration settings to a new template or individual
// settings, a draft configuration is created and DescribeConfigurationSettings
// for this environment returns two setting descriptions with different DeploymentStatus
// values.
//
// See the AWS API reference guide for AWS Elastic Beanstalk's
// API operation UpdateEnvironment for usage and error information.
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticbeanstalk-2010-12-01/UpdateEnvironment
//
func (c *ElasticBeanstalk) UpdateEnvironment(){
  c.CreateApplicationVersion()
  c.UpdateApplicationVersion()

  params := &elasticbeanstalk.UpdateEnvironmentInput{
      ApplicationName: aws.String(c.ApplicationName),
      EnvironmentName: aws.String(c.EnvironmentName),
      VersionLabel: aws.String(c.UniqueVersionLabel),
  }
  resp, err := c.Service.UpdateEnvironment(params)

  if err != nil {
      fmt.Println(err.Error())
      return
  }

  // Pretty-print the response data.
  fmt.Println(resp)
}

// UpdateApplicationVersion API operation for AWS Elastic Beanstalk.
//
// Updates the specified application version to have the specified properties.
//
// If a property (for example, description) is not provided, the value remains
// unchanged. To clear properties, specify an empty string.
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for AWS Elastic Beanstalk's
// API operation UpdateApplicationVersion for usage and error information.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticbeanstalk-2010-12-01/UpdateApplicationVersion
//
func (c *ElasticBeanstalk) UpdateApplicationVersion(){
  params := &elasticbeanstalk.UpdateApplicationVersionInput{
      ApplicationName: aws.String(c.ApplicationName), // Required
      VersionLabel:    aws.String(c.UniqueVersionLabel),    // Required
      Description:     aws.String(c.Description),
  }

  resp, err := c.Service.UpdateApplicationVersion(params)

  if err != nil {
    fmt.Println(err.Error())
    return
  }

  // Pretty-print the response data.
  fmt.Println(resp)
}

// CreateApplicationVersion API operation for AWS Elastic Beanstalk.
//
// Creates an application version for the specified application. You can create
// an application version from a source bundle in Amazon S3, a commit in AWS
// CodeCommit, or the output of an AWS CodeBuild build as follows:
//
// Specify a commit in an AWS CodeCommit repository with SourceBuildInformation.
//
// Specify a build in an AWS CodeBuild with SourceBuildInformation and BuildConfiguration.
//
// Specify a source bundle in S3 with SourceBundle
//
// Omit both SourceBuildInformation and SourceBundle to use the default sample
// application.
//
// Once you create an application version with a specified Amazon S3 bucket
// and key location, you cannot change that Amazon S3 location. If you change
// the Amazon S3 location, you receive an exception when you attempt to launch
// an environment from the application version.
//
// See the AWS API reference guide for AWS Elastic Beanstalk's
// API operation CreateApplicationVersion for usage and error information.
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticbeanstalk-2010-12-01/CreateApplicationVersion
//
func (c *ElasticBeanstalk) CreateApplicationVersion(){
  params := &elasticbeanstalk.CreateApplicationVersionInput{
      ApplicationName:        aws.String(c.ApplicationName), // Required
      VersionLabel:           aws.String(c.UniqueVersionLabel),    // Required
      Description:            aws.String(c.Description),
      AutoCreateApplication:  aws.Bool(true),
      Process:                aws.Bool(true),
      SourceBundle: &elasticbeanstalk.S3Location{
        S3Bucket:             aws.String(c.S3Bucket),
        S3Key:                aws.String(c.S3Key),
      },
  }

  resp, err := c.Service.CreateApplicationVersion(params)

  if err != nil {
    fmt.Println(err.Error())
    return
  }

  fmt.Println(resp)
}

// ListApplications API operation for AWS Elastic Beanstalk.
//
// Returns the descriptions of existing applications.
//
// See the AWS API reference guide for AWS Elastic Beanstalk's
// API operation DescribeApplications for usage and error information.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticbeanstalk-2010-12-01/DescribeApplications
func (c *ElasticBeanstalk) ListApplications(verbose bool){
  resp, err := c.Service.DescribeApplications(&elasticbeanstalk.DescribeApplicationsInput{})

  if err != nil {
    fmt.Println(err.Error())
    return
	}

  if verbose {
    fmt.Println(resp)
  } else {
    for j := range resp.Applications{
      fmt.Printf("%s", *resp.Applications[j].ApplicationName)
    }
  }
}

// ListEnvironments API operation for AWS Elastic Beanstalk.
//
// Returns descriptions for existing environments.
//
// See the AWS API reference guide for AWS Elastic Beanstalk's
// API operation DescribeEnvironments for usage and error information.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticbeanstalk-2010-12-01/DescribeEnvironments
func (c *ElasticBeanstalk) ListEnvironments(verbose bool, applicationName string, environmentName string){
  params := &elasticbeanstalk.DescribeEnvironmentsInput{}

  if applicationName != "" {
    params.ApplicationName = aws.String(applicationName)
  }

  if environmentName != "" {
    params.EnvironmentNames = []*string{aws.String(environmentName)}
  }

  resp, err := c.Service.DescribeEnvironments(params)

  if err != nil {
      fmt.Println(err.Error())
      return
  }

  if verbose {
    fmt.Println(resp)
  } else {
    for j := range resp.Environments{
      fmt.Printf("%s", *resp.Environments[j].EnvironmentName)
    }
  }
}

// DeleteApplication API operation for AWS Elastic Beanstalk.
//
// Deletes the specified application along with all associated versions and
// configurations. The application versions will not be deleted from your Amazon
// S3 bucket.
//
// You cannot delete an application that has a running environment.
//
// See the AWS API reference guide for AWS Elastic Beanstalk's
// API operation DeleteApplication for usage and error information.
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticbeanstalk-2010-12-01/DeleteApplication
func (c *ElasticBeanstalk) DeleteApplication(applicationName string){
  params := &elasticbeanstalk.DeleteApplicationInput{
      ApplicationName:     aws.String(applicationName), // Required
      //TerminateEnvByForce: aws.Bool(true),
  }
  resp, err := c.Service.DeleteApplication(params)

  if err != nil {
      fmt.Println(err.Error())
      return
  }

  fmt.Println(resp)
}

// CreateApplication API operation for AWS Elastic Beanstalk.
//
// Creates an application that has one configuration template named default
// and no application versions.
//
// See the AWS API reference guide for AWS Elastic Beanstalk's
// API operation CreateApplication for usage and error information.
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticbeanstalk-2010-12-01/CreateApplication
func (c *ElasticBeanstalk) CreateApplication(applicationName string) {
	params := &elasticbeanstalk.CreateApplicationInput{ ApplicationName: aws.String(applicationName) }
  resp, err := c.Service.CreateApplication(params)

  if err != nil {
      fmt.Println(err.Error())
      return
  }

  fmt.Println(resp)
}
