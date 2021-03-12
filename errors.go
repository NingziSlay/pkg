package pkg

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrorArrayOutOfRange = errors.New("array out of range")
	ErrorUnsupportedType = func(t string) error { return fmt.Errorf("unsupported type: %s", t) }
)

type ConfigError struct {
	msg string
	err error
}

func (c ConfigError) Error() string {
	if c.err != nil {
		return errors.Wrap(c.err, c.msg).Error()
	}
	return c.msg
}

func newE(msg string, args ...interface{}) ConfigError {
	return ConfigError{msg: fmt.Sprintf(msg, args...)}
}

func wrapE(msg string, err error) ConfigError {
	return ConfigError{msg: msg, err: err}
}

var (
	ErrorNonPointer = newE("parse of non-pointer")
	ErrorNilInput   = newE("parse of nil")
	ErrorNonStruct  = newE("parse of non-struct ptr")
)
