// Package main provides a client for Elastic Beanstalk deployments.
package main

import (
	"fmt"
	"log"
	"regexp"

	"./lib/aws/cloudformation"
	"./lib/aws/config"
	"./lib/aws/elasticbeanstalk"
	"./lib/aws/s3"
	"./lib/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	//Connection Commands
	lists                              = kingpin.Command("list", "Commands related to VPN connection profiles")
	_                                  = lists.Command("applications", "List applications")
	listEnvironments                   = lists.Command("environments", "List environments")
	listApplicationEnvironmentsFlag    = listEnvironments.Flag("app-name", "application_name.").Short('a').String()
	listApplicationEnvironmentNameFlag = listEnvironments.Flag("env-name", "Environment Name").Short('e').String()

	create = kingpin.Command("create", "Create EB")
	// Application
	createApplication        = create.Command("application", "Create Application")
	createApplicationNameArg = createApplication.Flag("app-name", "Application Name to create").Required().Short('a').String()
	// Environment
	createEnvironment                   = create.Command("environment", "Create Environment")
	createEnvironmentApplicationNameArg = createEnvironment.Flag("app-name", "Application Name").Required().Short('a').String()
	createEnvironmentNameArg            = createEnvironment.Flag("env-name", "Environment Name to create").Short('e').String()
	createEnvironmentTierArg            = createEnvironment.Flag("tier", "Environment Tier").Required().Short('t').String()
	createEnvironmentS3Arg              = createEnvironment.Flag("bucket", "S3 bucket name").Short('b').String()
	createEnvironmentConfigArg          = createEnvironment.Flag("config", "Config yml file").Short('c').String()
	createEnvironmentVersionArg         = createEnvironment.Flag("version-label", "Version").Short('l').String()
	createEnvironmentDescriptionArg     = createEnvironment.Flag("desc", "Description of Version").Short('d').String()

	delete = kingpin.Command("delete", "Destroy EB")
	// Application
	deleteApplication        = delete.Command("application", "Destroy application")
	deleteApplicationNameArg = deleteApplication.Flag("app-name", "Application to delete").Required().Short('a').String()

	// Environment
	deleteEnvironment                   = delete.Command("environment", "Destroy environment")
	deleteEnvironmentApplicationNameArg = deleteEnvironment.Flag("app-name", "Environment's application name").Required().Short('a').String()
	deleteEnvironmentNameArg            = deleteEnvironment.Flag("env-name", "Environment to delete").Required().Short('e').String()
	deleteEnvironmentTierArg            = deleteEnvironment.Flag("tier", "Tier Name (worker | webserver)").Required().Short('t').String()

	update                              = kingpin.Command("update", "Deploy EB")
	updateEnvironment                   = update.Command("environment", "Envrionment")
	updateEnvironmentApplicationNameArg = updateEnvironment.Flag("app-name", "The Application in which the environment lives under").Required().Short('a').String()
	updateEnvironmentNameArg            = updateEnvironment.Flag("env-name", "Envionment name").Required().Short('e').String()
	updateEnvironmentVersionArg         = updateEnvironment.Flag("version-label", "Version").Required().Short('l').String()
	updateEnvironmentTierArg            = updateEnvironment.Flag("tier", "Environment Tier").Required().Short('t').String()
	updateEnvironmentS3Arg              = updateEnvironment.Flag("bucket", "S3 Bucket Name + path").Short('b').String()
	updateEnvironmentDescriptionArg     = updateEnvironment.Flag("desc", "Description of Version").Short('d').String()

	rollback = kingpin.Command("rollback", "Rollback EB")

	// General Flags
	awsRegion          = kingpin.Flag("region", "AWS Region").Short('r').String()
	awsAccessKeyID     = kingpin.Flag("access-key-id", "AWS Access Key ID").Short('k').String()
	awsSecretAccessKey = kingpin.Flag("secret-access-key", "AWS Secret Access Key").Short('s').String()
	awsProfile         = kingpin.Flag("profile", "Aws Profile").Short('p').String()
	awsCredPath        = kingpin.Flag("cred-path", "Aws Cred Path").String()

	verbose = kingpin.Flag("verbose", "Verbose mode").Short('v').Bool()

	//Command Regex Section
	createRegex   = regexp.MustCompile(`^create`)
	listRegex     = regexp.MustCompile(`^list`)
	updateRegex   = regexp.MustCompile(`^update`)
	deleteRegex   = regexp.MustCompile(`^delete`)
	rollbackRegex = regexp.MustCompile(`^rollback`)

	//Global Vars
	cliVersion = "0.0.1"
)

// listFunction operation for Elastic Beanstalk Environment
//
// List all Elastic Beanstalk Applications AWS Elastic Beanstalk API
//
// List all Elastic Beanstalk Environments for the specified Elastic Beanstalk
// Application via AWS Elastic Beanstalk API
//
// Arguments:
//				- listMethod string
func listFunction(listMethod string, awsConfig config.Config) {
	ebService := elasticbeanstalk.New("", "", "", "", [2]string{}, "", "", awsConfig)

	switch listMethod {
	case "list applications":
		ebService.ListApplications(*verbose)
	case "list environments":
		ebService.ListEnvironments(*verbose, *listApplicationEnvironmentsFlag, *listApplicationEnvironmentNameFlag)
	default:
		log.Fatalf("not sure what to do with command: %s", listMethod)
	}
}

