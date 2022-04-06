package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hideA88/awsenv/pkg"
)

func GetCredentials(ctx pkg.Context, profile string) (*Credential, error) {
	cfg, err := config.LoadSharedConfigProfile(ctx.Context, profile)
	if err != nil {
		return nil, err
	}

	var cred *aws.Credentials
	if cfg.RoleARN != "" {
		cred, err = assumeRole(ctx, cfg)
		if err != nil {
			ctx.Logger.Errorf("assumeRole error: %v", err)
			return nil, err
		}
	} else {
		cred, err = loadConfig(ctx, profile)
		if err != nil {
			ctx.Logger.Errorf("load config error: %v", err)
			return nil, err
		}
	}

	return &Credential{
		AwsProfileName:     profile,
		AwsAccessKeyId:     cred.AccessKeyID,
		AwsSecretAccessKey: cred.SecretAccessKey,
		AwsSessionToken:    cred.SessionToken,
		AWSExpires:         cred.Expires,
	}, nil
}

func loadConfig(ctx pkg.Context, profile string) (*aws.Credentials, error) {
	c, err := config.LoadDefaultConfig(ctx.Context, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}
	cred, err := c.Credentials.Retrieve(ctx.Context)
	if err != nil {
		return nil, err
	}

	return &cred, nil
}

func assumeRole(ctx pkg.Context, scfg config.SharedConfig) (*aws.Credentials, error) {
	//TODO implement MFASerialがちゃんとあるかバリデーションはしたほうがいいかも
	var profile string
	if scfg.SourceProfileName != "" {
		profile = scfg.SourceProfileName
	} else {
		profile = scfg.Profile
	}

	cfg, err := config.LoadDefaultConfig(ctx.Context,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(scfg.Region),
	)

	if err != nil {
		return nil, err
	}

	client := sts.NewFromConfig(cfg)

	creds := stscreds.NewAssumeRoleProvider(client, scfg.RoleARN,
		func(o *stscreds.AssumeRoleOptions) {
			o.SerialNumber = aws.String(scfg.MFASerial)
			o.TokenProvider = myStdinTokenProvider
		})
	c, err := creds.Retrieve(ctx.Context)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func myStdinTokenProvider() (string, error) {
	var v string
	//標準出力の場合、evalの評価対象になってしまうため
	fmt.Fprintf(os.Stderr, "Assume Role MFA token code: ")
	_, err := fmt.Scanln(&v)

	return v, err
}
