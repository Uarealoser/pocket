package xlog

type Level uint32

const (
	DebugLevel Level = iota
	NoticeLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	panicLevel // 仅仅是为了异常情况，xlog未初始化成功等，包内使用，打印到stderr 留作以后扩展
)

const DefaultLevel = DebugLevel

var levelNames = [...]string{
	DebugLevel:  "DEBUG",
	NoticeLevel: "NOTICE",
	WarnLevel:   "WARN",
	ErrorLevel:  "ERROR",
	FatalLevel:  "FATAL",
	panicLevel:  "PANIC",
}

func (l Level) String() string {
	return levelNames[l]
}

var levelStrings = map[string]Level{
	"debug":  DebugLevel,
	"notice": NoticeLevel,
	"warn":   WarnLevel,
	"error":  ErrorLevel,
	"fatal":  FatalLevel,
	"panic":  panicLevel,
}

func StringToLevel(s string) Level {
	if lv, ok := levelStrings[s]; ok {
		return lv
	}
	return DefaultLevel
}

func (lv Level) IsEnabled(baseLevel Level) bool {
	return lv >= baseLevel
}
