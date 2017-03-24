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
	createApplicationNameArg = createApplication.Flag("app-name", "Application Name to create").Required().Short('a').Envar("APP_NAME").String()
	// Environment
	createEnvironment                   = create.Command("environment", "Create Environment")
	createEnvironmentApplicationNameArg = createEnvironment.Flag("app-name", "Application Name").Envar("APP_NAME").Required().Short('a').String()
	createEnvironmentNameArg            = createEnvironment.Flag("env-name", "Environment Name to create").Short('e').Envar("ENV_NAME").String()
	createEnvironmentTierArg            = createEnvironment.Flag("tier", "Environment Tier").Required().Short('t').Envar("ENV_TIER").String()
	createEnvironmentS3Arg              = createEnvironment.Flag("bucket", "S3 bucket name").Short('b').Envar("ENV_S3_BUCKET_LOCATION").String()
	createEnvironmentConfigArg          = createEnvironment.Flag("config", "Config yml file").Short('c').Envar("ENV_CONFIG_PATH").String()
	createEnvironmentVersionArg         = createEnvironment.Flag("version-label", "Version").Short('l').Envar("ENV_VERSION_LABEL").String()
	createEnvironmentLocalFilePathArg   = createEnvironment.Flag("local-file-prefix", "The directory (prefix) of the file to upload.").Short('d').Envar("ENV_LOCAL_FILE_PREFIX").String()
	createEnvironmentDescriptionArg     = createEnvironment.Flag("desc", "Description of Version").Envar("ENV_DESCRIPTION").String()

	delete = kingpin.Command("delete", "Destroy EB")
	// Application
	deleteApplication        = delete.Command("application", "Destroy application")
	deleteApplicationNameArg = deleteApplication.Flag("app-name", "Application to delete").Required().Short('a').String()

	// Environment
	deleteEnvironment                   = delete.Command("environment", "Destroy environment")
	deleteEnvironmentApplicationNameArg = deleteEnvironment.Flag("app-name", "Environment's application name").Required().Short('a').Envar("APP_NAME").String()
	deleteEnvironmentNameArg            = deleteEnvironment.Flag("env-name", "Environment to delete").Required().Short('e').Envar("ENV_NAME").String()
	deleteEnvironmentTierArg            = deleteEnvironment.Flag("tier", "Tier Name (worker | webserver)").Required().Short('t').Envar("ENV_TIER").String()

	update                              = kingpin.Command("update", "Deploy EB")
	updateEnvironment                   = update.Command("environment", "Envrionment")
	updateEnvironmentApplicationNameArg = updateEnvironment.Flag("app-name", "The Application in which the environment lives under").Required().Short('a').Envar("APP_NAME").String()
	updateEnvironmentNameArg            = updateEnvironment.Flag("env-name", "Envionment name").Required().Short('e').Envar("ENV_NAME").String()
	updateEnvironmentVersionArg         = updateEnvironment.Flag("version-label", "Version").Short('l').Envar("ENV_VERSION_LABEL").String()
	updateEnvironmentTierArg            = updateEnvironment.Flag("tier", "Environment Tier").Required().Short('t').Envar("ENV_TIER").String()
	updateEnvironmentS3Arg              = updateEnvironment.Flag("bucket", "S3 Bucket Name + path").Short('b').Envar("ENV_S3_BUCKET_LOCATION").String()
	updateEnvironmentLocalFilePathArg   = updateEnvironment.Flag("local-file-prefix", "The directory (prefix) of the file to upload.").Short('d').Envar("ENV_LOCAL_FILE_PREFIX").String()
	updateEnvironmentDescriptionArg     = updateEnvironment.Flag("desc", "Description of Version").Envar("ENV_DESCRIPTION").String()
	updateEnvironmentConfigArg          = updateEnvironment.Flag("config", "Config yml file").Short('c').Envar("ENV_CONFIG_PATH").String()

	rollback = kingpin.Command("rollback", "Rollback EB")

	// General Flags
	awsRegion          = kingpin.Flag("region", "AWS Region").Short('r').Envar("AWS_REGION").String()
	awsAccessKeyID     = kingpin.Flag("access-key-id", "AWS Access Key ID").Short('k').Envar("AWS_ACCESS_KEY_ID").String()
	awsSecretAccessKey = kingpin.Flag("secret-access-key", "AWS Secret Access Key").Short('s').Envar("AWS_SECRET_ACCESS_KEY").String()
	awsProfile         = kingpin.Flag("profile", "Aws Profile").Short('p').Envar("AWS_PROFILE").String()
	awsCredPath        = kingpin.Flag("cred-path", "Aws Cred Path").Envar("AWS_CRED_PATH").String()

	verbose = kingpin.Flag("verbose", "Verbose mode").Short('v').Envar("VERBOSE_MODE").Bool()

	//Command Regex Section
	createRegex   = regexp.MustCompile(`^create`)
	listRegex     = regexp.MustCompile(`^list`)
	updateRegex   = regexp.MustCompile(`^update`)
	deleteRegex   = regexp.MustCompile(`^delete`)
	rollbackRegex = regexp.MustCompile(`^rollback`)

	//Global Vars
	cliVersion = "0.0.2"
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

	switch createMethod {
	case "create application":
		ebService := elasticbeanstalk.New(
			*createEnvironmentApplicationNameArg, *createEnvironmentNameArg,
			*createEnvironmentVersionArg, *createEnvironmentDescriptionArg,
			bucketInfo, *createEnvironmentVersionArg+".zip",
			*createEnvironmentTierArg, awsConfig,
		)
		ebService.CreateApplication(*createApplicationNameArg)
	case "create environment":
		environment := utils.GetDefault(*createEnvironmentNameArg, *createEnvironmentApplicationNameArg)
		asset, err := Asset(fmt.Sprintf("resources/cloudformation/templates/%s_v1.json", *createEnvironmentTierArg))

		if err != nil {
			log.Fatalf("Asset not found: %s", err)
			return
		}

		s3Service.UploadSingleFile(*createEnvironmentS3Arg, *createEnvironmentLocalFilePathArg+"/"+*createEnvironmentVersionArg+".zip", *createEnvironmentVersionArg)

		additionalConfigOptions := make(map[string]string)
		additionalConfigOptions["EnvName"] = environment
		additionalConfigOptions["AppName"] = *createEnvironmentApplicationNameArg
		additionalConfigOptions["AppBucket"] = bucketInfo[0]
		additionalConfigOptions["AppKey"] = bucketInfo[1] + "/" + *createEnvironmentVersionArg + ".zip"
		additionalConfigOptions["VersionLabel"] = *createEnvironmentVersionArg
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

	switch updateMethod {
	case "update environment":
		environment := utils.GetDefault(*updateEnvironmentNameArg, *updateEnvironmentApplicationNameArg)
		asset, err := Asset(fmt.Sprintf("resources/cloudformation/templates/%s_v1.json", *updateEnvironmentTierArg))

		if err != nil {
			log.Fatalf("Asset not found: %s", err)
			return
		}

		if *updateEnvironmentVersionArg != "" {
			s3Service.UploadSingleFile(*updateEnvironmentS3Arg, *updateEnvironmentLocalFilePathArg+"/"+*updateEnvironmentVersionArg+".zip", *updateEnvironmentVersionArg)
		}

		configOptions := make(map[string]string)
		configOptions["EnvName"] = environment
		configOptions["AppName"] = *updateEnvironmentApplicationNameArg
		if len(bucketInfo) > 0 {
			configOptions["AppBucket"] = bucketInfo[0]
			configOptions["AppKey"] = bucketInfo[1] + "/" + *updateEnvironmentVersionArg + ".zip"
		}
		configOptions["VersionLabel"] = *updateEnvironmentVersionArg
		configOptions["Tier"] = *updateEnvironmentTierArg

		usePreviousTemplate := true

		if *updateEnvironmentConfigArg != "" {
			configOptions = utils.CombineConfigOptions(utils.GetConfig(*updateEnvironmentConfigArg), configOptions)
			usePreviousTemplate = false
		}

		cfServcie := cloudformation.New(awsConfig)

		cfServcie.UpdateStack(configOptions, asset, usePreviousTemplate)
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
