package api

import (
	"errors"
	"fmt"
	"strings"
)

// Code classifies an Error. Each code has its own exit code.
type Code string

// The four error codes.
const (
	CodeError    Code = "error"     // anything else; exit 1
	CodeUsage    Code = "usage"     // a bad request; exit 2
	CodeConflict Code = "conflict"  // the request clashes with the log; exit 3
	CodeNotFound Code = "not_found" // a ref names nothing; exit 4
)

// Error is an operation failure a caller can act on. Next, if set, is a
// command that would help, such as `rhiza change new "<intent>"`.
type Error struct {
	Code    Code
	Message string
	Next    string
}

// Error returns the message.
func (e *Error) Error() string {
	return e.Message
}

// WithNext returns a copy of e with Next set, leaving e unchanged:
//
//	return Usagef("no change set").WithNext(`rhiza change new "<intent>"`)
func (e *Error) WithNext(next string) *Error {
	c := *e
	c.Next = next
	return &c
}

func newError(code Code, format string, args []any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Errorf returns an Error with code error.
func Errorf(format string, args ...any) *Error { return newError(CodeError, format, args) }

// Usagef returns an Error with code usage.
func Usagef(format string, args ...any) *Error { return newError(CodeUsage, format, args) }

// Conflictf returns an Error with code conflict.
func Conflictf(format string, args ...any) *Error { return newError(CodeConflict, format, args) }

// NotFoundf returns an Error with code not_found.
func NotFoundf(format string, args ...any) *Error { return newError(CodeNotFound, format, args) }

// AsError returns the first *Error in err's chain, or, if there is none, an
// Error with code error and err's message. It returns nil for nil.
func AsError(err error) *Error {
	if err == nil {
		return nil
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{Code: CodeError, Message: err.Error()}
}

// ExitCode maps err to the process exit code: 0 for nil, 2 for usage, 3 for
// conflict, 4 for not_found and 1 for anything else, including errors that
// wrap no *Error.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	switch AsError(err).Code {
	case CodeUsage:
		return 2
	case CodeConflict:
		return 3
	case CodeNotFound:
		return 4
	default:
		return 1
	}
}

// ErrorBody is the JSON form of an Error.
type ErrorBody struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Next    string `json:"next"`
}

// ErrorResponse is what --json prints for a failure:
// {"error":{"code","message","next"}}.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorJSON returns the ErrorResponse for err (see AsError). Its message is
// the found *Error's, not the whole wrapped chain's.
func ErrorJSON(err error) ErrorResponse {
	e := AsError(err)
	if e == nil {
		e = &Error{Code: CodeError}
	}
	return ErrorResponse{Error: ErrorBody{Code: e.Code, Message: e.Message, Next: e.Next}}
}

// article is the indefinite article for s: "an" before a vowel, else "a".
func article(s string) string {
	if s != "" && strings.ContainsRune("aeiouAEIOU", rune(s[0])) {
		return "an"
	}
	return "a"
}
