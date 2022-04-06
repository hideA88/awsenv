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

	"github.com/hideA88/awsenv/pkg"
	"github.com/hideA88/awsenv/pkg/aws"
	"github.com/hideA88/awsenv/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awsenv",
	Short: "set environment variables aws credentials and config",
	Long: `set environment variables aws credentials and config.
you need set -p profile name For example: awsenv -p dev`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger := Logger(verbose)
		ctx := pkg.Context{
			Context: context.TODO(), Logger: logger,
		}

		isVer, _ := cmd.Flags().GetBool("version")
		if isVer {
			version.Show(ctx)
		} else {
			profile, _ := cmd.Flags().GetString("profile")
			aws.Auth(ctx, profile)
		}
	},
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
	rootCmd.Flags().BoolP("verbose", "v", false, "show detail logs")
	rootCmd.Flags().Bool("version", false, "show version and system info")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
