package goconfig

import "os"

var (
	cfg interface{}
)

// mergObject function merge 2 config objects in one
func mergeObject(dst interface{}, src interface{}) interface{} {
	switch dst := dst.(type) {
	case map[string]interface{}:
		switch src := src.(type) {
		case map[string]interface{}:
			r := make(map[string]interface{})

			for k, v := range dst {
				r[k] = v
			}

			for k, v := range src {
				r[k] = mergeObject(dst[k], v)
			}
			return r
		default:
			return src
		}
	default:
		return src
	}
}

// evaluateConfig function evaluate all env variables in object
func evaluateConfig(envCfg interface{}) interface{} {
	switch envCfg := envCfg.(type) {
	case map[string]interface{}:
		r := make(map[string]interface{})

		for k, v := range envCfg {
			r[k] = evaluateConfig(v)
		}
		return r
	case []interface{}:
		r := []interface{}{}

		for i := range envCfg {
			r = append(r, evaluateConfig(envCfg[i]))
		}
		return r
	case string:
		return os.Getenv(envCfg)
	default:
		return envCfg
	}
}

func getConfigFile(host string) string {
	return configDir + host + ".json"
}
