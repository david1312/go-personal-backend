package crashy

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	MaxStackLength = 50
)

type (
	ErrCode string
	Error   struct {
		Code       ErrCode `json:"code"`
		Message    string  `json:"message"`
		stacktrace string
		origin     error
	}
	CodeMapper interface {
		ErrCode() ErrCode
		ErrMessage() string
	}
	Wrapped interface {
		Unwrap() error
	}
	StackTracer interface {
		StackTrace() string
	}
)

// New create a pointer to a new instance of `Error`
func New(err error, code ErrCode, message string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		stacktrace: getStackTrace(2),
		origin:     err,
	}
}

func (e Error) Error() string {
	if e.Code == "" {
		e.Code = ErrCodeUnexpected
	}
	if e.Message == "" {
		e.Message = Message(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implementation of `Wrapped` interface
func (e Error) Unwrap() error {
	return e.origin
}

// StackTrace implementation of `StackTracer` interface
func (e Error) StackTrace() string {
	return e.stacktrace
}
func (e Error) ErrCode() ErrCode {
	if e.Code == "" {
		e.Code = ErrCodeUnexpected
	}
	return e.Code
}
func (e Error) ErrMessage() string {
	if e.Code == "" {
		e.Code = ErrCodeUnexpected
	}
	if e.Message == "" {
		e.Message = Message(e.Code)
	}
	return e.Message
}

func Is(err error, code ErrCode) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(CodeMapper); ok {
		return e.ErrCode() == code
	}
	return false
}

// Wrap error but keep origin and stacktrace
func Wrap(err error, code ErrCode, message string) error {
	if err == nil {
		return New(err, code, message)
	}

	if e, ok := err.(Wrapped); ok { //we do not support multi wrapper
		err = e.Unwrap()
	}

	var st string
	if e, ok := err.(StackTracer); ok { //we keep original stacktrace
		st = e.StackTrace()
	}
	if st == "" {
		st = getStackTrace(2)
	}

	return &Error{
		Code:       code,
		Message:    message,
		origin:     err,
		stacktrace: st,
	}
}

func Wrapc(err error, code ErrCode) error {
	return Wrap(err, code, "")
}

func Wrapf(err error, code ErrCode, format string, params ...interface{}) error {
	return Wrap(err, code, fmt.Sprintf(format, params...))
}

func getStackTrace(skip int) string {
	stackBuf := make([]uintptr, MaxStackLength)
	length := runtime.Callers(skip, stackBuf[:])
	stack := stackBuf[:length]

	trace := ""
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n\tFile: %s, Line: %d. Function: %s", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}
	return trace
}
