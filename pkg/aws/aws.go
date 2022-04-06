package aws

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hideA88/awsenv/pkg"
)

func Auth(ctx pkg.Context, profile string) {
	logger := ctx.Logger
	logger.Info("try authentication: ", profile)
	c, err := getCredentials(ctx, profile)
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}

	local, _ := time.LoadLocation("Local")
	expires := c.AWSExpires.In(local).Format(time.RFC3339)

	logger.Info("success.")
	logger.Info("AWS_PROFILE_NAME=", c.AwsProfileName)
	logger.Info("AWS_ACCESS_KEY_ID=", c.AwsAccessKeyId)
	logger.Info("AWS_SECRET_ACCESS_KEY=", c.AwsSecretAccessKey[:5]+"*****...")
	if len(c.AwsSessionToken) > 0 {
		logger.Info("AWS_SESSION_TOKEN=", c.AwsSessionToken[:10]+"...")
		logger.Info("AWS_SESSION_EXPIRES=", expires)
	}

	fmt.Printf("export AWS_PROFILE_NAME=%s\n", c.AwsProfileName)
	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", c.AwsAccessKeyId)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", c.AwsSecretAccessKey)
	if len(c.AwsSessionToken) > 0 {
		fmt.Printf("export AWS_SESSION_TOKEN=%s\n", c.AwsSessionToken)
		fmt.Printf("export AWS_SESSION_EXPIRES=%s\n", expires)
	}
}

func getCredentials(ctx pkg.Context, profile string) (*Credential, error) {
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
