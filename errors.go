package goconfig

// ConfigError represents an error related to configuration loading or extraction.
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
