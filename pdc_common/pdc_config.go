package pdc_common

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

type Lisensi struct {
	Email string `json:"email" yaml:"email"`
	Pwd   string `json:"pwd" yaml:"pwd"`
}

type PdcConfig struct {
	Lisensi    Lisensi `json:"lisensi" yaml:"lisensi"`
	ProjectID  string  `json:"project_id" yaml:"project_id"`
	Credential []byte  `json:"-" yaml:"-"`
	Version    string
	logname    string
	Hostname   string
}

func (conf *PdcConfig) CredOption() option.ClientOption {
	return option.WithCredentialsJSON(config.Credential)
}

var defaultCredential []byte

var config PdcConfig

func checkPath(fname string) (string, string) {
	exts := strings.Split(fname, ".")
	ext := exts[len(exts)-1]

	pwd, _ := os.Getwd()
	fullpath := filepath.Join(pwd, fname)

	if _, err := os.Stat(fullpath); errors.Is(err, os.ErrNotExist) {
		log.Info().Str(fullpath, "not exists")
	}

	return fullpath, ext

}

func SetConfig(fname string, version string, logname string, cred []byte) {
	namefull, ext := checkPath(fname)

	datas, err := os.ReadFile(namefull)
	if err != nil {
		panic(err)
	}

	if ext == "json" {
		json.Unmarshal(datas, &config)
	} else {
		yaml.Unmarshal(datas, &config)
	}

	if config.ProjectID == "" {
		config.ProjectID = "shopeepdc"
	}
	config.Version = version
	config.logname = logname

	hostname, _ := os.Hostname()
	config.Hostname = hostname

	if len(cred) >= 0 {
		config.Credential = cred
	} else {
		config.Credential = defaultCredential
	}

}

func GetConfig() *PdcConfig {
	return &config
}
