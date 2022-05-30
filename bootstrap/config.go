package bootstrap

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Env  string `yaml:"env"`
	Host struct {
		Address string `yaml:"address"`
	} `yaml:"host"`

	Database struct {
		Write string `yaml:"write"`
	} `yaml:"database"`

	Key struct {
		EncryptKey string `yaml:"encrypt_key"`
		JWT        string `yaml:"jwt"`
	} `yaml:"key"`

	Api struct {
		TimeOut int32 `yaml:"timeout"`
	} `yaml:"api"`

	Assets struct {
		Url string `yaml:"url"`
	} `yaml:"assets"`
}

func LoadConfig(file string) (cnfg Config, err error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return cnfg, err
	}

	err = yaml.Unmarshal([]byte(yamlFile), &cnfg)
	if err != nil {
		return cnfg, err
	}

	if cnfg.Env == "" && os.Getenv("ENV") == "" {
		cnfg.Env = "local"
	}

	if os.Getenv("ENV") != "" {
		cnfg.Env = os.Getenv("ENV")
	}

	os.Setenv("ENV", cnfg.Env)

	return cnfg, err
}
