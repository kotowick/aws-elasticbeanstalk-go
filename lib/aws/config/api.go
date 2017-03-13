// Package config provides a client for AWS Config.
//
package config

import (
   "github.com/aws/aws-sdk-go/aws/credentials"
   "../../utils"
   "os"
   "fmt"
)

// GetCredentials returns a new AWS credentials.Credentials
//
// credentials.Credentials can either be from:
//  - user-specified variables (command-line input)
//  - aws credentials file (typically under ~/.aws/credentials)
//  - environment variables
//
func GetCredentials(accessKeyId string, secretAccessKey string, cred_path string, profile string) *credentials.Credentials {
  if userCreds, err := GetUserSpecifiedCredentials(accessKeyId, secretAccessKey); err == nil {
    return userCreds
  } else if fileCreds, err := GetFileCredentials(cred_path, profile); err == nil {
    return fileCreds
  } else if envCreds, err := GetEnvCredentials(); err == nil {
    return envCreds
  } else {
    fmt.Println("Error: no AWS credentials found. Checked for ENV variables, credentials file, and command-line arguments.")
    os.Exit(1)
  }

  return &credentials.Credentials{}
}

// GetFileCredentials returns a new AWS credentials.Credentials
//
// credentials.Credentials aws_access_key_id and aws_secret_access_key are based
// off of the ~/.aws/credentials file
//
func GetFileCredentials(path string, profile string) (*credentials.Credentials, error){
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
// credentials.Credentials aws_access_key_id and aws_secret_access_key are based
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
// credentials.Credentials aws_access_key_id and aws_secret_access_key are based
// off a user input (command-line arguments)
//
func GetUserSpecifiedCredentials(aws_access_key_id string, aws_secret_access_key string) (*credentials.Credentials, error) {
  token := ""

  creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)

  if _, err := creds.Get(); err != nil {
		return &credentials.Credentials{}, err
	}
  return creds, nil
}

// GetUserSpecifiedCredentials prints out the credentials from a Config
//
func (c *Config) PrintInfo(){
  fmt.Println("#--")

  fmt.Println("Value of credentials : ", c.AwsConfig.Credentials)

  if c.Profile != "" {
    fmt.Println("Profile : ", c.Profile)
  }

  fmt.Println("--#")
}
