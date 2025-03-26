package errs

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Wrap - adding operation name(function name) to error.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	// get program counter, line number.
	pc, _, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("unknown : %w", err)
	}

	functionName := runtime.FuncForPC(pc).Name()  // get fn name.
	functionName = trimFunctionName(functionName) // trim fn name.

	// format fn details for wrap error.
	functionDetails := fmt.Sprintf("%s:%d", functionName, line)

	return fmt.Errorf("[%s] %w", functionDetails, err)
}

// Wrapf - adding operation name(function name) to error with custom message.
func Wrapf(err error, message string) error {
	if err == nil {
		return nil
	}

	// get program counter, line number.
	pc, _, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("unknown : %w", err)
	}

	functionName := runtime.FuncForPC(pc).Name()  // get fn name.
	functionName = trimFunctionName(functionName) // trim fn name.

	// format fn details for wrap error.
	functionDetails := fmt.Sprintf("%s:%d", functionName, line)

	if message == "" {
		return fmt.Errorf("[%s] : %w", functionDetails, err)
	}

	return fmt.Errorf("[%s] %s : %w", functionDetails, message, err)
}

// Unwrap - for unwrap error to raw error.
func Unwrap(err error) error {
	var lastErr = err

	for unwrappedErr := errors.Unwrap(err); unwrappedErr != nil; unwrappedErr = errors.Unwrap(unwrappedErr) {
		lastErr = unwrappedErr
	}

	return lastErr
}

// trimFunctionName - for trim function name.
// Clear the path, leaving only the package and function name.
func trimFunctionName(fullName string) string {
	lastSlashIndex := strings.LastIndex(fullName, "/")
	if lastSlashIndex != -1 {
		fullName = fullName[lastSlashIndex+1:]
	}

	return fullName
}

var (
	ErrNoKeys    = errors.New("key not found")
	ErrKeyExists = errors.New("key already exists")
)
