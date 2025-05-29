// Package provides error implementation for engine
package enerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"sort"

	"github.com/invopop/validation"
)

// MultiError contains a map indicating the parameter and its error.
// Serves to display errors in user interface forms
type MultiError struct {
	Errs map[string]error
}

func (e *MultiError) Error() string {
	return fmt.Sprintf("%+v", e.Errs)
}

// MarshalJSON converts the Errors into a valid JSON.
func (es *MultiError) MarshalJSON() ([]byte, error) {
	errs := map[string]interface{}{}
	for key, err := range es.Errs {
		if ms, ok := err.(json.Marshaler); ok {
			errs[key] = ms
		} else {
			errs[key] = err.Error()
		}
	}
	return json.Marshal(errs)
}

func mapStringToError(arg map[string]string) map[string]error {
	res := make(map[string]error, len(arg))
	for k, v := range arg {
		res[k] = errors.New(v)
	}
	return res
}

func errsToMap(errs []string) map[string]error {
	result := make(map[string]error, len(errs)/2+1)
	for i := 0; i < len(errs); i += 2 {
		key := errs[i]
		if i+1 < len(errs) {
			result[key] = errors.New(errs[i+1])
		} else {
			result[key] = errors.New("%not-defined%")
		}
	}
	return result
}

// EM() extends E() and builds an error value from its arguments with different types of input args:
//
// There must be at least one argument or EM panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//
//	map[string]error,
//	map[string]string,
//		The underlying multi-errors
func EM(args ...interface{}) *ApplicationError {

	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	var e = &MultiError{}
	var errPairs []string
	newArgs := make([]any, 0)
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			if errPairs == nil {
				errPairs = make([]string, 0, 2)
			}
			errPairs = append(errPairs, arg)
		case validation.Errors:
			e = &MultiError{arg}
		case map[string]error:
			e = &MultiError{arg}
		case map[string]string:
			e = &MultiError{mapStringToError(arg)}
		case Op, Kind, UserName:
			newArgs = append(newArgs, arg)
		default:
			_, file, line, _ := runtime.Caller(1)
			return &ApplicationError{Err: fmt.Errorf(
				"errors.E: bad call from %s:%d: %v, unknown type %T, value %v in error call",
				file, line, args, arg, arg)}
		}
	}
	if len(errPairs) > 0 {
		e = &MultiError{
			Errs: errsToMap(errPairs),
		}
	}

	newArgs = append(newArgs, e)

	return E(newArgs...)
}

// ApplicationError is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type ApplicationError struct {
	// Op is the operation being performed, usually the name of the method
	// being invoked.
	Op Op
	// User is the name of the user attempting the operation.
	User UserName
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	// The underlying error that triggered this one, if any.
	Err error
}

func (e *ApplicationError) isZero() bool {
	return e.User == "" && e.Kind == 0 && e.Err == nil
}

// Unwrap method allows for unwrapping errors using errors.As
func (e *ApplicationError) Unwrap() error {
	return e.Err
}

func (e *ApplicationError) Error() string {
	return e.Err.Error()
}

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// Return slog.Attr with "op" key
func OpAttr(op Op) slog.Attr {
	return slog.Attr{
		Key:   "op",
		Value: slog.StringValue(string(op)),
	}
}

// UserName is a string representing a user
type UserName string

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind uint8

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers. Do not reorder this list or remove
// any items since that will change their values.
// New items must be added only to the end.
const (
	Other          Kind = iota // Unclassified error. This value is not printed in the error message.
	Invalid                    // Invalid operation for this type of item.
	IO                         // External I/O error such as network failure.
	Exist                      // Item already exists.
	NotExist                   // Item does not exist.
	Private                    // Information withheld.
	Internal                   // Internal error or inconsistency.
	BrokenLink                 // Link target does not exist.
	Database                   // Error from database.
	Validation                 // Input validation error.
	Unanticipated              // Unanticipated error.
	InvalidRequest             // Invalid Request

	Unauthenticated // Unauthenticated Request
	// Unauthenticated is used when the user is not registered in the engine.
	Unidentified

	// Unauthorized is used when a user is authenticated, but is not authorized
	// to access the resource or perform an action.
	Unauthorized
	UnsupportedMediaType // Unsupported Media Type
	NotExistAction       // Action does not exist
	MaxPlayers           // Engine reached its maximum players
	AlreadyStarted       // Engine started so some functions blocked
	GameLogic
)

// E builds an error value from its arguments.
// There must be at least one argument or E panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//
//	UserName
//		The username of the user attempting the operation.
//	string
//		Treated as an error message and assigned to the
//		Err field after a call to errors.New.
//	enerr.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
func E(args ...interface{}) *ApplicationError {

	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &ApplicationError{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case UserName:
			e.User = arg
		case Kind:
			e.Kind = arg
		case string:
			e.Err = Str(arg)
		case *ApplicationError:
			// Make a copy
			errorCopy := *arg
			e.Err = &errorCopy
		case error:
			e.Err = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			return &ApplicationError{Err: fmt.Errorf(
				"errors.E: bad call from %s:%d: %v, unknown type %T, value %v in error call",
				file, line, args, arg, arg)}
		}
	}

	prev, ok := e.Err.(*ApplicationError)
	if !ok {
		return e
	}

	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}

	return e
}

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case IO:
		return "I/O error"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case BrokenLink:
		return "link target does not exist"
	case Private:
		return "information withheld"
	case Internal:
		return "internal error"
	case Database:
		return "database error"
	case Validation:
		return "input validation error"
	case Unanticipated:
		return "unanticipated error"
	case InvalidRequest:
		return "invalid request error"
	case Unauthenticated:
		return "unauthenticated request"
	case Unauthorized:
		return "unauthorized request"
	case UnsupportedMediaType:
		return "unsupported media type"
	default:
		return "unknown error kind"
	}
}

// Str returns an error that formats as the given text. It is intended to
// be used as the error-typed argument to the E function.
func Str(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// KindIs reports whether err is an *Error of the given Kind.
// If err is nil then KindIs returns false.
func KindIs(kind Kind, err error) bool {
	var e *ApplicationError
	if errors.As(err, &e) {
		if e.Kind != Other {
			return e.Kind == kind
		}
		if e.Err != nil {
			return KindIs(kind, e.Err)
		}
	}
	return false
}

// OpStack returns the op stack information for an error
func OpStack(err error) []string {
	type o struct {
		Op    string
		Order int
	}

	e := err
	i := 0
	var os []o

	// loop through all wrapped errors and add to struct
	// order will be from top to bottom of stack
	for errors.Unwrap(e) != nil {
		var errsError *ApplicationError
		if errors.As(e, &errsError) {
			if errsError.Op != "" {
				op := o{Op: string(errsError.Op), Order: i}
				os = append(os, op)
			}
		}
		e = errors.Unwrap(e)
		i++
	}

	// reverse the order of the stack (bottom to top)
	sort.Slice(os, func(i, j int) bool { return os[i].Order > os[j].Order })

	// pull out just the stack info, now in reversed order
	var ops []string
	for _, op := range os {
		ops = append(ops, op.Op)
	}

	return ops
}
