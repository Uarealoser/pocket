package ini

import "errors"

type Unmarshaler interface {
	Unmarshal(value string) error
}

type Sectioner interface {
	Init() error
}

// Errors
var (
	// ErrorNilPtr _
	ErrorNilPtr = errors.New("ptr is nil")
	// ErrorDst _
	ErrorDst = errors.New("dst must be ptr whitch point to a struct")
	// ErrorUnmarshaler _
	ErrorUnmarshaler = errors.New("struct must have Unmarshal method")
	// ErrorSection _
	ErrorSection = errors.New("section should be a struct")
)
