package config

import "reflect"

// 对比secion
func compareSections(oldCfg, newCfg Config, sections []string) (ok bool) {

	if oldCfg == nil || newCfg == nil {
		return oldCfg == newCfg
	}

	for _, section := range sections {
		oldSection, _ := oldCfg.GetSection(section)
		newSection, _ := newCfg.GetSection(section)
		if !reflect.DeepEqual(oldSection, newSection) {
			return
		}
	}

	ok = true
	return
}
