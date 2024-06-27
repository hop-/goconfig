package goconfig

import (
	"encoding/json"
	"os"
	"path"
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
func Get[T any](objectName string) (*T, error) {
	value := cfg

	for _, k := range strings.Split(objectName, ".") {
		value = value.(map[string]any)[k]
	}

	return convertValue[T](value)
}

func GetAny(objectName string) any {
	value := cfg

	for _, k := range strings.Split(objectName, ".") {
		value = value.(map[string]any)[k]
	}

	return value
}

func convertValue[T any](value any) (*T, error) {
	buf, err := json.Marshal(value)

	if err != nil {
		return nil, err
	}

	var newValue T
	json.Unmarshal(buf, &newValue)

	return &newValue, nil
}

// Has returns true if config exist and false if not
func Has(objectName string) bool {
	value := cfg

	for _, k := range strings.Split(objectName, ".") {
		value = value.(map[string]any)[k]

		if value == nil {
			return false
		}
	}

	return true
}

func getConfigFile(host string) string {
	return path.Join(configDir, host+".json")
}

func loadDefaultConfig() error {
	defaultConfigPath := path.Join(configDir, defaultConfigName)
	file, err := os.Open(defaultConfigPath)
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

	var overwriteCfg any
	if err := json.NewDecoder(file).Decode(&overwriteCfg); err != nil {
		return err
	}

	// merge with existing config object
	cfg = mergeObject(cfg, overwriteCfg)
	return nil
}

// load custom environment configuration from file
func loadCustomEnvConfig() error {
	configPath := path.Join(configDir, customEnvConfigName)
	file, err := os.Open(configPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var envCfg any
	if err := json.NewDecoder(file).Decode(&envCfg); err != nil {
		return err
	}

	// evaluate env variables in config object
	envCfg, _ = evaluateConfig(envCfg)
	// merge with existing config object
	cfg = mergeObject(cfg, envCfg)
	return nil
}
