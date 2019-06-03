package config

import (
	"encoding/json"
	"os"
)

// Config represents the configuration data
type Config struct {
	VolonsPlatform string `json:"volons_platform"`
	HTTPAddr       string `json:"http"`
	Database       string `json:"database"`
}

// Init conf with defaults
var _conf = Config{
	VolonsPlatform: "", //"https://api.volons.fr/gcs",
	HTTPAddr:       "0.0.0.0:8081",
	Database:       "./database/",
}

// Get returns the global config
func Get() Config {
	return _conf
}

// Read reads the global config from a json file
func Read(filePath string) {
	if filePath == "" {
		_conf.VolonsPlatform = getEnv("VOLONS_PLATFORM", _conf.VolonsPlatform)
		_conf.HTTPAddr = getEnv("VOLONS_HTTP", _conf.HTTPAddr)
		_conf.Database = getEnv("VOLONS_DATABASE", _conf.Database)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&_conf); err != nil {
		panic(err.Error())
	}
}

func getEnv(name string, defaultVal string) string {
	val := os.Getenv(name)

	if val == "" {
		return defaultVal
	}

	return val
}
