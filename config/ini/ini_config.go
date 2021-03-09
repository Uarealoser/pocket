package ini

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type IniConfig struct {
	Sections map[string]IniSection
}

func NewIniConfig() (ini *IniConfig) {
	return &IniConfig{
		Sections: map[string]IniSection{},
	}
}

/*
	从Reader获取IniConfig
*/
func NewIniConfigFromReader(reader io.Reader) (ini *IniConfig, err error) {
	p, err := parseIni(reader)
	if err != nil {
		return
	}
	sections := map[string]IniSection{}
	for k, v := range p.Sections {
		sections[k] = v
	}
	ini = &IniConfig{
		Sections: sections,
	}
	return
}

/*
	从文件获取IniConfig
*/
func NewIniConfigFromFile(filePath string) (ini *IniConfig, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	return NewIniConfigFromReader(f)
}

/*
	从String获取IniConfig
*/
func NewIniConfigFromText(content string) (ini *IniConfig, err error) {
	reader := strings.NewReader(content)
	return NewIniConfigFromReader(reader)
}

/*
	从字节数组获取IniConfig
*/
func NewIniConfigFromBytes(content []byte) (ini *IniConfig, err error) {
	reader := bytes.NewReader(content)
	return NewIniConfigFromReader(reader)
}

func (ini *IniConfig) GetSection(sectionName string) (section map[string]interface{}, err error) {
	section, ok := ini.Sections[sectionName]
	if !ok {
		err = fmt.Errorf("no such section %s", sectionName)
	}
	return
}

func (ini *IniConfig) GetIniSection(sectionName string) (section IniSection, err error) {
	section, ok := ini.Sections[sectionName]
	if !ok {
		err = fmt.Errorf("no such section %s", sectionName)
	}
	return
}

func (ini *IniConfig) GetSectionKeys(section string) (keys []string) {

	sec, err := ini.GetSection(section)
	if err != nil {
		return
	}

	for k, _ := range sec {
		keys = append(keys, k)
	}
	return
}

func (ini *IniConfig) GetValue(section, key string) (value interface{}, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetValue(key)
}

func (ini *IniConfig) GetString(section, key string) (value string, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetString(key)
}

func (ini *IniConfig) GetStringDefault(section, key string, defval string) (value string) {

	var err error
	value, err = ini.GetString(section, key)
	if err != nil {
		return defval
	}

	return
}

func (ini *IniConfig) GetInt(section, key string) (value int, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetInt(key)
}

func (ini *IniConfig) GetIntDefault(section, key string, defval int) (value int) {

	var err error
	value, err = ini.GetInt(section, key)
	if err != nil {
		return defval
	}

	return
}

func (ini *IniConfig) GetInt64(section, key string) (value int64, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetInt64(key)
}

func (ini *IniConfig) GetInt64Default(section, key string, defval int64) (value int64) {

	var err error
	value, err = ini.GetInt64(section, key)
	if err != nil {
		return defval
	}

	return
}

func (ini *IniConfig) GetSlice(section, key string) (value []string, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetSlice(key)
}

func (ini *IniConfig) GetSliceDefault(section, key string, defval []string) (value []string) {

	var err error
	value, err = ini.GetSlice(section, key)
	if err != nil {
		return defval
	}

	return
}

func (ini *IniConfig) GetMap(section, key string) (value map[string]string, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetMap(key)
}

func (ini *IniConfig) GetMapDefault(section, key string, defval map[string]string) (value map[string]string) {

	var err error
	value, err = ini.GetMap(section, key)
	if err != nil {
		return defval
	}

	return
}

func (ini *IniConfig) GetBool(section, key string) (value bool, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetBool(key)
}

func (ini *IniConfig) GetBoolDefault(section, key string, defval bool) (value bool) {

	var err error
	value, err = ini.GetBool(section, key)
	if err != nil {
		return defval
	}

	return
}

func (ini *IniConfig) GetFloat64(section, key string) (value float64, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetFloat64(key)
}

func (ini *IniConfig) GetFloat64Default(section, key string, defval float64) (value float64) {

	var err error
	value, err = ini.GetFloat64(section, key)
	if err != nil {
		return defval
	}

	return
}

func (ini *IniConfig) GetSections() (sections []string) {
	sections = make([]string, 0)
	for section := range ini.Sections {
		sections = append(sections, section)
	}
	return sections
}

func (ini *IniConfig) GetDuration(section, key string) (value time.Duration, err error) {

	sec, err := ini.GetIniSection(section)
	if err != nil {
		return
	}

	return sec.GetDuration(key)
}

func (ini *IniConfig) GetDurationDefault(section, key string, defval time.Duration) (value time.Duration) {

	var err error
	value, err = ini.GetDuration(section, key)
	if err != nil {
		return defval
	}

	return
}
