package goconfig

import "os"

var (
	cfg any
)

// mergObject function merge 2 config objects in one
func mergeObject(dst any, src any) any {
	switch dst := dst.(type) {
	case map[string]any:
		switch src := src.(type) {
		case map[string]any:
			r := make(map[string]any)

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
func evaluateConfig(envCfg any) (any, bool) {
	switch envCfg := envCfg.(type) {
	case map[string]any:
		r := make(map[string]any)

		for k, v := range envCfg {
			value, status := evaluateConfig(v)
			if status {
				r[k] = value
			}
		}
		return r, true
	case []any:
		r := []any{}

		for i := range envCfg {
			value, status := evaluateConfig(envCfg[i])
			if status {
				r = append(r, value)
			}
		}
		return r, true
	case string:
		return os.LookupEnv(envCfg)
	default:
		return envCfg, true
	}
}