// createFunction operation for Elastic Beanstalk Environment
//
// Create the provided AWS Elastic Beanstalk Application via AWS Elastic
// Beanstalk API
//
// Create the provided AWS Elastic Beanstalk Environment via AWS Cloudformation
//
// Arguments:
//				- createMethod string
func createFunction(createMethod string, awsConfig config.Config) {
	s3Service := s3.New(awsConfig)
	bucketInfo := s3Service.ParseS3Bucket(*createEnvironmentS3Arg)

	ebService := elasticbeanstalk.New(
		*createEnvironmentApplicationNameArg, *createEnvironmentNameArg,
		*createEnvironmentVersionArg, *createEnvironmentDescriptionArg,
		bucketInfo, "versions/"+*createEnvironmentVersionArg+".zip",
		*createEnvironmentTierArg, awsConfig,
	)

	switch createMethod {
	case "create application":
		//aws.ShellOut(fmt.Sprintf("Creating application : %s", *createApplicationNameArg), aws.ShellOutParams{CmdName: "aws", CmdArgs:  []string{"elasticbeanstalk", "create-application", "--application-name", *createApplicationNameArg}})
		ebService.CreateApplication(*createApplicationNameArg)
	case "create environment":
		environment := utils.GetDefault(*createEnvironmentNameArg, *createEnvironmentApplicationNameArg)
		asset, err := Asset(fmt.Sprintf("resources/cloudformation/templates/%s_v1.json", *createEnvironmentTierArg))

		if err != nil {
			log.Fatalf("Asset not found: %s", err)
			return
		}

		s3Service.UploadSingleFile(*createEnvironmentS3Arg, "versions/"+*createEnvironmentVersionArg+".zip", *createEnvironmentVersionArg)

		additionalConfigOptions := make(map[string]string)
		additionalConfigOptions["EnvName"] = environment
		additionalConfigOptions["AppName"] = *createEnvironmentApplicationNameArg
		additionalConfigOptions["AppKey"] = s3Service.ParseS3Bucket(*createEnvironmentS3Arg)[1] + "/" + "versions/" + *createEnvironmentVersionArg + ".zip"
		additionalConfigOptions["Tier"] = *createEnvironmentTierArg

		configOptions := utils.GetConfig(*createEnvironmentConfigArg)
		configOptions = utils.CombineConfigOptions(configOptions, additionalConfigOptions)

		cfServcie := cloudformation.New(awsConfig)

		cfServcie.CreateStack(configOptions, asset)
	}
}

// deleteFunction operation for Elastic Beanstalk Environment
//
// Delete the provided Elastic Beanstalk Environment through its related
// cloudformation stack name.
//
// Arguments:
//				- deleteMethod string
func deleteFunction(deleteMethod string, awsConfig config.Config) {
	switch deleteMethod {
	case "delete application":
		ebService := elasticbeanstalk.New("", "", "", "", [2]string{}, "", "", awsConfig)
		ebService.DeleteApplication(*deleteApplicationNameArg)
	case "delete environment":
		cfServcie := cloudformation.New(awsConfig)

		cfServcie.DeleteStack(fmt.Sprintf("%s-%s-%s", *deleteEnvironmentApplicationNameArg, *deleteEnvironmentNameArg, *deleteEnvironmentTierArg))
	}
}

// updateFunction operation for Elastic Beanstalk Environment
//
// Update the provided Elastic Beanstalk Environment version through the AWS
// Elastic Beanstalk API.
//
// TODO:
//		Update the provided Elastic Beanstalk Environment settings through its
// 		related Cloudformation stack name.
//
// Arguments:
//				- updateMethod string
func updateFunction(updateMethod string, awsConfig config.Config) {
	s3Service := s3.New(awsConfig)

	bucketInfo := s3Service.ParseS3Bucket(*updateEnvironmentS3Arg)

	svc := elasticbeanstalk.New(
		*updateEnvironmentApplicationNameArg, *updateEnvironmentNameArg,
		*updateEnvironmentVersionArg, *updateEnvironmentDescriptionArg,
		bucketInfo, "versions/"+*updateEnvironmentVersionArg+".zip",
		*updateEnvironmentTierArg, awsConfig,
	)

	switch updateMethod {
	case "update environment":
		if *updateEnvironmentS3Arg != "" {
			s3Service.UploadSingleFile(*updateEnvironmentS3Arg, "versions/"+*updateEnvironmentVersionArg+".zip", *updateEnvironmentVersionArg)
		}
		svc.UpdateEnvironment()
	}
}

//	Entry Controller
func main() {
	kingpin.Version(cliVersion)

	parsedArg := kingpin.Parse()

	awsConfig := config.New(
		*awsRegion,
		*awsAccessKeyID,
		*awsSecretAccessKey,
		*awsCredPath,
		*awsProfile,
	)

	switch {
	case createRegex.MatchString(parsedArg):
		createFunction(parsedArg, *awsConfig)
	case listRegex.MatchString(parsedArg):
		listFunction(parsedArg, *awsConfig)
	case updateRegex.MatchString(parsedArg):
		updateFunction(parsedArg, *awsConfig)
	case deleteRegex.MatchString(parsedArg):
		deleteFunction(parsedArg, *awsConfig)
	default:
		log.Fatalf("Command signature not recognized: %s", parsedArg)
	}
}
