package goconfig

import (
	"encoding/json"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

var (
	configDir = "config"
)

const (
	defaultConfigName   = "default.json"
	customEnvConfigName = "custom-environment-variables.json"
	tagName             = "goconfig"
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

// Extract config to annotated struct
func Extract(root string, object any) error {
	typeOfObject := reflect.TypeOf(object).Elem()
	valueOfObject := reflect.ValueOf(object).Elem()

	// Iterate over struct fields
	for i := 0; i < typeOfObject.NumField(); i++ {
		typeField := typeOfObject.Field(i)

		// Get the tag value for the field
		tag := typeField.Tag.Get(tagName)

		// Parse the tag to get the config path and optional flag
		path, optional := parseTag(tag)

		// If path is empty, skip this field
		if path == "" {
			continue
		}

		if root != "" {
			path = root + "." + path
		}

		if Has(path) {
			valueOfField := valueOfObject.Field(i)

			if valueOfField.Kind() == reflect.Pointer {
				valueOfField = valueOfField.Elem()
			}

			if valueOfField.Kind() == reflect.Struct {
				// If the field is a struct, recursively extract its fields
				err := Extract(path, valueOfField.Addr().Interface())
				if err != nil {
					return err
				}
			} else {
				// Extract the config value and set it to the struct field
				err := extractValueByPath(path, valueOfField.Addr().Interface())
				if err != nil {
					return err
				}
			}
		} else if !optional {
			return &ConfigError{Message: "required config " + path + " is missing"}
		}
	}

	return nil
}

// Get config
func Get[T any](objectName string) (*T, error) {
	value := new(T)

	err := extractValueByPath(objectName, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// GetAny config as any
func GetAny(objectName string) any {
	value, _ := resolveValueByPath(objectName)

	return value
}

// Has returns true if config exist and false if not
func Has(objectName string) bool {
	_, ok := resolveValueByPath(objectName)

	return ok
}

func resolveValueByPath(path string) (any, bool) {
	value := cfg
	if path == "" {
		return value, true
	}

	for _, k := range strings.Split(path, ".") {
		var ok bool
		value, ok = value.(map[string]any)[k]

		if !ok {
			return nil, false
		}
	}

	return value, true
}

func extractValueByPath(path string, valueObj any) error {
	value, ok := resolveValueByPath(path)

	if !ok {
		return &ConfigError{Message: "config " + path + " is missing"}
	}

	return convertValue(value, valueObj)
}

// convertValue converts value to given type
func convertValue(value any, valueObj any) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, valueObj)
	if err != nil {
		// try to mitigate quoted values from custom envs
		s, unquoteErr := strconv.Unquote(string(buf))
		if unquoteErr != nil {
			return err
		}

		err = json.Unmarshal([]byte(s), valueObj)
		if err != nil {
			return err
		}
	}

	return nil
}

// getConfigFile returns file path by host
func getConfigFile(host string) string {
	return path.Join(configDir, host+".json")
}

// loadDefaultConfig loads default config
func loadDefaultConfig() error {
	defaultConfigPath := path.Join(configDir, defaultConfigName)
	file, err := os.Open(defaultConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&cfg)
}

// loadFile loads configurtion from file
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

// loadCustomEnvConfig loads custom environment configuration from file
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

func parseTag(tag string) (string, bool) {
	parts := strings.Split(tag, ",")
	path := strings.TrimSpace(parts[0])
	optional := false

	if len(parts) > 1 {
		for _, part := range parts[1:] {
			part := strings.TrimSpace(part)
			if part == "optional" {
				optional = true
				break
			} else if part == "required" {
				optional = false
				break
			}
		}
	}

	return path, optional
}
