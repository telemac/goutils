package stacktrace

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// StackTrace represents a stack trace with file, line, and function name information.
type StackTrace struct {
	File     string // The file name of the frame.
	Line     int    // The line number of the frame.
	Function string // The name of the function for the frame.
}

func (st StackTrace) String() string {
	// if st.File contains github.com, remove left part of string
	file := st.File
	githubIdx := strings.Index(file, "github.com")
	if githubIdx > 0 {
		file = st.File[githubIdx:]
	}
	return fmt.Sprintf("File: %s:%d, Function: %s", file, st.Line, st.Function)
}

// StackTraces represents a slice of StackTrace values.
type StackTraces []StackTrace

// String formats each stack trace in the slice into a string, joined by newline characters.
func (st StackTraces) String() string {
	var result []string
	for _, s := range st {
		result = append(result, s.String())
	}
	return strings.Join(result, "\n")
}

// GetCaller returns a StackTrace value representing the file, line, and function name for the caller of the function that calls GetCaller.
func GetCaller(skip int) StackTrace {
	// Get the file and line number of the caller's caller.
	_, file, line, _ := runtime.Caller(skip)

	// Parse the function name from the file path.
	function := path.Base(file)

	// Return a new StackTrace struct.
	return StackTrace{
		File:     file,
		Line:     line,
		Function: function,
	}
}

// GetStackTrace returns a slice of StackTrace values representing the file, line, and function name for the stack trace.
func GetStackTrace() StackTraces {
	// Create a slice to hold the StackTrace values.
	var stackTraces StackTraces

	// Create a slice to hold the program counters for each caller.
	var callers [1024]uintptr

	// Get the number of callers in the stack trace and their program counters.
	numCallers := runtime.Callers(2, callers[:])

	// Loop through the callers and append the file, line, and function name to the slice.
	for i := 0; i < numCallers; i++ {
		// Get the function for the caller's program counter.
		funcInfo := runtime.FuncForPC(callers[i])

		// Get the file and line number of the caller.
		file, line := funcInfo.FileLine(callers[i])

		// Get the name of the function.
		funcName := funcInfo.Name()

		// Create a StackTrace value for the frame.
		frame := StackTrace{
			File:     file,
			Line:     line,
			Function: funcName,
		}

		// Append the StackTrace value to the slice.
		stackTraces = append(stackTraces, frame)
	}

	return stackTraces
}
