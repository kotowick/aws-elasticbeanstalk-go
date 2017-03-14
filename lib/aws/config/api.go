// Package config provides a client for AWS Config.
//
package config

import (
	"fmt"

	"../../utils"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// GetCredentials returns a new AWS credentials.Credentials
//
// credentials.Credentials can either be from:
//  - user-specified variables (command-line input)
//  - aws credentials file (typically under ~/.aws/credentials)
//  - environment variables
//
func GetCredentials(accessKeyID string, secretAccessKey string, credPath string, profile string) *credentials.Credentials {
	if userCreds, err := GetUserSpecifiedCredentials(accessKeyID, secretAccessKey); err == nil {
		return userCreds
	} else if fileCreds, err := GetFileCredentials(credPath, profile); err == nil {
		return fileCreds
	} else if envCreds, err := GetEnvCredentials(); err == nil {
		return envCreds
	}

	fmt.Println("Error: no AWS credentials found. Checked for ENV variables, credentials file, and command-line arguments.")

	return &credentials.Credentials{}
}

// GetFileCredentials returns a new AWS credentials.Credentials
//
// credentials.Credentials awsAccessKeyID and awsSecretAccessKey are based
// off of the ~/.aws/credentials file
//
func GetFileCredentials(path string, profile string) (*credentials.Credentials, error) {
	credPath := utils.GetDefault(path, "/Users/spkotowick/.aws/credentials")
	profilePath := utils.GetDefault(profile, "default")

	// the file location and load default profile
	creds := credentials.NewSharedCredentials(credPath, profilePath)
	if _, err := creds.Get(); err != nil {
		return &credentials.Credentials{}, err
	}
	return creds, nil
}

// GetEnvCredentials returns a new AWS credentials.Credentials
//
// credentials.Credentials awsAccessKeyID and awsSecretAccessKey are based
// off of environment variables
//
func GetEnvCredentials() (*credentials.Credentials, error) {
	creds := credentials.NewEnvCredentials()
	if _, err := creds.Get(); err != nil {
		return &credentials.Credentials{}, err
	}
	return creds, nil
}

// GetUserSpecifiedCredentials returns a new AWS credentials.Credentials
//
// credentials.Credentials awsAccessKeyID and awsSecretAccessKey are based
// off a user input (command-line arguments)
//
func GetUserSpecifiedCredentials(awsAccessKeyID string, awsSecretAccessKey string) (*credentials.Credentials, error) {
	token := ""

	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, token)

	if _, err := creds.Get(); err != nil {
		return &credentials.Credentials{}, err
	}
	return creds, nil
}

// PrintInfo prints out the credentials from a Config
//
func (c *Config) PrintInfo() {
	fmt.Println("#--")

	fmt.Println("Value of credentials : ", c.AwsConfig.Credentials)

	if c.Profile != "" {
		fmt.Println("Profile : ", c.Profile)
	}

	fmt.Println("--#")
}
