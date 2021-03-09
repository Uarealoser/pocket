package config

import "time"

type Config interface {
	GetValue(section, key string) (value interface{}, err error)
	GetString(section, key string) (value string, err error)
	GetStringDefault(section, key string, defval string) (value string)
	GetInt(section, key string) (value int, err error)
	GetIntDefault(section, key string, defval int) (value int)
	GetInt64(section, key string) (value int64, err error)
	GetInt64Default(section, key string, defval int64) (value int64)
	GetSlice(section, key string) (value []string, err error)
	GetSliceDefault(section, key string, defval []string) (value []string)
	GetMap(section, key string) (value map[string]string, err error)
	GetMapDefault(section, key string, defval map[string]string) (value map[string]string)
	GetBool(section, key string) (value bool, err error)
	GetBoolDefault(section, key string, defval bool) (value bool)
	GetFloat64(section, key string) (value float64, err error)
	GetFloat64Default(section, key string, defval float64) (value float64)
	GetDuration(section, key string) (value time.Duration, err error)
	GetDurationDefault(section, key string, defval time.Duration) (value time.Duration)
	GetSectionKeys(section string) (keys []string)
	GetSection(section string) (sec map[string]interface{}, err error)
	GetSections() (sections []string)
}
