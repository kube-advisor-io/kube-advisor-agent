package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type MQTTConfig struct {
	Broker             string `yaml:"broker"`
	Topic              string `yaml:"topic"`
	Qos                int    `yaml:"qos"`
	TlsKeyFile         string `yaml:"tlsKeyFile"`
	TlsCertificateFile string `yaml:"tlsCertificateFile"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	CleanSession       bool   `yaml:"cleanSession"`
}

type Config struct {
	ClusterID         string     `yaml:"clusterID"`
	CustomerID        string     `yaml:"customerID"`
	DisabledProviders []string   `yaml:"disabledProviders"`
	MQTT              MQTTConfig `yaml:"mqtt"`
}

func ReadConfig() Config {
	viper.SetConfigName("default_config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/")      // path to look for the config file in
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
    viper.AddConfigPath("/etc/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
    viper.MergeInConfig()

	config := Config{}
	viper.Unmarshal(&config)
	return config
}
