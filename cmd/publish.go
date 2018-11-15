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
	"bytes"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/maple-leaf/happydoc/models"
	"github.com/mholt/archiver"
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

		path := args[0]

		indexExist := isDocFolderHasIndexFile(path)
		if !indexExist {
			return fmt.Errorf("index.html not exist inside folder: %v", path)
		}

		setting.Path = path
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

		// _url, err := setting.ToURL()
		// if err != nil {
		// 	return err
		// }
		params := map[string]string{
			"version": setting.Version,
			"project": setting.Project,
			"type":    setting.DocType,
			"account": docConfig.Account,
		}
		zipPath := setting.Project + "_v" + setting.Version + ".zip"
		err := zipFolder(setting.Path, zipPath)
		if err != nil {
			return err
		}

		defer os.Remove(zipPath)

		uri := setting.Server + "/document/publish"
		req, _ := newFileUploadRequest(uri, params, "file", zipPath)
		req.Header.Add("X-Token", docConfig.Token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("Failed! Server return status code: %v", resp.StatusCode)
		}
		resp.Body.Close()

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

func isDocFolderHasIndexFile(docPath string) bool {
	return isFileExist(docPath + "/index.html")
}

func isFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func zipFolder(sourcePath string, destPath string) error {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}

	files, err := getFilesListInFolder(sourcePath)
	if err != nil {
		return err
	}

	err = z.Archive(files, destPath)

	return err
}

func getFilesListInFolder(folderPath string) (files []string, err error) {
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	return
}

// https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
// Creates a new file upload http request with optional extra params
func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
