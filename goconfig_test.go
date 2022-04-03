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

	some_int_as_float := Get("some_int")
	if some_int_as_float != 100.0 {
		t.Errorf("Integer should be 100 but got %v", some_int_as_float)
	}

	some_string := Get("some_string")
	if some_string != "test_string" {
		t.Errorf("string should be 'test_string' but got %v", some_string)
	}

	some_array := Get("some_array")
	kind_of_some_array := reflect.TypeOf(some_array).Kind()
	if kind_of_some_array != reflect.Slice {
		t.Errorf("Is not slice but %v", kind_of_some_array)
	} else {
		len_of_some_array := len(some_array.([]interface{}))
		if len_of_some_array != 4 {
			t.Errorf("slice lenght should be 4 but got %v", len_of_some_array)
		}
	}

	some_object := Get("some_object")
	kind_of_some_object := reflect.TypeOf(some_object).Kind()
	if kind_of_some_object != reflect.Map {
		t.Errorf("Is not slice but %v", kind_of_some_object)
	} else {
		len_of_some_array := len(some_object.(map[string]interface{}))
		if len_of_some_array != 2 {
			t.Errorf("slice lenght should be 2 but got %v", len_of_some_array)
		}
	}
}

func TestLoadWithHost(t *testing.T) {
	// Setting env
	os.Setenv("HOST_CONFIG_DIR", "test_data/with_configs_of_all_types")
	os.Setenv("HOST_ENV", "some_host_name")
	os.Setenv("ENV_VAR_NAME", "some_value")

	Load()

	var some_int int
	GetObject("some_int", &some_int)
	if some_int != 10 {
		t.Errorf("Integer should be 10 but got %v", some_int)
	}

	some_string := Get("some_string")
	if some_string != "some_text" {
		t.Errorf("string should be 'some_text' but got %v", some_string)
	}

	some_other_string := Get("some_other_string")
	if some_other_string != "some_other_text" {
		t.Errorf("string should be 'some_other_text' but got %v", some_other_string)
	}

	some_env_config := Get("some_env_config")
	if some_env_config != "some_value" {
		t.Errorf("string should be 'some_value' but got %v", some_env_config)
	}

	some_array := Get("some_array")
	kind_of_some_array := reflect.TypeOf(some_array).Kind()
	if kind_of_some_array != reflect.Slice {
		t.Errorf("Is not slice but %v", kind_of_some_array)
	} else {
		len_of_some_array := len(some_array.([]interface{}))
		if len_of_some_array != 1 {
			t.Errorf("slice lenght should be 1 but got %v", len_of_some_array)
		}
	}

	some_object := Get("some_object")
	kind_of_some_object := reflect.TypeOf(some_object).Kind()
	if kind_of_some_object != reflect.Map {
		t.Errorf("Is not slice but %v", kind_of_some_object)
	} else {
		len_of_some_array := len(some_object.(map[string]interface{}))
		if len_of_some_array != 3 {
			t.Errorf("slice length should be 3 but got %v", len_of_some_array)
		}
	}
}
