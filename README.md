# Elastic Beanstalk in GOlang!

## Install

## Commands

** Create Application
```
./bin/ebd create application -a <app-name>
-r <aws-region>
```

** Create Environment
```
./bin/ebd create environment -a <env_name> -e <env-name> -t <tier> -b <s3-bucket/s3-folder> -c <config-file> -l v1
-r <aws-region>
-p <aws-profile>
```

** Update Environment
```
./bin/ebd update environment -a <env_name> -e <env-name> -t <tier> -b <s3-bucket/s3-folder> -l v1

Optional Args:

-r <aws-region>
-p <aws-profile>
```

** List Application
```
./bin/ebd list applications

Optional Args:

-r <aws-region>
-v
```
Optional Args:


** List Environments
```
./bin/ebd list environments

Optional Args:

-a <app-name>
-r <aws-region>
-v
```

** Delete Application
```
./bin/ebd delete application -a <app-name>

Optional Args:

-r <aws-region>
```

** Delete Environment
```
./bin/ebd delete environment -a <app-name> -e <env-name> -t <worker>
```
