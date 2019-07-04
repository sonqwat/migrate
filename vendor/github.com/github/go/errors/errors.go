// Package errors provides simple error handling primitives.
//
// The traditional error handling idiom in Go is roughly akin to
//
//     if err != nil {
//             return err
//     }
//
// which applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error. For example
//
//     _, err := ioutil.ReadAll(r)
//     if err != nil {
//             return errors.Wrap(err, "read failed")
//     }
//
// Retrieving the cause of an error
//
// Using errors.Wrap constructs a stack of errors, adding context to the
// preceding error. Depending on the nature of the error it may be necessary
// to reverse the operation of errors.Wrap to retrieve the original error
// for inspection. Any error value which implements this interface
//
//     type Causer interface {
//             Cause() error
//     }
//
// can be inspected by errors.Cause. errors.Cause will recursively retrieve
// the topmost error which does not implement causer, which is assumed to be
// the original cause. For example:
//
//     switch err := errors.Cause(err).(type) {
//     case *MyError:
//             // handle specifically
//     default:
//             // unknown error
//     }
//
// Formatted printing of errors
//
// All error values returned from this package implement fmt.Formatter and can
// be formatted by the fmt package. The following verbs are supported
//
//     %s    print the error. If the error has a Cause it will be
//           printed recursively
//     %v    see %s
//     %+v   extended format. Each Frame of the error's StackTrace will
//           be printed in detail.
//
// Retrieving the stack trace of an error or wrapper
//
// New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
// invoked. This information can be retrieved with the following interface.
//
//     type stackTracer interface {
//             StackTrace() errors.StackTrace
//     }
//
// Where errors.StackTrace is defined as
//
//     type StackTrace []Frame
//
// The Frame type represents a call site in the stack trace. Frame supports
// the fmt.Formatter interface that can be used for printing information about
// the stack trace of this error. For example:
//
//     if err, ok := err.(stackTracer); ok {
//             for _, f := range err.StackTrace() {
//                     fmt.Printf("%+s:%d", f)
//             }
//     }
//
// See the documentation for Frame.Format for more details.
// Deprecated: Use github.com/pkg/errors instead.
package errors

import (
	"fmt"
	"io"
)

// _error is an error implementation returned by New and Errorf
// that implements its own fmt.Formatter.
type _error struct {
	msg string
	*stack
}

// Deprecated: Use github.com/pkg/errors.Error instead.
func (e _error) Error() string { return e.msg }

// Deprecated: Use github.com/pkg/errors.Format instead.
func (e _error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msg)
			fmt.Fprintf(s, "%+v", e.StackTrace())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	}
}

// New returns an error with the supplied message.
// Deprecated: Use github.com/pkg/errors.New instead.
func New(message string) error {
	return _error{
		message,
		callers(0),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Deprecated: Use github.com/pkg/errors.Errorf instead.
func Errorf(format string, args ...interface{}) error {
	return _error{
		fmt.Sprintf(format, args...),
		callers(0),
	}
}

// Panic formats the output of a `recover` call, be it either an error, a
// string message, or an arbitrary value. It creates a wrapped error with a
// suitable stack trace pointing at the source of the panic.
// Deprecated: Use github.com/pkg/errors.Errorf("panic: %v", err) instead.
func Panic(panic interface{}) error {
	stack := callers(3)
	switch p := panic.(type) {
	case error:
		return _error{p.Error(), stack}
	case string:
		return _error{p, stack}
	default:
		return _error{fmt.Sprintf("panic: %+v", panic), stack}
	}
}

type cause struct {
	cause error
	msg   string
}

// Deprecated: Use github.com/pkg/errors instead.
func (c cause) Error() string { return fmt.Sprintf("%s: %v", c.msg, c.Cause()) }
// Deprecated: Use github.com/pkg/errors instead.
func (c cause) Cause() error  { return c.cause }

// wrapper is an error implementation returned by Wrap and Wrapf
// that implements its own fmt.Formatter.
type wrapper struct {
	cause
	*stack
}

// Deprecated: Use github.com/pkg/errors instead.
func (w wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			fmt.Fprintf(s, "%+v", w.StackTrace())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// Wrap returns an error annotating err with message.
// If err is nil, Wrap returns nil.
// Deprecated: Use github.com/pkg/errors.WithMessage instead.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   message,
		},
		stack: callers(0),
	}
}

// Wrapf returns an error annotating err with the format specifier.
// If err is nil, Wrapf returns nil.
// Deprecated: Use github.com/pkg/errors.WithMessagef instead.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		stack: callers(0),
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type Causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
// Deprecated: Use github.com/pkg/errors.Cause instead.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}