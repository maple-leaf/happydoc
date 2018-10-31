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
	"net/url"

	"github.com/spf13/cobra"
)

type docSetting struct {
	server  string
	version string
	path    string
	project string
	docType string
}

func (setting docSetting) toURL() (string, error) {
	_url, err := url.Parse(setting.server)
	if setting.server == "" || err != nil {
		return "", errors.New("server is invalid")
	}

	query := _url.Query()
	if setting.version == "" {
		return "", errors.New("version is invalid")
	}
	if setting.project == "" {
		return "", errors.New("project is invalid")
	}

	query.Add("project", setting.project)
	query.Add("version", setting.version)
	query.Add("type", setting.docType)
	_url.RawQuery = query.Encode()

	return _url.String(), nil
}

var setting = docSetting{}
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish your document",
	Long:  `usage example: happydoc publish path/to/docs -s http://127.0.0.1:8000 -p awesomeProject -v 1.0.0`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("path to document folder not provided")
		}
		setting.path = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		version := cmd.Flag("version")
		if version != nil {
			setting.version = version.Value.String()
		}
		project := cmd.Flag("project")
		if project != nil {
			setting.project = project.Value.String()
		}
		server := cmd.Flag("server")
		if server != nil {
			setting.server = server.Value.String()
		}
		setting.docType = cmd.Flag("type").Value.String()

		_url, err := setting.toURL()
		if err != nil {
			return err
		}
		http.Get(_url)
		fmt.Printf("%v", _url)
		fmt.Printf("%v", setting)
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
