package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	Bucket             string `ini:"bucket"`
	AwsAccessKeyID     string `ini:"aws_access_key_id"`
	AwsSecretAccessKey string `ini:"aws_secret_access_key"`
	Region             string `ini:"region" default:"ru-central1"`
	EndpointUrl        string `ini:"endpoint_url" default:"https://storage.yandexcloud.net"`
}

func NewConfig() (*Config, error) {

	cfg, err := ini.Load("config/local.ini")
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	config := new(Config)
	err = cfg.Section("DEFAULT").MapTo(config)
	fmt.Print(*config)
	return config, nil
}

func (c *Config) UpdateConfig() {
	cfg, err := ini.Load("local.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	cfg.Section("DEFAULT").Key("bucket").SetValue(c.Bucket)
	cfg.Section("DEFAULT").Key("aws_access_key_id").SetValue(c.AwsAccessKeyID)
	cfg.Section("DEFAULT").Key("aws_secret_access_key").SetValue(c.AwsSecretAccessKey)
	cfg.SaveTo("local.ini")
}
