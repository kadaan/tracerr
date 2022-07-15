// Package tracerr makes error output more informative.
// It adds stack trace to error and can display error with source fragments.
//
// Check example of output here https://github.com/kadaan/tracerr
package tracerr

import (
	"errors"
	"fmt"
	"runtime"
)

// DefaultFrameCapacity is a default capacity for frames array.
// It can be changed to number of expected frames
// for purpose of performance optimisation.
var DefaultFrameCapacity = 20

// DefaultFrameSkipCount is a number of frames to skip
// when retrieving the stack frames for the error.
var DefaultFrameSkipCount = 2

type Tracerr interface {
	CustomError(err error, frames []Frame) Error
	Errorf(message string, args ...interface{}) Error
	New(message string) Error
	Wrap(err error) Error
	Unwrap(err error) error
}

func NewTracerr(frameCapacity int, stackFrameSkipCount int) Tracerr {
	return &tracerr{
		frameCapacity:       frameCapacity,
		stackFrameSkipCount: stackFrameSkipCount,
	}
}

type tracerr struct {
	frameCapacity       int
	stackFrameSkipCount int
}

var Default = NewTracerr(DefaultFrameCapacity, DefaultFrameSkipCount)

func (t *tracerr) CustomError(err error, frames []Frame) Error {
	return &errorData{
		err:    err,
		frames: frames,
	}
}

func (t *tracerr) Errorf(message string, args ...interface{}) Error {
	return t.trace(fmt.Errorf(message, args...))
}

func (t *tracerr) New(message string) Error {
	return t.trace(fmt.Errorf(message))
}

func (t *tracerr) Wrap(err error) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	if wrapped := errors.Unwrap(err); wrapped != nil {
		e, ok := wrapped.(*errorData)
		err := fmt.Errorf("%w", Unwrap(err))
		if ok {
			return &errorData{
				err:    err,
				frames: e.frames,
			}
		}
	}
	return t.trace(err)
}

func (t *tracerr) Unwrap(err error) error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if !ok {
		return err
	}
	return e.Unwrap()
}

func (t *tracerr) trace(err error) Error {
	skip := t.stackFrameSkipCount
	frames := make([]Frame, 0, t.frameCapacity)
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := Frame{
			Func: fn.Name(),
			Line: line,
			Path: path,
		}
		frames = append(frames, frame)
		skip++
	}
	return &errorData{
		err:    err,
		frames: frames,
	}
}

// Error is an error with stack trace.
type Error interface {
	Error() string
	StackTrace() []Frame
	Unwrap() error
}

type errorData struct {
	// err contains original error.
	err error
	// frames contains stack trace of an error.
	frames []Frame
}

// CustomError creates an error with provided frames.
func CustomError(err error, frames []Frame) Error {
	return &errorData{
		err:    err,
		frames: frames,
	}
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
func Errorf(message string, args ...interface{}) Error {
	return Default.Errorf(message, args...)
}

// New creates new error with stacktrace.
func New(message string) Error {
	return Default.New(message)
}

// Wrap adds stacktrace to existing error.
func Wrap(err error) Error {
	return Default.Wrap(err)
}

// Unwrap returns the original error.
func Unwrap(err error) error {
	return Default.Unwrap(err)
}

// Error returns error message.
func (e *errorData) Error() string {
	return e.err.Error()
}

// StackTrace returns stack trace of an error.
func (e *errorData) StackTrace() []Frame {
	return e.frames
}

// Unwrap returns the original error.
func (e *errorData) Unwrap() error {
	return e.err
}

// Frame is a single step in stack trace.
type Frame struct {
	// Func contains a function name.
	Func string
	// Line contains a line number.
	Line int
	// Path contains a file path.
	Path string
}

// StackTrace returns stack trace of an error.
// It will be empty if err is not of type Error.
func StackTrace(err error) []Frame {
	e, ok := err.(Error)
	if !ok {
		return nil
	}
	return e.StackTrace()
}

// String formats Frame to string.
func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}
