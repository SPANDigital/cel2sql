package cel2sql

import (
	"fmt"
)

// ConversionError represents an error that occurred during CEL to SQL conversion.
// It provides a sanitized user-facing message while preserving detailed information
// for logging and debugging. This prevents information disclosure through error messages
// (CWE-209: Information Exposure Through Error Message).
type ConversionError struct {
	// UserMessage is the sanitized error message safe to display to end users
	UserMessage string

	// InternalDetails contains detailed information for logging and debugging
	// This should NEVER be exposed to end users
	InternalDetails string

	// WrappedErr is the underlying error, if any
	WrappedErr error
}

// Error returns the user-facing error message.
// This is what gets displayed when the error is returned to callers.
func (e *ConversionError) Error() string {
	return e.UserMessage
}

// Unwrap returns the wrapped error for use with errors.Is and errors.As
func (e *ConversionError) Unwrap() error {
	return e.WrappedErr
}

// Internal returns the full internal details for logging purposes.
// This should only be used with structured logging, never displayed to users.
func (e *ConversionError) Internal() string {
	if e.InternalDetails != "" {
		return e.InternalDetails
	}
	return e.UserMessage
}

// newConversionError creates a ConversionError with separate user and internal messages
func newConversionError(userMsg string, internalDetails string) *ConversionError {
	return &ConversionError{
		UserMessage:     userMsg,
		InternalDetails: internalDetails,
	}
}

// newConversionErrorf creates a ConversionError with formatted internal details
func newConversionErrorf(userMsg string, internalFormat string, args ...interface{}) *ConversionError {
	return &ConversionError{
		UserMessage:     userMsg,
		InternalDetails: fmt.Sprintf(internalFormat, args...),
	}
}

// wrapConversionError wraps an existing error with a generic user-facing message
// Always uses errMsgConversionFailed as the user message to prevent information leakage
func wrapConversionError(err error, internalContext string) *ConversionError {
	internalDetails := internalContext
	if err != nil {
		if internalContext != "" {
			internalDetails = fmt.Sprintf("%s: %v", internalContext, err)
		} else {
			internalDetails = err.Error()
		}
	}

	return &ConversionError{
		UserMessage:     errMsgConversionFailed,
		InternalDetails: internalDetails,
		WrappedErr:      err,
	}
}

// Common error messages (sanitized for end users)
const (
	errMsgUnsupportedExpression    = "unsupported expression type"
	errMsgInvalidOperator          = "invalid operator in expression"
	errMsgUnsupportedType          = "unsupported type in expression"
	errMsgUnsupportedComprehension = "unsupported comprehension operation"
	errMsgInvalidFieldAccess       = "invalid field access in expression"
	errMsgConversionFailed         = "failed to convert expression component"
	errMsgInvalidTimestampOp       = "invalid timestamp operation"
	errMsgInvalidDuration          = "invalid duration value"
	errMsgInvalidArguments         = "invalid function arguments"
	errMsgUnknownType              = "unknown type in schema"
	errMsgUnknownEnum              = "unknown enum value"
	errMsgInvalidPattern           = "invalid pattern in expression"
)
