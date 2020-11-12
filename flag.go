package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"unicode"
)

// IntValue implement the flag.Value interface
type IntValue struct {
	ptr reflect.Value
}

// Set assign ptr
func (i IntValue) Set(s string) error {
	if "" == s {
		return nil
	}
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	i.ptr.SetInt(v)
	return err
}

func (i IntValue) String() string {
	if k := i.ptr.Kind(); k == reflect.Invalid {
		return "0"
	}
	return strconv.FormatInt(i.ptr.Int(), 10)
}

// StringValue implement the flag.Value interface
type StringValue struct {
	ptr reflect.Value
}

// Set assign ptr
func (i StringValue) Set(s string) error {
	if "" == s {
		return nil
	}
	i.ptr.SetString(s)
	return nil
}

func (i StringValue) String() string {
	if k := i.ptr.Kind(); k == reflect.Invalid {
		return ""
	}
	return i.ptr.String()
}

// BoolValue implement the flag.Value interface
type BoolValue struct {
	ptr reflect.Value
}

// Set assign ptr
func (i BoolValue) Set(s string) error {
	if "" == s {
		return nil
	}
	v, err := strconv.ParseBool(s)
	i.ptr.SetBool(v)
	return err
}

func (i BoolValue) String() string {
	if k := i.ptr.Kind(); k == reflect.Invalid {
		return "false"
	}
	return strconv.FormatBool(i.ptr.Bool())
}

// FloatValue implement the flag.Value interface
type FloatValue struct {
	ptr reflect.Value
}

// Set assign ptr
func (i FloatValue) Set(s string) error {
	if "" == s {
		return nil
	}
	v, err := strconv.ParseFloat(s, 0)
	i.ptr.SetFloat(v)
	return err
}

func (i FloatValue) String() string {
	if k := i.ptr.Kind(); k == reflect.Invalid {
		return "0.0"
	}
	return strconv.FormatFloat(i.ptr.Float(), 'E', -1, 64)
}

// UintValue implement the flag.Value interface
type UintValue struct {
	ptr reflect.Value
}

// Set assign ptr
func (i UintValue) Set(s string) error {
	if "" == s {
		return nil
	}
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	i.ptr.SetUint(v)
	return err
}

func (i UintValue) String() string { return strconv.FormatUint(i.ptr.Uint(), 0) }

// Lower lowercase first letter
func Lower(s string) string {
	first := rune(s[0])
	return string(unicode.ToLower(first)) + s[1:]
}

var fields = make(map[string]reflect.Value)

// Flag only supports basic types of flag configuration
func Flag(name string, value interface{}, usage string) {
	valueField := reflect.ValueOf(value).Elem()
	flagVar(name, valueField, usage)
}

// FlagVar only supports struct type flag configuration
func FlagVar(config interface{}) {
	v := reflect.ValueOf(config).Elem()

	// Name of the struct tag used
	const tagName = "flag"

	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		typeField := v.Type().Field(i)

		name := Lower(typeField.Name)
		usage := typeField.Tag.Get(tagName)

		flagVar(name, valueField, usage)
	}
}

// FlagParse parse flag
func FlagParse(name string, usage string) {
	var configPath string
	flag.StringVar(&configPath, name, "", usage)

	flag.Parse()
	flag.Usage()

	defer func() {
		for k := range fields {
			delete(fields, k)
		}
	}()

	if "" == configPath {
		return
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var values map[string]interface{}
	err = json.Unmarshal(byteValue, &values)
	if err != nil {
		panic(err)
	}

	for k, field := range fields {
		if v, ok := values[k]; ok {
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val, ok := interface{}(v).(float64)
				if !ok {
					panic(k + " isn't of type int")
				}
				field.SetInt(int64(val))
			case reflect.String:
				val, ok := interface{}(v).(string)
				if !ok {
					panic(k + " isn't of type string")
				}
				field.SetString(val)
			case reflect.Bool:
				val, ok := interface{}(v).(bool)
				if !ok {
					panic(k + " isn't of type bool")
				}
				field.SetBool(val)
			case reflect.Float32, reflect.Float64:
				val, ok := interface{}(v).(float64)
				if !ok {
					panic(k + " isn't of type float")
				}
				field.SetFloat(val)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val, ok := interface{}(v).(float64)
				if !ok {
					panic(k + " isn't of type uint")
				}
				field.SetUint(uint64(val))
			}
		}
	}
}

func flagVar(name string, value reflect.Value, usage string) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		flag.Var(IntValue{value}, name, usage)
	case reflect.String:
		flag.Var(StringValue{value}, name, usage)
	case reflect.Bool:
		flag.Var(BoolValue{value}, name, usage)
	case reflect.Float32, reflect.Float64:
		flag.Var(FloatValue{value}, name, usage)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		flag.Var(FloatValue{value}, name, usage)
	}
	fields[name] = value
}
