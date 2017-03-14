// Package cloudformation provides a client for AWS CloudFormation.
//
package cloudformation

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/segmentio/go-camelcase"
)

func (c *CloudFormation) CreateStack(paramaters map[string]string, template []byte) {
	var cf_parameters []*cloudformation.Parameter

	for k, v := range paramaters {
		cf_parameters = append(cf_parameters, &cloudformation.Parameter{
			ParameterKey:   aws.String(camelcase.Camelcase(k)),
			ParameterValue: aws.String(fmt.Sprintf("%s", v)),
		})
	}

	fmt.Printf("%s", cf_parameters)

	params := &cloudformation.CreateStackInput{
		StackName:       aws.String(fmt.Sprintf("%s-%s-%s", paramaters["AppName"], paramaters["EnvName"], paramaters["Tier"])), // Required
		DisableRollback: aws.Bool(false),
		//OnFailure: aws.String("OnFailure"),
		Parameters:       cf_parameters,
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

func (c *CloudFormation) DeleteStack(stack_name string) {
	params := &cloudformation.DeleteStackInput{
		StackName: aws.String(stack_name),
	}
	resp, err := c.Service.DeleteStack(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}
