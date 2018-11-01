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
	"errors"
	"fmt"
	"net/http"

	"github.com/maple-leaf/happydoc/models"
	"github.com/spf13/cobra"
)

var setting = models.PublishSetting{}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish your document",
	Long:  `usage example: happydoc publish path/to/docs -s http://127.0.0.1:8000 -p awesomeProject -v 1.0.0`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("path to document folder not provided")
		}
		setting.Path = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		setting.Server = docConfig.Server
		setting.Project = docConfig.Project

		version := cmd.Flag("version")
		if version != nil {
			setting.Version = version.Value.String()
		}
		project := cmd.Flag("project").Value.String()
		if project != "" {
			setting.Project = project
		}
		server := cmd.Flag("server").Value.String()
		if server != "" {
			setting.Server = server
		}
		setting.DocType = cmd.Flag("type").Value.String()

		_url, err := setting.ToURL()
		if err != nil {
			return err
		}
		http.Get(_url)
		fmt.Printf("%v", _url)
		return nil
	},
}

func initPublishCmd() {
	publishCmd.Flags().StringP("version", "v", "", "current document version")
	publishCmd.Flags().StringP("project", "p", "", "project that document belongs to")
	publishCmd.Flags().StringP("type", "t", "default", "document type")
	publishCmd.Flags().StringP("server", "s", "", "document server url")
	publishCmd.Flags().BoolP("force", "f", false, "force update version of document with given type and project")
}

func validPublishArgs(args []string) {

}
