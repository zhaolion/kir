package env

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestParamValidate(t *testing.T) {
	unsupported := []interface{}{
		0,
		"",
		[]string{},
		map[string]string{},
		true,
	}
	for _, v := range unsupported {
		if err := Encode(v); err == nil {
			t.Error("expected: pointer to struct type required error")
		}
		if err := Encode(&v); err == nil {
			t.Error("expected: pointer to struct type required error")
		}
	}
}

func TestEncode(t *testing.T) {
	type AppMode string
	type Basic struct {
		Mode AppMode `env:"struct1_mode"`
		Addr string  `env:"struct1_addr"`
	}
	type struct1 struct {
		Basic
		Name   string    `env:"struct1_name"`  // string
		Hosts  []string  `env:"struct1_hosts"` // slice
		Hosts1 [3]string `env:"struct1_hosts"` // array
		Port   int       `env:"struct1_port"`  // int
		Debug  bool      `env:"struct1_debug"` // bool
		Rate   float64   `env:"struct1_rate"`  // float
		Omit   string    `env:"-"`             // omit
		DB     struct {
			Drive string   `env:"struct1_db_drive"`
			Hosts []string `env:"struct1_db_hosts"`
			Port  uint     `env:"struct1_db_port"`
		}
	}
	envs := map[string]string{
		"struct1_mode":     "test",
		"struct1_addr":     "localhost",
		"struct1_name":     "app",
		"struct1_hosts":    "192.168.33.1,192.168.33.2,192.168.33.3",
		"struct1_port":     "8080",
		"struct1_debug":    "true",
		"struct1_rate":     "0.618",
		"struct1_omit":     "omit",
		"struct1_db_drive": "mysql",
		"struct1_db_hosts": "1,2,3",
		"struct1_db_port":  "3306",
	}
	value := struct1{
		Name:   "app",
		Hosts:  []string{"192.168.33.1", "192.168.33.2", "192.168.33.3"},
		Hosts1: [3]string{"192.168.33.1", "192.168.33.2", "192.168.33.3"},
		Port:   8080,
		Debug:  true,
		Rate:   0.618,
		Omit:   "",
	}
	value.Mode = "test"
	value.Addr = "localhost"
	value.DB.Drive = "mysql"
	value.DB.Hosts = []string{"1", "2", "3"}
	value.DB.Port = 3306

	for k, v := range envs {
		_ = os.Setenv(k, v)
	}
	var s struct1
	err := Encode(&s)
	if err != nil {
		t.Errorf("expected: %v, got: %v", nil, err)
	}
	if err != nil {
		return
	}
	eq, err := equal(s, value)
	if err != nil {
		t.Fatal(err)
	}
	if !eq {
		t.Errorf("expected: %+v, got: %+v", value, s)
	}
}

func TestSetBoolVal(t *testing.T) {
	type TC struct {
		Value  string
		HasErr bool
		Bool   bool
	}
	tcs := []TC{
		TC{"", false, false},
		TC{"true", false, true},
		TC{"false", false, false},
		TC{"True", true, false},
		TC{"True", true, false},
		TC{"TRUE", true, false},
		TC{"1", true, false},
	}
	for _, tc := range tcs {
		var b bool
		v := reflect.ValueOf(&b).Elem()
		err := setBoolVal(v, tc.Value)
		hasErr := err != nil
		if tc.HasErr != hasErr {
			t.Errorf("expected: %t, got: %t", tc.HasErr, hasErr)
		}
		b = v.Bool()
		if b != tc.Bool {
			t.Errorf("expected: %t, got: %t", tc.Bool, b)
		}
	}
}

func equal(s1, s2 interface{}) (bool, error) {
	data1, err := json.Marshal(s1)
	if err != nil {
		return false, err
	}
	data2, err := json.Marshal(s2)
	if err != nil {
		return false, err
	}

	return string(data1) == string(data2), nil
}
