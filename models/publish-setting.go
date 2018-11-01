package models

import (
	"errors"
	"net/url"
)

type PublishSetting struct {
	Server  string
	Version string
	Path    string
	Project string
	DocType string
}

func (setting PublishSetting) ToURL() (string, error) {
	_url, err := url.Parse(setting.Server)
	if setting.Server == "" || err != nil {
		return "", errors.New("server is invalid")
	}

	query := _url.Query()
	if setting.Version == "" {
		return "", errors.New("version is invalid")
	}
	if setting.Project == "" {
		return "", errors.New("project is invalid")
	}

	query.Add("project", setting.Project)
	query.Add("version", setting.Version)
	query.Add("type", setting.DocType)
	_url.RawQuery = query.Encode()

	return _url.String(), nil
}
