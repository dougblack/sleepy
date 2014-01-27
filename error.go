package sleepy

// ErrorString is an Error that returns a string
// representation.
type ErrorString struct {
    s string
}

// Returns the string associated with this ErrorString.
func (e *ErrorString) Error() string {
    return e.s
}
