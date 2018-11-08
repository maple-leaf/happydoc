// Copyright © 2018 maple-leaf <tjfdfs.88@outlook.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/maple-leaf/happydoc/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var docConfig = models.DocConfig{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "happydoc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Join(args, ""))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(initCmd)

	rootCmd.AddCommand(serverCmd)

	rootCmd.AddCommand(publishCmd)
	initPublishCmd()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Search config in working directory with name ".happydoc.json"
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetConfigName(".happydoc")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		docConfig = models.DocConfig{
			Project: viper.GetString("project"),
			Server:  viper.GetString("server"),
		}
	}
}
