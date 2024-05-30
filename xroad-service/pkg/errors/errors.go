package errors

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotFound   = errors.New("NOT_FOUND")
	ErrBadRequest = errors.New("BAD_REQUEST")
	ErrVotingEnd  = errors.New("VOTING_END")
)

// EHSError is an any error returned by the EHS backend.
// It implements Error and Is methods so could be safely used in
// errors.Is function as well as being used as a regular error interface.
type EHSError struct {
	Err error
}

// Error overrides error interface's Error method.
func (v EHSError) Error() string {
	return v.Err.Error()
}

// Is overrides anonymous interface{ Is(error) bool } interface's method,
// used in errors.Is.
func (v EHSError) Is(err error) bool {
	return v.Err.Error() == err.Error()
}

type FieldError struct {
	Code  string      `json:"code"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (e FieldError) ToErr() error {
	marshal, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return errors.New(string(marshal))
}
