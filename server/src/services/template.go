package services

import (
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"time"

	"github.com/gin-contrib/multitemplate"
)

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	funcMap := template.FuncMap{
		"formatAsDate":     formatAsDate,
		"lastIndexOfArray": lastIndexOfArray,
	}

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	pages, err := filepath.Glob(templatesDir + "/pages/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, page := range pages {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, page)
		r.AddFromFilesFuncs(filepath.Base(page), funcMap, files...)
	}
	return r
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func lastIndexOfArray(in interface{}) int {
	inVal := reflect.ValueOf(in)
	inType := inVal.Type()
	if inType.Kind() == reflect.Slice {
		length := inVal.Len()
		return length - 1
	}

	return -1
}
