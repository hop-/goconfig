package goconfig

import (
	"encoding/json"
	"os"
	"strings"
)

var (
	configDir = "config"
)

const (
	defaultConfigName   = "default.json"
	customEnvConfigName = "custom-environment-variables.json"
)

// Load configurations
func Load() error {
	hostConfigDir := os.Getenv("HOST_CONFIG_DIR")
	if hostConfigDir != "" {
		configDir = hostConfigDir
	}

	if err := loadDefaultConfig(); err != nil {
		return err
	}

	host := os.Getenv("HOST_ENV")

	if host != "" {
		if err := loadFile(getConfigFile(host)); err != nil {
			return err
		}
	}

	return loadCustomEnvConfig()
}

// Get config
func Get(objectName string) interface{} {
	value := cfg

	for _, k := range strings.Split(objectName, ".") {
		value = value.(map[string]interface{})[k]
	}

	return value
}

// GetObject assign value to object with already structured
func GetObject(objectName string, object interface{}) error {
	value := Get(objectName)
	buf, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return json.Unmarshal(buf, object)
}

// Has return true if config exist and false if not
func Has(objectName string) bool {
	value := cfg

	for _, k := range strings.Split(objectName, ".") {
		value = value.(map[string]interface{})[k]

		if value == nil {
			return false
		}
	}

	return true
}

func loadDefaultConfig() error {
	file, err := os.Open(configDir + "/" + defaultConfigName)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&cfg)
}

// load configurtion form file
func loadFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var overwriteCfg interface{}
	if err := json.NewDecoder(file).Decode(&overwriteCfg); err != nil {
		return err
	}

	// merge with existing config object
	cfg = mergeObject(cfg, overwriteCfg)
	return nil
}

// load custon environment configuration from file
func loadCustomEnvConfig() error {
	file, err := os.Open(configDir + "/" + customEnvConfigName)
	if err != nil {
		return nil
	}
	defer file.Close()

	var envCfg interface{}
	if err := json.NewDecoder(file).Decode(&envCfg); err != nil {
		return err
	}

	// evaluate env variables in config object
	envCfg = evaluateConfig(envCfg)
	// merge with existing config object
	cfg = mergeObject(cfg, envCfg)
	return nil
}
