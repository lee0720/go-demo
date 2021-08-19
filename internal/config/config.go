package config

import cfg "gitlab.mvalley.com/adam/common/pkg/config"

type ConfigType struct {
	NSQConfig NSQConfiguration
}

type NSQConfiguration struct {
	Host            string
	Port            string
	Topic           string
	ConsumerChannel string
}

var c *ConfigType

// InitConfig ...
func InitConfig(configName string, configPaths []string) error {
	c = new(ConfigType)
	return cfg.InitConfiguration(configName, configPaths, c)
}

func GetConfig() ConfigType {
	return *c
}
