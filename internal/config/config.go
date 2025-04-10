package config

import (
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"gopkg.in/yaml.v3"
	"os"
)

// Config конфигурация БД.
type Config struct {
	MSSQL MSSQLConfig `yaml:"MSSQL"`
	WLogg *wlogger.CustomWLogg
}

// MSSQLConfig настройки БД MSSQL.
type MSSQLConfig struct {
	MSSQLDriver   string `yaml:"MSSQLDriver"`
	MSSQLServer   string `yaml:"MSSQLServer"`
	MSSQLPort     string `yaml:"MSSQLPort"`
	MSSQLDatabase string `yaml:"MSSQLDatabase"`
	MSSQLUsername string `yaml:"MSSQLUsername"`
	MSSQLPassword string `yaml:"MSSQLPassword"`
}

func (c *Config) ConfigMSSQL(pathToYamlFile string) error {
	return c.unmarshallYaml(pathToYamlFile)
}

// unmarshallYaml загружает конфигурационный файл из YAML-файла в структуру Config.
func (c *Config) unmarshallYaml(pathToYaml string) error {
	yamlFile, err := os.ReadFile(pathToYaml)
	if err != nil {
		c.WLogg.LogE(msg.E3100, err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		c.WLogg.LogE(msg.E3101, err)
		return err
	}

	return nil
}
