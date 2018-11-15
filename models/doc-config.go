package models

type DocConfig struct {
	Project string `json:"project"`
	Server  string `json:"server"`
	Account string `json:"account"`
	Token   string `json:"token"`
}
