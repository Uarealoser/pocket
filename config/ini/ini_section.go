package ini

import (
	"fmt"
	"strconv"
	"time"
)

type IniSection map[string]interface{}

func (sec IniSection) GetValue(key string) (value interface{}, err error) {
	value, ok := sec[key]
	if !ok {
		err = fmt.Errorf("no such key %s", key)
		return
	}
	return
}

func (sec IniSection) GetString(key string) (value string, err error) {
	v, err := sec.GetValue(key)
	if err != nil {
		return
	}
	value, ok := v.(string)
	if !ok {
		err = fmt.Errorf("%s is not a string and cann't parse to any other type", key)
		return
	}
	return
}

func (sec IniSection) GetSlice(key string) (value []string, err error) {
	v, err := sec.GetValue(key)
	if err != nil {
		return
	}
	var ok bool
	value, ok = v.([]string)
	if !ok {
		err = fmt.Errorf("%s is not a []string", key)
		return
	}
	return
}

func (sec IniSection) GetMap(key string) (value map[string]string, err error) {

	v, err := sec.GetValue(key)
	if err != nil {
		return
	}

	var ok bool
	value, ok = v.(map[string]string)
	if !ok {
		err = fmt.Errorf("%s is not a map[string]string", key)
		return
	}

	return
}

func (sec IniSection) GetInt(key string) (value int, err error) {

	v, err := sec.GetString(key)
	if err != nil {
		return
	}

	value, err = strconv.Atoi(v)
	return
}

func (sec IniSection) GetInt64(key string) (value int64, err error) {

	v, err := sec.GetString(key)
	if err != nil {
		return
	}

	value, err = strconv.ParseInt(v, 10, 64)
	return
}

func (sec IniSection) GetFloat64(key string) (value float64, err error) {

	v, err := sec.GetString(key)
	if err != nil {
		return
	}

	value, err = strconv.ParseFloat(v, 64)
	return
}

func (sec IniSection) GetBool(key string) (value bool, err error) {

	v, err := sec.GetString(key)
	if err != nil {
		return
	}

	value, err = strconv.ParseBool(v)
	return
}

func (sec IniSection) GetDuration(key string) (value time.Duration, err error) {

	v, err := sec.GetString(key)
	if err != nil {
		return
	}

	value, err = time.ParseDuration(v)
	return
}
