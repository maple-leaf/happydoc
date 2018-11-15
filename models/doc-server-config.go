package models

type DocServerConfig struct {
	Port   uint64 `json:"port"`
	PassWD string `json:"password"`
}
