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
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/maple-leaf/happydoc/helpers"
	"github.com/maple-leaf/happydoc/models"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run happydoc server, will init server when first run",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		data := []byte{}
		port := uint64(9876)
		questions := getServerQuestions(port)
		answers := []string{}

		for {
			answers, err = questions.Ask()
			if err != nil {
				break
			}

			port, err = strconv.ParseUint(answers[1], 10, 32)
			if err != nil {
				break
			}
			docServerConfig := models.DocServerConfig{
				Port: port,
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

		if err == nil {
			folder := answers[0]
			if folder != "" {
				_, _err := helpers.RunShellCmd("mkdir "+folder, false)
				if _err != nil {
					err = _err
					return
				}
				os.Chdir(folder)
			}
			// TODO: get docker config; init docker; start server
		}

		return
	},
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

	return
}
