package pdc_config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Lisensi struct {
	Email string `json:"email" yaml:"email"`
	Pwd   string `json:"pwd" yaml:"pwd"`
}

type BotConfig struct {
	Lisensi Lisensi `json:"lisensi" yaml:"lisensi"`
}

func LoadConfigFile(path string) BotConfig {
	var configbot BotConfig
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	ext := filepath.Ext(path)
	if ext == ".yml" {
		err = yaml.NewDecoder(file).Decode(&configbot)
	} else {
		err = json.NewDecoder(file).Decode(&configbot)
	}
	if err != nil {
		panic(err)
	}

	return configbot

}

var botConfig BotConfig
var botConfigOnce sync.Once

func GetConfig(path string) *BotConfig {

	botConfigOnce.Do(func() {
		botConfig = LoadConfigFile(path)
	})

	return &botConfig
}
