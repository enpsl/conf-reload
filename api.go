package conf_reload

import (
	"time"
)

var defaultEngine = NewEngine()

func LoadEngine(path string, opts ...Option) {
	err := defaultEngine.Load(path, opts...)
	if err != nil {
		panic(err)
	}
}

// Get external exposure api to get any type value.
func Get(key string) interface{} {
	return defaultEngine.Get(key)
}

// GetString external exposure api to get string type value.
func GetString(key string) string {
	return defaultEngine.GetString(key)
}

// GetBool external exposure api to get bool type value.
func GetBool(key string) bool {
	return defaultEngine.GetBool(key)
}

// GetInt external exposure api to get int type value.
func GetInt(key string) int {
	return defaultEngine.GetInt(key)
}

// GetInt64 external exposure api to get int64 type value.
func GetInt64(key string) int64 {
	return defaultEngine.GetInt64(key)
}

// GetFloat64 external exposure api to get float64 type value.
func GetFloat64(key string) float64 {
	return defaultEngine.GetFloat64(key)
}

// GetTime external exposure api to get time.Time type value.
func GetTime(key string) time.Time {
	return defaultEngine.GetTime(key)
}

// GetDuration external exposure api to get time.Duration type value.
func GetDuration(key string) time.Duration {
	return defaultEngine.GetDuration(key)
}

// GetStringSlice external exposure api to get []string type value.
func GetStringSlice(key string) []string {
	return defaultEngine.GetStringSlice(key)
}

// GetSlice external exposure api to get []interface{} type value.
func GetSlice(key string) []interface{} {
	return defaultEngine.GetSlice(key)
}

// GetStringMap external exposure api to get map[string]interface{} type value.
func GetStringMap(key string) map[string]interface{} {
	return defaultEngine.GetStringMap(key)
}

// GetStringMapString external exposure api to get map[string]string type value.
func GetStringMapString(key string) map[string]string {
	return defaultEngine.GetStringMapString(key)
}

// GetStringMapStringSlice external exposure api to get map[string][]string type value.
func GetStringMapStringSlice(key string) map[string][]string {
	return defaultEngine.GetStringMapStringSlice(key)
}

// DecodeToStruct The external exposure api is used for decoding,
// which can decode the value of the key map to the out variable
func DecodeToStruct(key string, out interface{}) error {
	return defaultEngine.DecodeToStruct(key, out)
}
