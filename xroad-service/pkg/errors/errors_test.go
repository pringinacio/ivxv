package errors

import (
	"errors"
	"fmt"
	"testing"
)

type NestedError struct {
	Err error
}

func (n *NestedError) Error() string {
	return n.Err.Error()
}

func methodThatAlwaysReturnsErrNotFound() error   { return errors.New("NOT_FOUND") }
func methodThatAlwaysReturnsErrBadRequest() error { return errors.New("BAD_REQUEST") }
func methodThatAlwaysReturnsErrVotingEnd() error  { return errors.New("VOTING_END") }

func TestEHSErrorCanBeUsedInErrorsIsFunc(t *testing.T) {
	errNotFound := methodThatAlwaysReturnsErrNotFound()
	errBadRequest := methodThatAlwaysReturnsErrBadRequest()
	errVotingEnd := methodThatAlwaysReturnsErrVotingEnd()

	wrapErrNotFound := &EHSError{Err: errNotFound}
	wrapErrBadRequest := &EHSError{Err: errBadRequest}
	wrapErrVotingEnd := &EHSError{Err: errVotingEnd}

	data := []error{wrapErrNotFound, wrapErrBadRequest, wrapErrVotingEnd}
	expected := []error{ErrNotFound, ErrBadRequest, ErrVotingEnd}

	for i, err := range data {
		areEqual := errors.Is(err, expected[i])
		if areEqual != true {
			msg := "Expected err to be equal to %v, but got %v"
			t.Fatal(fmt.Sprintf(msg, expected[i], err))
		}
	}
}

func TestNestedEHSErrorCanBeUsedInErrorsIsFunc(t *testing.T) {
	errNotFound := methodThatAlwaysReturnsErrNotFound()
	errBadRequest := methodThatAlwaysReturnsErrBadRequest()
	errVotingEnd := methodThatAlwaysReturnsErrVotingEnd()

	threeLevelsNestedErrNotFound := &NestedError{
		Err: &NestedError{
			Err: &NestedError{
				Err: errNotFound,
			},
		},
	}
	threeLevelsNestedErrBadRequest := &NestedError{
		Err: &NestedError{
			Err: &NestedError{
				Err: errBadRequest,
			},
		},
	}
	threeLevelsNestedErrVotingEnd := &NestedError{
		Err: &NestedError{
			Err: &NestedError{
				Err: errVotingEnd,
			},
		},
	}

	wrapThreeLevelsNestedErrNotFound := &EHSError{Err: threeLevelsNestedErrNotFound}
	wrapThreeLevelsNestedErrBadRequest := &EHSError{Err: threeLevelsNestedErrBadRequest}
	wrapThreeLevelsNestedErrVotingEnd := &EHSError{Err: threeLevelsNestedErrVotingEnd}

	data := []error{
		wrapThreeLevelsNestedErrNotFound,
		wrapThreeLevelsNestedErrBadRequest,
		wrapThreeLevelsNestedErrVotingEnd,
	}

	expected := []error{ErrNotFound, ErrBadRequest, ErrVotingEnd}

	for i, err := range data {
		areEqual := errors.Is(err, expected[i])
		if areEqual != true {
			msg := "Expected err to be equal to %v, but got %v"
			t.Fatal(fmt.Sprintf(msg, expected[i], err))
		}
	}
}

func TestRegularErrorCannotBeUsedInErrorsIsFunc(t *testing.T) {
	errNotFound := methodThatAlwaysReturnsErrNotFound()
	errBadRequest := methodThatAlwaysReturnsErrBadRequest()
	errVotingEnd := methodThatAlwaysReturnsErrVotingEnd()

	data := []error{errNotFound, errBadRequest, errVotingEnd}
	expected := []error{ErrNotFound, ErrBadRequest, ErrVotingEnd}

	for i, err := range data {
		areEqual := errors.Is(err, expected[i])
		if areEqual == true {
			msg := "Expected err not to be equal to %v"
			t.Fatal(fmt.Sprintf(msg, err))
		}
	}
}

func TestNestedRegularErrorCannotBeUsedInErrorsIsFunc(t *testing.T) {
	errNotFound := methodThatAlwaysReturnsErrNotFound()
	errBadRequest := methodThatAlwaysReturnsErrBadRequest()
	errVotingEnd := methodThatAlwaysReturnsErrVotingEnd()

	threeLevelsNestedErrNotFound := &NestedError{
		Err: &NestedError{
			Err: &NestedError{
				Err: errNotFound,
			},
		},
	}
	threeLevelsNestedErrBadRequest := &NestedError{
		Err: &NestedError{
			Err: &NestedError{
				Err: errBadRequest,
			},
		},
	}
	threeLevelsNestedErrVotingEnd := &NestedError{
		Err: &NestedError{
			Err: &NestedError{
				Err: errVotingEnd,
			},
		},
	}

	data := []error{
		threeLevelsNestedErrNotFound,
		threeLevelsNestedErrBadRequest,
		threeLevelsNestedErrVotingEnd,
	}

	expected := []error{ErrNotFound, ErrBadRequest, ErrVotingEnd}

	for i, err := range data {
		areEqual := errors.Is(err, expected[i])
		if areEqual == true {
			msg := "Expected err not to be equal to %v"
			t.Fatal(fmt.Sprintf(msg, err))
		}
	}
}
