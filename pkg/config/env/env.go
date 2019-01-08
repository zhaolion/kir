package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Encode encode struct field with system env
func Encode(i interface{}) error {
	val := reflect.ValueOf(i)
	typ := val.Type()
	if kind := typ.Kind(); kind != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("pointer to struct type required")
	}
	return setStructValue(val.Elem())
}

func setStructValue(v reflect.Value) (err error) {
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		if err != nil {
			return err
		}
		f := typ.Field(i)
		fval := v.Field(i)
		tag := f.Tag.Get("env")
		if tag == "-" {
			continue
		}
		switch fval.Kind() {
		case reflect.Struct:
			err = setStructValue(fval)
		default:
			name, _ := ParseTag(tag)
			if name == "" {
				continue
			}
			if value := os.Getenv(name); value != "" {
				err = SetValue(fval, value)
			}
		}
	}
	return err
}

// SetValue ...
func SetValue(v reflect.Value, value string) (err error) {
	switch v.Kind() {
	case reflect.Bool:
		err = setBoolVal(v, value)
	case reflect.String:
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = setIntVal(v, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = setUintVal(v, value)
	case reflect.Float32, reflect.Float64:
		err = setFloatVal(v, value)
	case reflect.Array, reflect.Slice:
		err = setSliceVal(v, value)
	default:
		err = fmt.Errorf("unknown supported type: %s, %s", v.Kind(), v.Type().Name())
	}
	return err
}

// ParseTag ...
func ParseTag(tag string) (string, []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

func setBoolVal(v reflect.Value, value string) error {
	if value == "" {
		return nil
	} else if value == "true" {
		v.SetBool(true)
	} else if value == "false" {
		v.SetBool(false)
	} else {
		return fmt.Errorf("invalid bool value: %s", value)
	}
	return nil
}

func setIntVal(v reflect.Value, value string) error {
	x, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(x)
	return nil
}

func setUintVal(v reflect.Value, value string) error {
	x, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	v.SetUint(x)
	return nil
}

func setFloatVal(v reflect.Value, value string) error {
	x, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	v.SetFloat(x)
	return nil
}

func setSliceVal(v reflect.Value, value string) error {
	vv := strings.Split(value, ",")
	var nv reflect.Value
	if v.Kind() == reflect.Array {
		if v.Len() != len(vv) {
			return errors.New("unmatched length of array")
		}
		nv = v
	} else {
		nv = reflect.MakeSlice(v.Type(), len(vv), len(vv))
	}
	for i := 0; i < len(vv); i++ {
		if err := SetValue(nv.Index(i), vv[i]); err != nil {
			return err
		}
	}
	v.Set(nv)
	return nil
}
