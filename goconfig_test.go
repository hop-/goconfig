package goconfig

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadDefaultDirNotFound(t *testing.T) {
	err := Load()
	if err == nil {
		t.Fatal("Load() default should return error if config dir is not faund")
	}
}

func TestLoadFromEmptyDir(t *testing.T) {
	// Setting env
	os.Setenv("HOST_CONFIG_DIR", "test_data/empty")

	err := Load()
	if err == nil {
		t.Fatal("Load() should return error if default.json is not faund in dir")
	}
}

func TestLoadFromDefaultJson(t *testing.T) {
	// Setting env
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_empty_default")

	err := Load()
	if err != nil {
		t.Fatal("Load() should not return error if default.json is empty")
	}
}

func TestLoadWithEmptyHost(t *testing.T) {
	// Setting env
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_empty_host")
	os.Setenv("HOST_ENV", "some_host_name")

	err := Load()
	if err != nil {
		t.Fatal("Load() should not return error if host is empty")
	}
}

func TestLoadWithWrongOrNotExistingHost(t *testing.T) {
	// Setting env
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_empty_host")
	os.Setenv("HOST_ENV", "some_not_existing_host_name")

	err := Load()
	if err == nil {
		t.Fatal("Load() should return error if host is not faund")
	}
}

func TestLoadWithDefault(t *testing.T) {
	// Setting env
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")

	Load()

	some_int_as_float, err := Get[float32]("some_int")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_int_as_float != 100.0 {
		t.Errorf("Float should be 100.0 but got %v", *some_int_as_float)
	}

	some_int, err := Get[int32]("some_int")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_int != 100 {
		t.Errorf("Integer should be 100 but got %v", *some_int_as_float)
	}

	some_string, err := Get[string]("some_string")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_string != "test_string" {
		t.Errorf("String should be 'test_string' but got %v", *some_string)
	}

	some_array := GetAny("some_array")
	kind_of_some_array := reflect.TypeOf(some_array).Kind()
	if kind_of_some_array != reflect.Slice {
		t.Errorf("Is not slice but %v", kind_of_some_array)
	} else {
		len_of_some_array := len(some_array.([]any))
		if len_of_some_array != 4 {
			t.Errorf("Slice length should be 4 but got %v", len_of_some_array)
		}
	}

	some_array_of_any, err := Get[[]any]("some_array")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else {
		len_of_some_array := len(*some_array_of_any)
		if len_of_some_array != 4 {
			t.Errorf("Slice length should be 4 but got %v", len_of_some_array)
		}
	}

	some_object := GetAny("some_object")
	kind_of_some_object := reflect.TypeOf(some_object).Kind()
	if kind_of_some_object != reflect.Map {
		t.Errorf("Is not slice but %v", kind_of_some_object)
	} else {
		len_of_some_array := len(some_object.(map[string]any))
		if len_of_some_array != 2 {
			t.Errorf("Slice length should be 2 but got %v", len_of_some_array)
		}
	}

	some_object_of_defined_type, err := Get[struct {
		Key1 string `json:"key_1"`
		Key2 string `json:"key_2"`
	}]("some_object")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else {

		if some_object_of_defined_type.Key1 != "value_1" || some_object_of_defined_type.Key2 != "value_2" {
			t.Errorf(`Expected { Key1: "value_1", Key2: "value_2" } but got %+v`, *some_object_of_defined_type)
		}
	}
}

func TestLoadWithHost(t *testing.T) {
	// Setting env
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Setenv("HOST_ENV", "some_host_name")
	os.Setenv("ENV_VAR_NAME", "some_value")
	os.Setenv("ENV_BOOL_VAR", "true")
	os.Setenv("ENV_INT_VAR", "2")

	Load()

	some_int, err := Get[int]("some_int")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_int != 10 {
		t.Errorf("Integer should be 10 but got %v", *some_int)
	}

	some_string, err := Get[string]("some_string")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_string != "some_text" {
		t.Errorf("string should be 'some_text' but got %v", *some_string)
	}

	some_other_string, err := Get[string]("some_other_string")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_other_string != "some_other_text" {
		t.Errorf("string should be 'some_other_text' but got %v", *some_other_string)
	}

	some_env_config := GetAny("some_env_config")
	if some_env_config != "some_value" {
		t.Errorf("string should be 'some_value' but got %v", some_env_config)
	}

	some_bool_env_config, err := Get[bool]("some_bool_env_config")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if !*some_bool_env_config {
		t.Errorf("bool should be 'true' but got %v", *some_bool_env_config)
	}

	some_int_env_config, err := Get[int]("some_int_env_config")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_int_env_config != 2 {
		t.Errorf("bool should be 'true' but got %v", *some_int_env_config)
	}

	some_array := GetAny("some_array")
	kind_of_some_array := reflect.TypeOf(some_array).Kind()
	if kind_of_some_array != reflect.Slice {
		t.Errorf("Is not slice but %v", kind_of_some_array)
	} else {
		len_of_some_array := len(some_array.([]any))
		if len_of_some_array != 1 {
			t.Errorf("slice length should be 1 but got %v", len_of_some_array)
		}
	}

	some_array_of_any, err := Get[[]any]("some_array")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else {
		len_of_some_array := len(*some_array_of_any)
		if len_of_some_array != 1 {
			t.Errorf("slice length should be 1 but got %v", len_of_some_array)
		}
	}

	some_object := GetAny("some_object")
	kind_of_some_object := reflect.TypeOf(some_object).Kind()
	if kind_of_some_object != reflect.Map {
		t.Errorf("Is not slice but %v", kind_of_some_object)
	} else {
		len_of_some_array := len(some_object.(map[string]any))
		if len_of_some_array != 3 {
			t.Errorf("slice length should be 3 but got %v", len_of_some_array)
		}
	}

	some_object_of_defined_type, err := Get[struct {
		Key1 string `json:"key_1"`
		Key2 string `json:"key_2"`
		Key3 string `json:"key_3"`
	}]("some_object")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else {
		if some_object_of_defined_type.Key1 != "value_1" ||
			some_object_of_defined_type.Key2 != "new_value_2" ||
			some_object_of_defined_type.Key3 != "value_3" {
			t.Errorf(`Expected { Key1: "value_1", Key2: "new_value_2", Key3: "value_3" } but got %+v`, *some_object_of_defined_type)
		}
	}
}

