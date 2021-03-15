package xlog

import "testing"

func TestStringToLevel(t *testing.T) {
	cases := []map[string]Level{
		{
			"debug":  DebugLevel,
			"notice": NoticeLevel,
			"warn":   WarnLevel,
			"error":  ErrorLevel,
			"fatal":  FatalLevel,
			"panic":  panicLevel,
		},
	}

	for _, cs := range cases {
		for params, lv := range cs {
			res := StringToLevel(params)
			if res != lv {
				t.Errorf("%s expected:%s,get:%s", params, lv, res)
			}
			// t.Logf("check level %s", res)
		}
	}
}
