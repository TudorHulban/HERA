package pghera

/*
File contains package constants.
*/

import (
	"errors"
)

// ErrorNotAPointer Error message for data not being a pointer.
// Pointer needed for ....
var ErrorNotAPointer = errors.New("passed data is not a pointer")
