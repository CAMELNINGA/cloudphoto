package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

func NewConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error while getting home dir")
	}
	cfg, err := ini.Load(home + "/.config/cloudphoto/cloudphotorc")
	if err != nil {
		config := new(Config)
		return config, nil
	}
	config := new(Config)
	err = cfg.Section("DEFAULT").MapTo(config)
	return config, nil
}

func InitConfig(bucket, awskey, awssec string) (*Config, error) {
	home, err := os.UserHomeDir()
	path := home + "/.config/cloudphoto"
	if err := os.MkdirAll(path, 0770); err != nil {
		return nil, fmt.Errorf("Error creating dir %s \n", err)
	}
	f, err := os.OpenFile(home+"/.config/cloudphoto/cloudphotorc", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			f, err = os.Open(home + "/.config/cloudphoto/cloudphotorc")
			if err != nil {
				return nil, fmt.Errorf("Error open file %s \n", err)
			}
		} else {
			return nil, fmt.Errorf("Error open file %s \n", err)
		}
	}
	f.Close()
	config := new(Default)
	config.AwsAccessKeyID = awskey
	config.AwsSecretAccessKey = awssec
	config.Bucket = bucket
	config.Region = "ru-central1"
	config.EndpointUrl = "https://storage.yandexcloud.net"
	def := &Config{config}

	cfg := ini.Empty()
	err = cfg.ReflectFrom(def)
	cfg.SaveTo(home + "/.config/cloudphoto/cloudphotorc")
	return def, nil
}

type Config struct {
	*Default
}
type Default struct {
	Bucket             string `ini:"bucket"`
	AwsAccessKeyID     string `ini:"aws_access_key_id"`
	AwsSecretAccessKey string `ini:"aws_secret_access_key"`
	Region             string `ini:"region" default:"ru-central1"`
	EndpointUrl        string `ini:"endpoint_url" default:"https://storage.yandexcloud.net"`
}
