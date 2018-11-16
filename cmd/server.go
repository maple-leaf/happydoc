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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/manifoldco/promptui"
	"github.com/maple-leaf/happydoc/helpers"
	"github.com/maple-leaf/happydoc/models"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run happydoc server, will init server when first run",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Println("checking if command docker-compose exist")
		_, err = helpers.RunShellCmd("which docker-compose", false)
		if err != nil {
			return
		}
		fmt.Println("DONE")
		data := []byte{}
		port := uint64(9876)
		questions := getServerQuestions(port)
		answers := []string{}
		var docServerConfig models.DocServerConfig

		for {
			answers, err = questions.Ask()
			if err != nil {
				break
			}

			port, err = strconv.ParseUint(answers[1], 10, 32)
			if err != nil {
				break
			}
			docServerConfig = models.DocServerConfig{
				Port:   port,
				PassWD: answers[2],
			}

			data, err = json.MarshalIndent(docServerConfig, "", "    ")
			if err != nil {
				break
			}

			fmt.Println(string(data))
			confirmed, _err := models.Confirm{
				Title:        "Is this looks good?",
				DefaultToYes: true,
			}.Ask()

			if confirmed == "N" {
				err = errors.New("cancel")
			} else if _err != nil {
				err = _err
			}

			if err == nil || err == promptui.ErrInterrupt {
				break
			}
		}

		if err == promptui.ErrInterrupt {
			err = errors.New("Canceled")
		}

		if err != nil {
			return
		}
		folder := answers[0]
		if folder != "" {
			_, _err := helpers.RunShellCmd("mkdir "+folder, false)
			if _err != nil {
				err = _err
				return
			}
			os.Chdir(folder)
		}

		sessionKey, err := generateSessionKey()
		if err != nil {
			return
		}

		content, err := fetchDockerCompose()
		if err != nil {
			return err
		}
		contentStr := string(content)
		contentStr = strings.Replace(contentStr, "${DB_PASSWD}", docServerConfig.PassWD, 2)
		contentStr = strings.Replace(contentStr, "${HAPPYDOC_PORT}", strconv.FormatUint(docServerConfig.Port, 10), 1)
		contentStr = strings.Replace(contentStr, "${SESSION_KEY}", strconv.FormatUint(sessionKey, 10), 1)
		fmt.Println(contentStr)
		composeFile, err := os.Create("docker-compose.yml")
		if err != nil {
			return err
		}
		_, err = composeFile.WriteString(contentStr)

		if err != nil {
			return err
		}

		fmt.Println("\n===starting server===")
		_, err = helpers.RunShellCmd("docker-compose up", false)

		return
	},
}

func fetchDockerCompose() (body []byte, err error) {
	uri := "https://raw.githubusercontent.com/maple-leaf/happydoc/master/server/docker-compose.yml"
	fmt.Printf("fetching docker-compose.yml from %v", uri)
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	return
}

func getServerQuestions(defaultServerPort uint64) (questions models.Questions) {
	questions = models.Questions{
		Items: []models.QA{},
	}
	questions.Items = append(questions.Items, models.Question{
		Title: "Init server at folder(leave blank will init at current folder, else create one)",
	})
	questions.Items = append(questions.Items, models.Question{
		Title:      "which port will this server run on",
		DefaultVal: strconv.FormatUint(defaultServerPort, 10),
	})
	questions.Items = append(questions.Items, models.Question{
		Title:      "set password for postgresql db in container(not postgresql on your host)",
		DefaultVal: "happydoc",
	})

	return
}

func generateSessionKey() (key string, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	key = id.String()

	return
}