func TestHasExistingKey(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	if !Has("some_int") {
		t.Error("Has() should return true for existing key")
	}
}

func TestHasNonExistingKey(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	if Has("non_existing_key") {
		t.Error("Has() should return false for non-existing key")
	}
}

func TestHasNestedKey(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	if !Has("some_object.key_1") {
		t.Error("Has() should return true for existing nested key")
	}
}

func TestHasNonExistingNestedKey(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	if Has("some_object.non_existing_key") {
		t.Error("Has() should return false for non-existing nested key")
	}
}

func TestGetMissingKey(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	_, err := Get[string]("non_existing_key")
	if err == nil {
		t.Error("Get() should return error for missing key")
	}
}

func TestGetNestedKey(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	val, err := Get[string]("some_object.key_1")
	if err != nil {
		t.Error("Get() should not return error for existing nested key", err)
	} else if *val != "value_1" {
		t.Errorf("Expected 'value_1' but got %v", *val)
	}
}

func TestGetAnyMissingKey(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	val := GetAny("non_existing_key")
	if val != nil {
		t.Errorf("GetAny() should return nil for missing key, got %v", val)
	}
}

func TestGetAnyEmptyPath(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	val := GetAny("")
	if val == nil {
		t.Error("GetAny() should return entire config for empty path")
	}
	kind := reflect.TypeOf(val).Kind()
	if kind != reflect.Map {
		t.Errorf("GetAny() with empty path should return map, got %v", kind)
	}
}

