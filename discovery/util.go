package discovery

import (
	"fmt"
	"strconv"
)

func getTimeoutFromMap(data map[string]interface{}) (
	result Timeout, err error) {

	connTimeoutMs, err := getIntFromMap(data, "connTimeoutMs")
	if err != nil {
		connTimeoutMs = ServiceDefaultTimeout
	}

	readTimeoutMs, err := getIntFromMap(data, "readTimeoutMs")
	if err != nil {
		readTimeoutMs = ServiceDefaultTimeout
	}

	writeTimeoutMs, err := getIntFromMap(data, "writeTimeoutMs")
	if err != nil {
		writeTimeoutMs = ServiceDefaultTimeout
	}

	err = nil
	result = Timeout{
		ConnTimeoutMs:  uint32(connTimeoutMs),
		ReadTimeoutMs:  uint32(readTimeoutMs),
		WriteTimeoutMs: uint32(writeTimeoutMs),
	}
	return
}

func getMapFromMap(data map[string]interface{}, key string) (
	result map[string]interface{}, err error) {

	value, ok := data[key]
	if !ok {
		err = fmt.Errorf("invalid config, not found:%s", key)
		return
	}

	switch inst := value.(type) {
	case map[string]interface{}:
		result = inst
		return
	default:
		err = fmt.Errorf("invald config, invalid value, key:%s", key)
		return
	}
}

func getArrayFromMap(data map[string]interface{}, key string) (
	result []interface{}, err error) {

	value, ok := data[key]
	if !ok {
		err = fmt.Errorf("invalid config, not found:%s", key)
		return
	}

	switch inst := value.(type) {
	case []interface{}:
		result = inst
		return
	default:
		err = fmt.Errorf("invald config, invalid value, key:%s", key)
		return
	}
}

func getFloatFromMap(data map[string]interface{}, key string) (
	result float64, err error) {

	value, ok := data[key]
	if !ok {
		err = fmt.Errorf("invalid config, not found: %s", key)
		return
	}

	switch inst := value.(type) {
	case string:
		result, err = strconv.ParseFloat(inst, 64)
		return
	case int:
		result = float64(inst)
		return
	case int32:
		result = float64(inst)
		return
	case int64:
		result = float64(inst)
		return
	case float32:
		result = float64(inst)
	case float64:
		result = inst
	default:
		err = fmt.Errorf("read int failed, unknown type")
		return
	}

	return
}

func getIntFromMap(data map[string]interface{}, key string) (
	result int64, err error) {

	value, ok := data[key]
	if !ok {
		err = fmt.Errorf("invalid config, not found:%s", key)
		return
	}

	switch inst := value.(type) {
	case string:
		result, err = strconv.ParseInt(inst, 10, 64)
		return
	case int:
		result = int64(inst)
		return
	case int32:
		result = int64(inst)
		return
	case int64:
		result = inst
		return
	case float32:
		result = int64(inst)
	case float64:
		result = int64(inst)
	default:
		err = fmt.Errorf("read int failed, unknown type")
		return
	}

	return
}

func getStringFromMap(data map[string]interface{}, key string) (
	result string, err error) {

	value, ok := data[key]
	if !ok {
		err = fmt.Errorf("invalid config,not found %s", key)
		return
	}

	result, ok = value.(string)
	if !ok {
		err = fmt.Errorf("invalid config, config key:%s", key)
		return
	}

	return
}

func getBoolFromMap(data map[string]interface{}, key string) (
	result bool, err error) {

	value, ok := data[key]
	if !ok {
		err = fmt.Errorf("invalid config,not found %s", key)
		return
	}

	result, ok = value.(bool)
	if !ok {
		err = fmt.Errorf("invalid config, config key:%s", key)
		return
	}

	return
}
