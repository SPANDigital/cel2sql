package cel2sql

import (
	"errors"
	"fmt"
)

// Sentinel errors for common error conditions.
// These exported errors allow callers to use errors.Is() for specific error handling.
var (
	// ErrUnsupportedExpression indicates an unsupported CEL expression type
	ErrUnsupportedExpression = errors.New("unsupported expression type")

	// ErrInvalidFieldName indicates an invalid or empty field name
	ErrInvalidFieldName = errors.New("invalid field name")

	// ErrInvalidSchema indicates a problem with the provided schema
	ErrInvalidSchema = errors.New("invalid schema")

	// ErrInvalidRegexPattern indicates an invalid regex pattern
	ErrInvalidRegexPattern = errors.New("invalid regex pattern")

	// ErrMaxDepthExceeded indicates recursion depth limit exceeded
	ErrMaxDepthExceeded = errors.New("maximum recursion depth exceeded")

	// ErrMaxOutputLengthExceeded indicates output length limit exceeded
	ErrMaxOutputLengthExceeded = errors.New("maximum output length exceeded")

	// ErrInvalidComprehension indicates an invalid comprehension expression
	ErrInvalidComprehension = errors.New("invalid comprehension expression")

	// ErrMaxComprehensionDepthExceeded indicates comprehension nesting depth exceeded
	ErrMaxComprehensionDepthExceeded = errors.New("maximum comprehension depth exceeded")

	// ErrInvalidArguments indicates invalid function arguments
	ErrInvalidArguments = errors.New("invalid function arguments")

	// ErrInvalidTimestampOperation indicates an invalid timestamp operation
	ErrInvalidTimestampOperation = errors.New("invalid timestamp operation")

	// ErrInvalidDuration indicates an invalid duration value
	ErrInvalidDuration = errors.New("invalid duration value")

	// ErrInvalidJSONPath indicates an invalid JSON path expression
	ErrInvalidJSONPath = errors.New("invalid JSON path")

	// ErrInvalidOperator indicates an invalid operator
	ErrInvalidOperator = errors.New("invalid operator")

	// ErrUnsupportedType indicates an unsupported type
	ErrUnsupportedType = errors.New("unsupported type")

	// ErrContextCanceled indicates the operation was cancelled via context
	ErrContextCanceled = errors.New("operation cancelled")

	// ErrInvalidByteArrayLength indicates byte array exceeds maximum length
	ErrInvalidByteArrayLength = errors.New("byte array exceeds maximum length")
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
// Always uses errMsgConversionFailed as the user message to prevent information leakage,
// unless the error is already a formatted error with a specific message that should be preserved
func wrapConversionError(err error, internalContext string) *ConversionError {
	if err == nil {
		return &ConversionError{
			UserMessage:     errMsgConversionFailed,
			InternalDetails: internalContext,
		}
	}

	// Build internal details with context if provided
	var internalDetails string
	if internalContext != "" {
		internalDetails = fmt.Sprintf("%s: %v", internalContext, err)
	} else {
		internalDetails = err.Error()
	}

	// Check if this is a plain error (not ConversionError) - preserve its message
	if _, isConversionError := err.(*ConversionError); !isConversionError {
		return &ConversionError{
			UserMessage:     err.Error(), // Preserve the original error message
			InternalDetails: internalDetails,
			WrappedErr:      err,
		}
	}

	// If it's a ConversionError, check if it has a specific (non-generic) user message
	// that should be preserved through the wrapping chain
	convErr := err.(*ConversionError)
	if convErr.UserMessage != errMsgConversionFailed {
		// This is a specific error message that should be preserved
		return &ConversionError{
			UserMessage:     convErr.UserMessage,
			InternalDetails: internalDetails,
			WrappedErr:      err,
		}
	}

	// For generic ConversionErrors, use generic message
	return &ConversionError{
		UserMessage:     errMsgConversionFailed,
		InternalDetails: internalDetails,
		WrappedErr:      err,
	}
}

// Common error messages (sanitized for end users)
// These use sentence case (capitalize first word) for consistency
const (
	errMsgUnsupportedExpression      = "Unsupported expression type"
	errMsgInvalidOperator            = "Invalid operator in expression"
	errMsgUnsupportedType            = "Unsupported type in expression"
	errMsgUnsupportedComprehension   = "Unsupported comprehension operation"
	errMsgComprehensionDepthExceeded = "Comprehension nesting exceeds maximum depth"
	errMsgInvalidFieldAccess         = "Invalid field access in expression"
	errMsgConversionFailed           = "Failed to convert expression component"
	errMsgInvalidTimestampOp         = "Invalid timestamp operation"
	errMsgInvalidDuration            = "Invalid duration value"
	errMsgInvalidArguments           = "Invalid function arguments"
	errMsgUnknownType                = "Unknown type in schema"
	errMsgUnknownEnum                = "Unknown enum value"
	errMsgInvalidPattern             = "Invalid pattern in expression"
)
