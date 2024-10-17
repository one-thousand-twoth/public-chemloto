package models

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type EngineError struct {
	// Op is the operation being performed, usually the name of the method
	// being invoked.
	Op Op
	// User is the name of the user attempting the operation.
	User UserName
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	// Param represents the parameter related to the error.
	Param Parameter
	// Code is a human-readable, short representation of the error
	Code Code
	// The underlying error that triggered this one, if any.
	Err error
}

func (e *EngineError) isZero() bool {
	return e.User == "" && e.Kind == 0 && e.Param == "" && e.Code == "" && e.Err == nil
}

// Unwrap method allows for unwrapping errors using errors.As
func (e *EngineError) Unwrap() error {
	return e.Err
}

func (e *EngineError) Error() string {
	return e.Err.Error()
}

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// UserName is a string representing a user
type UserName string

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind uint8

// Parameter represents the parameter related to the error.
type Parameter string

// Code is a human-readable, short representation of the error
type Code string

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
	// Unauthenticated is used when a request lacks valid authentication credentials.
	//
	// For Unauthenticated errors, the response body will be empty.
	// The error is logged and http.StatusUnauthorized (401) is sent.
	Unauthenticated // Unauthenticated Request
	// Unauthorized is used when a user is authenticated, but is not authorized
	// to access the resource.
	//
	// For Unauthorized errors, the response body should be empty.
	// The error is logged and http.StatusForbidden (403) is sent.
	Unauthorized
	UnsupportedMediaType // Unsupported Media Type
)
