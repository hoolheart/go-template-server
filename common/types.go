package common

import "fmt"

// ModuleError represents an error from specific module and package
type ModuleError struct {
	module 		string//module name
	pack 			string//package name
	function	string//function name
	body			string//error body
}

// Error string of given ModuleError
func (err ModuleError) Error() (string) {
	return fmt.Sprintf("[%s/%s:%s] %s",err.module,err.pack,err.function,err.body)
}

// NewError creates a new ModuleError
func NewError(mod,pack,foo,body string) (err error) {
	return ModuleError{mod,pack,foo,body}
}
