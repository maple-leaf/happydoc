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

	"github.com/spf13/viper"

	"github.com/manifoldco/promptui"
	"github.com/maple-leaf/happydoc/models"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate setting of happydoc",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		project, _ := getConfigFromPackageJSONFile()
		data := []byte{}

		for {
			docConfig, err = askQuestions(project)
			if err != nil {
				break
			}

			data, err = json.MarshalIndent(docConfig, "", "    ")
			if err != nil {
				break
			}

			fmt.Println(string(data))
			_, err = confirm("Is this looks good?", true)

			if err == nil || err == promptui.ErrInterrupt {
				break
			}
		}

		if err == promptui.ErrInterrupt {
			err = errors.New("Canceled")
		}

		if err == nil {
			ioutil.WriteFile(".happydoc.json", data, 0644)
		}

		return
	},
}

func getConfigFromPackageJSONFile() (project string, version string) {
	viper.SetConfigName("package")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	project = viper.GetString("name")
	version = viper.GetString("version")

	return
}

func askQuestions(project string) (docConfig models.DocConfig, err error) {
	project, err = askQuestionWithDefault("name of project", project)
	if err != nil {
		return
	}

	server, err := askQuestion("server api address")
	if err != nil {
		return
	}

	account, err := askQuestion("what's your account")
	if err != nil {
		return
	}

	token, err := askQuestion("what's your token")
	if err != nil {
		return
	}

	docConfig = models.DocConfig{
		Project: project,
		Server:  server,
		Account: account,
		Token:   token,
	}
	return
}

func askQuestion(question string) (answer string, err error) {
	if question == "" {
		return
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . | green }}: ",
		Invalid: "{{ . | red }}: ",
		Success: "{{ . | bold }}: ",
	}

	qa := promptui.Prompt{
		Label:     question,
		Templates: templates,
	}

	answer, err = qa.Run()

	return
}

func askQuestionWithDefault(question string, defaultVal string) (answer string, err error) {
	if question == "" {
		return "", nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . | green }}: ",
		Invalid: "{{ . | red }}: ",
		Success: "{{ . | bold }}: ",
	}

	qa := promptui.Prompt{
		Label:     question,
		Templates: templates,
		Default:   defaultVal,
	}

	answer, err = qa.Run()

	return
}

func confirm(question string, defaultToYes bool) (answer string, err error) {
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
		Default:   "Y",
	}

	answer, err = prompt.Run()

	return
}
