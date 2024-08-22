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
		Write         string `yaml:"write"`
		MigrationUrl  string `yaml:"migration_url"`
		MigrationPath string `yaml:"migration_path"`
	} `yaml:"database"`

	Key struct {
		EncryptKey string `yaml:"encrypt_key"`
		JWT        string `yaml:"jwt"`
		Anonymous  string `yaml:"anonymous"`
	} `yaml:"key"`

	Api struct {
		BaseUrl string `yaml:"base_url"`
		TimeOut int32  `yaml:"timeout"`
	} `yaml:"api"`

	Assets struct {
		Common struct {
			MaxFileSize int `yaml:"max_file_size"`
		} `yaml:"common"`
		BaseUrl    string `yaml:"base_url"`
		UploadPath string `yaml:"upload_path"`
		ProfilePic struct {
			MaxSize int    `yaml:"max_size"`
			Path    string `yaml:"path"`
		} `yaml:"profile_pic"`
	} `yaml:"assets"`

	Midtrans struct {
		Timeout    int    `yaml:"timeout"`
		MerchantId string `yaml:"merchant_id"`
		ClientKey  string `yaml:"client_key"`
		ServerKey  string `yaml:"server_key"`
		BaseUrl    string `yaml:"base_url"`
		Api        struct {
			Charge  string `yaml:"charge"`
			Inquiry string `yaml:"inquiry"`
		} `yaml:"api"`
	} `yaml:"midtrans"`

	SMTP struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		SenderName   string `yaml:"sender_name"`
		AuthEmail    string `yaml:"auth_email"`
		AuthPassword string `yaml:"auth_password"`
	} `yaml:"smtp"`

	FCM struct {
		NotifUrl  string `yaml:"notif_url"`
		ClientKey string `yaml:"client_key"`
	} `yaml:"fcm"`
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
