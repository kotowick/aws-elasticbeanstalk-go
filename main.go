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
	//--
	// Common Flags
	//--
	//
	// General Flags
	awsRegion          = kingpin.Flag("region", "AWS Region").Short('r').Envar("AWS_REGION").String()
	awsAccessKeyID     = kingpin.Flag("access-key-id", "AWS Access Key ID").Short('k').Envar("AWS_ACCESS_KEY_ID").String()
	awsSecretAccessKey = kingpin.Flag("secret-access-key", "AWS Secret Access Key").Short('s').Envar("AWS_SECRET_ACCESS_KEY").String()
	awsProfile         = kingpin.Flag("profile", "Aws Profile").Short('p').Envar("AWS_PROFILE").String()
	awsCredPath        = kingpin.Flag("cred-path", "Aws Cred Path").Envar("AWS_CRED_PATH").String()
	verbose            = kingpin.Flag("verbose", "Verbose mode").Short('v').Envar("VERBOSE_MODE").Bool()

	// Environment commands
	environmentApplicationNameArg = kingpin.Flag("app-name", "Application Name").Envar("APP_NAME").Short('a').String()
	environmentNameArg            = kingpin.Flag("env-name", "Environment Name to create").Short('e').Envar("ENV_NAME").String()
	environmentTierArg            = kingpin.Flag("tier", "Environment Tier").Short('t').Envar("ENV_TIER").String()
	environmentS3Arg              = kingpin.Flag("bucket", "S3 bucket name").Short('b').Envar("ENV_S3_BUCKET_LOCATION").String()
	environmentConfigArg          = kingpin.Flag("config", "Config yml file").Short('c').Envar("ENV_CONFIG_PATH").String()
	environmentVersionArg         = kingpin.Flag("version-label", "Version").Short('l').Envar("ENV_VERSION_LABEL").String()
	environmentLocalFilePathArg   = kingpin.Flag("local-file-prefix", "The directory (prefix) of the file to upload.").Short('d').Envar("ENV_LOCAL_FILE_PREFIX").String()
	environmentDescriptionArg     = kingpin.Flag("desc", "Description of Version").Envar("ENV_DESCRIPTION").String()

	//--
	// Command Regex Section
	//--
	//
	createRegex   = regexp.MustCompile(`^create`)
	listRegex     = regexp.MustCompile(`^list`)
	updateRegex   = regexp.MustCompile(`^update`)
	deleteRegex   = regexp.MustCompile(`^delete`)
	rollbackRegex = regexp.MustCompile(`^rollback`)
	upsertRegex   = regexp.MustCompile(`^upsert`)

	//--
	// Global Vars
	//--
	//
	cliVersion = "1.0.3"

	// --
	// Top Level Commands
	// --
	//
	lists    = kingpin.Command("list", "Commands related to VPN connection profiles")
	create   = kingpin.Command("create", "Create EB")
	delete   = kingpin.Command("delete", "Destroy EB")
	update   = kingpin.Command("update", "Deploy EB")
	rollback = kingpin.Command("rollback", "Rollback EB")
	upsert   = kingpin.Command("upsert", "Upsert EB")

	// --
	// List
	// --
	//
	_                = lists.Command("applications", "List applications")
	listEnvironments = lists.Command("environments", "List environments")

	// --
	// Create
	// --
	//
	createApplication = create.Command("application", "Create Application")
	createEnvironment = create.Command("environment", "Create Environment")

	// --
	// Delete
	// --
	//
	deleteApplication = delete.Command("application", "Destroy application")
	deleteEnvironment = delete.Command("environment", "Destroy environment")

	// --
	// update
	// --
	//
	updateEnvironment = update.Command("environment", "Envrionment")
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
func listFunction(listMethod string, awsConfig config.Config, ebService elasticbeanstalk.ElasticBeanstalk) {
	switch listMethod {
	case "list applications":
		ebService.ListApplications(*verbose, true)
	case "list environments":
		if utils.VerifyParamatersWithAnd(true, map[string]string{"Application Name": *environmentApplicationNameArg}) {
			ebService.ListEnvironments(*verbose, true, *environmentApplicationNameArg, *environmentNameArg)
		}
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
func createFunction(createMethod string, awsConfig config.Config, ebService elasticbeanstalk.ElasticBeanstalk, s3Service s3.S3) {
	bucketInfo := s3Service.ParseS3Bucket(*environmentS3Arg)

	if !utils.VerifyParamatersWithAnd(true, map[string]string{"Application Name": *environmentApplicationNameArg}) {
		return
	}

	switch createMethod {
	case "create application":
		ebService.CreateApplication(*environmentApplicationNameArg)
	case "create environment":
		if !utils.VerifyParamatersWithAnd(true, map[string]string{"Local File Path": *environmentLocalFilePathArg, "Environment Name": *environmentNameArg, "Tier": *environmentTierArg, "Configuration File": *environmentConfigArg, "S3 Bucket/path": *environmentS3Arg, "Version": *environmentVersionArg}) {
			return
		}

		environment := utils.GetDefault(*environmentNameArg, *environmentApplicationNameArg)
		asset, err := Asset(fmt.Sprintf("resources/cloudformation/templates/%s_v1.json", *environmentTierArg))

		if err != nil {
			log.Fatalf("Asset not found: %s", err)
			return
		}

		s3Service.UploadSingleFile(*environmentS3Arg, *environmentLocalFilePathArg+"/"+*environmentVersionArg+".zip", *environmentVersionArg)

		additionalConfigOptions := make(map[string]string)
		additionalConfigOptions["EnvName"] = environment
		additionalConfigOptions["AppName"] = *environmentApplicationNameArg
		additionalConfigOptions["AppBucket"] = bucketInfo[0]
		additionalConfigOptions["AppKey"] = bucketInfo[1] + "/" + *environmentVersionArg + ".zip"
		additionalConfigOptions["VersionLabel"] = *environmentVersionArg
		additionalConfigOptions["Tier"] = *environmentTierArg

		configOptions := utils.GetConfig(*environmentConfigArg)
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
func deleteFunction(deleteMethod string, awsConfig config.Config, ebService elasticbeanstalk.ElasticBeanstalk) {
	switch deleteMethod {
	case "delete application":
		if utils.VerifyParamatersWithAnd(true, map[string]string{"Application Name": *environmentApplicationNameArg}) {
			ebService.DeleteApplication(*environmentApplicationNameArg)
		}
	case "delete environment":
		if utils.VerifyParamatersWithAnd(true, map[string]string{"Application Name": *environmentApplicationNameArg, "Environment Name": *environmentNameArg, "Tier": *environmentTierArg}) {
			cfServcie := cloudformation.New(awsConfig)
			cfServcie.DeleteStack(fmt.Sprintf("%s-%s-%s", *environmentApplicationNameArg, *environmentNameArg, *environmentTierArg))
		}
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
func updateFunction(updateMethod string, awsConfig config.Config, ebService elasticbeanstalk.ElasticBeanstalk, s3Service s3.S3) {
	bucketInfo := s3Service.ParseS3Bucket(*environmentS3Arg)

	switch updateMethod {
	case "update environment":
		if !utils.VerifyParamatersWithAnd(true, map[string]string{"Application Name": *environmentApplicationNameArg, "Environment Name": *environmentNameArg, "Tier": *environmentTierArg}) {
			return
		}

		if !utils.VerifyParamatersWithOr(false, map[string]string{"Configuration File": *environmentConfigArg, "Version": *environmentVersionArg}) {
			return
		}

		environment := utils.GetDefault(*environmentNameArg, *environmentApplicationNameArg)
		asset, err := Asset(fmt.Sprintf("resources/cloudformation/templates/%s_v1.json", *environmentTierArg))

		if err != nil {
			log.Fatalf("Asset not found: %s", err)
			return
		}

		if *environmentVersionArg != "" {
			s3Service.UploadSingleFile(*environmentS3Arg, *environmentLocalFilePathArg+"/"+*environmentVersionArg+".zip", *environmentVersionArg)
		}

		configOptions := make(map[string]string)
		configOptions["EnvName"] = environment
		configOptions["AppName"] = *environmentApplicationNameArg
		if len(bucketInfo) > 0 {
			configOptions["AppBucket"] = bucketInfo[0]
			configOptions["AppKey"] = bucketInfo[1] + "/" + *environmentVersionArg + ".zip"
		}
		configOptions["VersionLabel"] = *environmentVersionArg
		configOptions["Tier"] = *environmentTierArg

		usePreviousTemplate := true

		if *environmentConfigArg != "" {
			configOptions = utils.CombineConfigOptions(utils.GetConfig(*environmentConfigArg), configOptions)
			usePreviousTemplate = false
		}

		cfServcie := cloudformation.New(awsConfig)

		cfServcie.UpdateStack(configOptions, asset, usePreviousTemplate)
	}
}

func upsertFunction(updateMethod string, awsConfig config.Config, ebService elasticbeanstalk.ElasticBeanstalk, s3Service s3.S3) {
	if !ebService.ApplicationExists() {
		fmt.Println("App DOES NOT exist.. creating")
		createFunction("create application", awsConfig, ebService, s3Service)
	}

	if ebService.EnvironmentExists() {
		updateFunction("update environment", awsConfig, ebService, s3Service)
	} else {
		fmt.Println("Env DOES NOT exist..creating")
		createFunction("create environment", awsConfig, ebService, s3Service)
	}
}

//	Entry Controller
func main() {
	kingpin.Version(cliVersion)

	parsedArg := kingpin.Parse()

	if !utils.VerifyParamatersWithAnd(true, map[string]string{"AWS Region": *awsRegion}) {
		return
	}

	if !utils.VerifyParamatersWithOr(false, map[string]string{"Aws Profile": *awsProfile}) && !utils.VerifyParamatersWithAnd(false, map[string]string{"AWS Access Key ID": *awsAccessKeyID, "Aws Secret Access Key": *awsSecretAccessKey}) {
		return
	}

	awsConfig := config.New(
		*awsRegion,
		*awsAccessKeyID,
		*awsSecretAccessKey,
		*awsCredPath,
		*awsProfile,
	)

	s3Service := s3.New(*awsConfig)

	bucketInfo := [2]string{}

	if *environmentS3Arg != "" {
		bucketInfo = s3Service.ParseS3Bucket(*environmentS3Arg)
	}

	// if Application does not exist, create it
	ebService := elasticbeanstalk.New(
		*environmentApplicationNameArg, *environmentNameArg,
		*environmentVersionArg, *environmentDescriptionArg,
		bucketInfo, *environmentVersionArg+".zip",
		*environmentTierArg, *awsConfig,
	)

	switch {
	case createRegex.MatchString(parsedArg):
		createFunction(parsedArg, *awsConfig, *ebService, *s3Service)
	case listRegex.MatchString(parsedArg):
		listFunction(parsedArg, *awsConfig, *ebService)
	case updateRegex.MatchString(parsedArg):
		updateFunction(parsedArg, *awsConfig, *ebService, *s3Service)
	case deleteRegex.MatchString(parsedArg):
		deleteFunction(parsedArg, *awsConfig, *ebService)
	case upsertRegex.MatchString(parsedArg):
		upsertFunction(parsedArg, *awsConfig, *ebService, *s3Service)
	default:
		log.Fatalf("Command signature not recognized: %s", parsedArg)
	}
}
