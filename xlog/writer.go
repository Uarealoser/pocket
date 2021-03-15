package xlog

type XlogWriter interface {
	Name() string
	Level() Level
	Output(entry) error
}
