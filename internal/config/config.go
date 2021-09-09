package config

import cfg "gitlab.mvalley.com/datapack/cain/pkg/config"

type ConfigType struct {
	NSQConfig   NSQConfiguration
	MysqlConfig MysqlConfiguration
}

type MysqlConfiguration struct {
	RecruitmentDatapackMySQLConfig cfg.MySQLConfiguration
}

type NSQConfiguration struct {
	Host            string
	Port            string
	Topic           string
	ConsumerChannel string
}

var Config *ConfigType

// InitConfig ...
func InitConfig(configName string, configPaths []string) error {
	Config = new(ConfigType)
	return cfg.InitConfiguration(configName, configPaths, Config)
}

func GetConfig() ConfigType {
	return *Config
}
