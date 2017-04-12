// Package cloudformation provides a client for AWS CloudFormation.
//
package cloudformation

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/segmentio/go-camelcase"
)

// DescribeStack describes the specified Cloudformation Stack.
//
// Example:
//     // Describe a Cloudformation stack.
//     svc := cloudformation.DescribeStack(stackName)
//
func (c *CloudFormation) DescribeStack(stackName string) bool {
	if stackName == "" {
		return false
	}

	params := &cloudformation.DescribeStacksInput{
		NextToken: aws.String("NextToken"),
		StackName: aws.String(stackName),
	}

	resp, err := c.Service.DescribeStacks(params)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println(resp)

	if len(resp.Stacks) > 0 {
		return true
	}
	return false
}

// UpdateStack updates an existing Cloudformation Stack
//
// Example:
//     // Update a Cloudformation Stack.
//     svc := cloudformation.UpdateStack(paramaters, template)
//
func (c *CloudFormation) UpdateStack(paramaters map[string]string, template []byte, usePreviousTemplate bool) {
	var cfParameters []*cloudformation.Parameter

	if usePreviousTemplate {
		delete(paramaters, "VPC_PRIVATE_SUBNETS")
		delete(paramaters, "SOLUTION_STACK")
		delete(paramaters, "VPC_ELB_SUBNETS")
		delete(paramaters, "VPC_ID")

		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("vpcPrivateSubnets"),
			UsePreviousValue: aws.Bool(true),
		})

		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("solutionStack"),
			UsePreviousValue: aws.Bool(true),
		})

		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("vpcElbSubnets"),
			UsePreviousValue: aws.Bool(true),
		})

		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("vpcId"),
			UsePreviousValue: aws.Bool(true),
		})
	}

	if paramaters["VersionLabel"] == "" {
		delete(paramaters, "VersionLabel")
		delete(paramaters, "AppBucket")
		delete(paramaters, "AppKey")
		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("versionLabel"),
			UsePreviousValue: aws.Bool(true),
		})
		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("appBucket"),
			UsePreviousValue: aws.Bool(true),
		})
		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("appKey"),
			UsePreviousValue: aws.Bool(true),
		})
	}
	if paramaters["AppBucket"] == "" {
		delete(paramaters, "AppBucket")
		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("appBucket"),
			UsePreviousValue: aws.Bool(true),
		})
	}
	if paramaters["AppKey"] == "" {
		delete(paramaters, "AppKey")
		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String("appKey"),
			UsePreviousValue: aws.Bool(true),
		})
	}

	for k, v := range paramaters {
		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:   aws.String(camelcase.Camelcase(k)),
			ParameterValue: aws.String(fmt.Sprintf("%s", v)),
		})
	}

	fmt.Printf("%s", cfParameters)

	params := &cloudformation.UpdateStackInput{
		StackName:  aws.String(fmt.Sprintf("%s-%s", paramaters["AppName"], paramaters["EnvName"])), // Required
		Parameters: cfParameters,
	}

	if usePreviousTemplate {
		fmt.Println("Using previous template")
		params.UsePreviousTemplate = aws.Bool(usePreviousTemplate)
	} else {
		fmt.Println("Using template file")
		params.TemplateBody = aws.String(fmt.Sprintf("%s", template))
	}

	resp, err := c.Service.UpdateStack(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}

// CreateStack creates a new Cloudformation Stack
//
// Example:
//     // Create a Cloudformation Stack.
//     svc := cloudformation.CreateStack(paramaters, template)
//
func (c *CloudFormation) CreateStack(paramaters map[string]string, template []byte) {
	var cfParameters []*cloudformation.Parameter

	for k, v := range paramaters {
		cfParameters = append(cfParameters, &cloudformation.Parameter{
			ParameterKey:   aws.String(camelcase.Camelcase(k)),
			ParameterValue: aws.String(fmt.Sprintf("%s", v)),
		})
	}

	fmt.Printf("%s", cfParameters)

	params := &cloudformation.CreateStackInput{
		StackName:       aws.String(fmt.Sprintf("%s-%s", paramaters["AppName"], paramaters["EnvName"])), // Required
		DisableRollback: aws.Bool(false),
		//OnFailure: aws.String("OnFailure"),
		Parameters:       cfParameters,
		TemplateBody:     aws.String(fmt.Sprintf("%s", template)),
		TimeoutInMinutes: aws.Int64(25),
	}

	resp, err := c.Service.CreateStack(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}

// DeleteStack deletes a specified Stack
//
// Example:
//     // Delete a Cloudformation stack.
//     svc := cloudformation.DeleteStack(stackName)
//
func (c *CloudFormation) DeleteStack(stackName string) {
	params := &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	}
	resp, err := c.Service.DeleteStack(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}
