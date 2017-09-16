// Package config configuration related
package config

import (
	"encoding/json"
	"os"
	"path"
)

// GetConfig get config defined in config.json
func GetConfig() (config *Config) {
	pwd, _ := os.Getwd()
	path := path.Join(pwd, "config.json")
	configFile, err := os.Open(path)
	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		panic(err)
	}

	return
}

// MongoDB mongo config
type MongoDB struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,,omitempty"`
}

// Sender sender
type Sender struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// Mail Mail service config
type Mail struct {
	Receiver string `json:"receiver"`
	Sender   `json:"sender"`
}

// Config config entry
type Config struct {
	MongoDB `json:"mongodb"`
	Mail    `json:"mail"`
}
