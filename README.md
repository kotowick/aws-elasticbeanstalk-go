# Elastic Beanstalk in GOlang!

## Install
1. You can download the binary directly and use it
2. You can download the docker images and use it (see at the very bottom for instructions)

## Commands

**Create Application**
```
./bin/ebd create application -a <app-name>

Optional Args:

-r <aws-region>
```

**Create Environment**
```
./bin/ebd create environment -a <env_name> -e <env-name> -t <tier> -b <s3-bucket/s3-folder> -c <config-file> -l v1

Optional Args:

-r <aws-region>
-p <aws-profile>
```

**Update Environment**
```
./bin/ebd update environment -a <env_name> -e <env-name> -t <tier> -b <s3-bucket/s3-folder> -l v1

Optional Args:

-r <aws-region>
-p <aws-profile>
```

**List Application**
```
./bin/ebd list applications

Optional Args:

-r <aws-region>
-v
```
Optional Args:


**List Environments**
```
./bin/ebd list environments

Optional Args:

-a <app-name>
-r <aws-region>
-v
```

**Delete Application**
```
./bin/ebd delete application -a <app-name>

Optional Args:

-r <aws-region>
```

**Delete Environment**
```
./bin/ebd delete environment -a <app-name> -e <env-name> -t <worker>
```

## Docker

Four things are required:

1. the docker image
  ```
  docker pull kotowick/go-deploy:latest
  ```
2. environment variables
  ```
  export APP_NAME='' # the Elastic Beanstalk Application Name
  export ENV_NAME='' # the Elastic Beanstalk environment Name
  export ENV_TIER='' # can either be (worker | webserver)
  export ENV_S3_BUCKET_LOCATION='{{bucket}}/{{app_name}}/{{versions}}' # example - S3 location to upload the version (.zip)
  export ENV_CONFIG_PATH='/tmp/vpc9838983.yml' # example - config settings for EB itself
  export ENV_LOCAL_FILE_PREFIX='/tmp' # can be anywhere, but has to match the VOLUME_PATH variable in the command to run the container
  export AWS_REGION=""
  export AWS_ACCESS_KEY_ID=""
  export AWS_SECRET_ACCESS_KEY=""
  export ENV_VERSION_LABEL="v1.0" # same name as the zip file
  ```
  
  These can either be set when running the container `(using -e options)`, or by having a `*.env` file under the `VOLUME_PATH`. The container will `source` any `*.env` files located under `ENV_LOCAL_FILE_PREFIX`.
 
3. EB config file (contains settings for the EB env itself.. such as `VPC ID, subnets...`)
4. the `version` of the file you want `deploy`
  ```
  if the version is v1.0, then a .zip file must match that name under ENV_LOCAL_FILE_PREFIX (listed above)
  ```

The startup command for this container is `go-deploy upsert`, which creates resoursed if they don't exist, and updates them if they do exist.