func TestExtractAllRequiredFieldsPresent(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type Config struct {
		SomeInt    int    `goconfig:"some_int"`
		SomeString string `goconfig:"some_string"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Errorf("Extract() should not return error when all required fields are present, got: %v", err)
	}
}

func TestExtractMissingRequiredField(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type Config struct {
		Missing string `goconfig:"non_existing_key"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err == nil {
		t.Error("Extract() should return error when required field is missing")
	}
}

func TestExtractMissingOptionalField(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type Config struct {
		Missing string `goconfig:"non_existing_key,optional"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Errorf("Extract() should not return error when optional field is missing, got: %v", err)
	}
}

func TestExtractFieldWithNoTag(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type Config struct {
		NoTag string
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Errorf("Extract() should not return error for fields without tag, got: %v", err)
	}
}

func TestExtractExplicitRequiredField(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type Config struct {
		Missing string `goconfig:"non_existing_key,required"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err == nil {
		t.Error("Extract() should return error when explicitly required field is missing")
	}
}

func TestLoadResetsConfig(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Setenv("HOST_ENV", "some_host_name")
	os.Setenv("ENV_VAR_NAME", "some_value")
	os.Setenv("ENV_BOOL_VAR", "true")
	os.Setenv("ENV_INT_VAR", "2")
	Load()

	// Now reload with different env (no host override)
	os.Unsetenv("HOST_ENV")
	Load()

	some_int, err := Get[int32]("some_int")
	if err != nil {
		t.Error("Error happened", err.Error())
	} else if *some_int != 100 {
		t.Errorf("After reload without host, integer should be 100 but got %v", *some_int)
	}
}

func TestConfigErrorMessage(t *testing.T) {
	err := &ConfigError{Message: "test error"}
	if err.Error() != "test error" {
		t.Errorf("ConfigError.Error() should return message, got: %v", err.Error())
	}
}

func TestExtractNestedStruct(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type ObjectConfig struct {
		Key1 string `goconfig:"key_1"`
		Key2 string `goconfig:"key_2"`
	}

	type Config struct {
		SomeInt    int          `goconfig:"some_int"`
		SomeObject ObjectConfig `goconfig:"some_object"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Fatalf("Extract() should not return error for nested struct, got: %v", err)
	}

	if obj.SomeInt != 100 {
		t.Errorf("Expected SomeInt=100 but got %v", obj.SomeInt)
	}
	if obj.SomeObject.Key1 != "value_1" {
		t.Errorf("Expected SomeObject.Key1='value_1' but got %v", obj.SomeObject.Key1)
	}
	if obj.SomeObject.Key2 != "value_2" {
		t.Errorf("Expected SomeObject.Key2='value_2' but got %v", obj.SomeObject.Key2)
	}
}

func TestExtractNestedStructPointer(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type ObjectConfig struct {
		Key1 string `goconfig:"key_1"`
		Key2 string `goconfig:"key_2"`
	}

	type Config struct {
		SomeObject *ObjectConfig `goconfig:"some_object"`
	}

	obj := &Config{SomeObject: &ObjectConfig{}}
	err := Extract("", obj)
	if err != nil {
		t.Fatalf("Extract() should not return error for nested struct pointer, got: %v", err)
	}

	if obj.SomeObject.Key1 != "value_1" {
		t.Errorf("Expected SomeObject.Key1='value_1' but got %v", obj.SomeObject.Key1)
	}
	if obj.SomeObject.Key2 != "value_2" {
		t.Errorf("Expected SomeObject.Key2='value_2' but got %v", obj.SomeObject.Key2)
	}
}

func TestExtractNestedStructMissingRequiredField(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type ObjectConfig struct {
		Key1       string `goconfig:"key_1"`
		MissingKey string `goconfig:"non_existing_key"`
	}

	type Config struct {
		SomeObject ObjectConfig `goconfig:"some_object"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err == nil {
		t.Error("Extract() should return error when nested struct has missing required field")
	}
}

func TestExtractNestedStructMissingOptionalField(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type ObjectConfig struct {
		Key1       string `goconfig:"key_1"`
		MissingKey string `goconfig:"non_existing_key,optional"`
	}

	type Config struct {
		SomeObject ObjectConfig `goconfig:"some_object"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Errorf("Extract() should not return error when nested struct has missing optional field, got: %v", err)
	}

	if obj.SomeObject.Key1 != "value_1" {
		t.Errorf("Expected SomeObject.Key1='value_1' but got %v", obj.SomeObject.Key1)
	}
}

func TestExtractDeeplyNestedStruct(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type InnerConfig struct {
		Key1 string `goconfig:"key_1"`
	}

	type OuterConfig struct {
		Inner InnerConfig `goconfig:"some_object"`
	}

	type Config struct {
		Outer OuterConfig `goconfig:""`
	}

	inner := &InnerConfig{}
	err := Extract("some_object", inner)
	if err != nil {
		t.Fatalf("Extract() should not return error for deeply nested struct, got: %v", err)
	}

	if inner.Key1 != "value_1" {
		t.Errorf("Expected Key1='value_1' but got %v", inner.Key1)
	}
}

func TestExtractNestedStructWithNoTag(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type ObjectConfig struct {
		NoTag string
		Key1  string `goconfig:"key_1"`
	}

	type Config struct {
		SomeObject ObjectConfig `goconfig:"some_object"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Errorf("Extract() should not return error when nested struct has field without tag, got: %v", err)
	}

	if obj.SomeObject.Key1 != "value_1" {
		t.Errorf("Expected SomeObject.Key1='value_1' but got %v", obj.SomeObject.Key1)
	}
	if obj.SomeObject.NoTag != "" {
		t.Errorf("Expected NoTag='' but got %v", obj.SomeObject.NoTag)
	}
}

func TestExtractPreservesUntaggedFields(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type Config struct {
		SomeInt int    `goconfig:"some_int"`
		NoTag   string // should remain zero value
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Fatalf("Extract() should not return error, got: %v", err)
	}

	if obj.SomeInt != 100 {
		t.Errorf("Expected SomeInt=100 but got %v", obj.SomeInt)
	}
	if obj.NoTag != "" {
		t.Errorf("Expected NoTag to be zero value but got %v", obj.NoTag)
	}
}

func TestExtractNestedStructReflectKind(t *testing.T) {
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Unsetenv("HOST_ENV")
	Load()

	type ObjectConfig struct {
		Key1 string `goconfig:"key_1"`
		Key2 string `goconfig:"key_2"`
	}

	type Config struct {
		SomeObject ObjectConfig `goconfig:"some_object"`
	}

	obj := &Config{}
	err := Extract("", obj)
	if err != nil {
		t.Fatalf("Extract() should not return error, got: %v", err)
	}

	kind := reflect.TypeOf(obj.SomeObject).Kind()
	if kind != reflect.Struct {
		t.Errorf("Expected SomeObject to be a struct but got %v", kind)
	}
}
