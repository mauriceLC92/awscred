# awscred

## What is this

The `awscred` package is an easy to use CLI for checking what [AWS](https://aws.amazon.com/) are available in your `.aws/credentials` file and which of these credentials are valid. Once a check is done, you can easily apply a valid profile and clean up credentials which are no longer of use.

## How to use

## Goal
- Be able to run `go run cmd/main.go apply workshop-admin echo $AWS_PROFILE` and see the profile AWS_PROFILE=workshop-admin
- be able to run that command with an environment variable set
