package models

import (
	"github.com/manifoldco/promptui"
)

type QA interface {
	Ask() (answer string, err error)
}

type Question struct {
	Title      string
	DefaultVal string
}

func (q Question) Ask() (answer string, err error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . | green }}: ",
		Invalid: "{{ . | red }}: ",
		Success: "{{ . | bold }}: ",
	}

	qa := promptui.Prompt{
		Label:     q.Title,
		Templates: templates,
		Default:   q.DefaultVal,
	}

	answer, err = qa.Run()

	return
}

type Confirm struct {
	Title        string
	DefaultToYes bool
}

func (c Confirm) Ask() (answer string, err error) {
	defaultVal := "N"
	if c.DefaultToYes {
		defaultVal = "Y"
	}

	prompt := promptui.Prompt{
		Label:     c.Title,
		IsConfirm: true,
		Default:   defaultVal,
	}

	_, err = prompt.Run()
	if err != nil && err != promptui.ErrInterrupt {
		answer = "N"
		err = nil
	} else {
		answer = "Y"
	}

	return
}

type Questions struct {
	Items []QA
}

func (qs Questions) Ask() (answers []string, err error) {
	answer := ""
	for _, item := range qs.Items {
		answer, err = item.Ask()
		if err != nil {
			return
		}
		answers = append(answers, answer)
	}

	return
}
