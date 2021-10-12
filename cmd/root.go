/*
Copyright © 2021 Hideaki Tarumi hideakit803@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/hideA88/awsenv/pkg"
	"github.com/hideA88/awsenv/pkg/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awsenv",
	Short: "set environment variables aws credentials and config",
	Long: `set environment variables aws credentials and config.
you need set -p profile name For example: awsenv -p dev`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger()
		ctx := pkg.Context{
			Context: context.TODO(), Logger: logger,
		}

		profile, _ := cmd.Flags().GetString("profile")

		logger.Info("try authentication: ", profile)
		c, err := aws.GetCredentials(ctx, profile)
		if err != nil {
			logger.Errorf("error: %v", err)
			return
		}

		local, _ := time.LoadLocation("Local")
		expires := c.AWSExpires.In(local).Format(time.RFC3339)

		logger.Info("success. credential ", c.AwsProfileName)
		logger.Info("AWS_ACCESS_KEY_ID=", c.AwsAccessKeyId)
		logger.Info("AWS_SECRET_ACCESS_KEY=", c.AwsSecretAccessKey[:5]+"*****...")
		logger.Info("AWS_SESSION_TOKEN=", c.AwsSessionToken[:10]+"...")
		logger.Info("AWS_SESSION_EXPIRES=", expires)

		fmt.Printf("export AWS_PROFILE_NAME=%s\n", c.AwsProfileName)
		fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", c.AwsAccessKeyId)
		fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", c.AwsSecretAccessKey)
		fmt.Printf("export AWS_SESSION_TOKEN=%s\n", c.AwsSessionToken)
		fmt.Printf("export AWS_SESSION_EXPIRES=%s\n", expires)
	},
}

func logger() *zap.SugaredLogger {
	rowlogger, _ := zap.NewProduction() //TODO implement

	//nolint:errcheck
	defer rowlogger.Sync() // flushes buffer, if any

	logger := rowlogger.Sugar()
	return logger
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().StringP("profile", "p", "default", "aws profile name")
	rootCmd.Flags().StringP("file", "f", "$HOME/.aws/credentials)", "credentials file location")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
