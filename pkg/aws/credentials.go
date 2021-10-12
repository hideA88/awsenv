package aws

import "time"

type Credential struct {
	AwsProfileName string
	AwsRegion      string
	AwsAccessKeyId string
	AwsSecretAccessKey string
	AwsSessionToken string
	AWSExpires time.Time
}