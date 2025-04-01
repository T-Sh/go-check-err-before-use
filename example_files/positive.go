package src

import (
	"errors"
	"fmt"
	"log/syslog"
	"testing"

	"github.com/stretchr/testify/require"
)

func right() {
	logger, err := syslog.New(syslog.Priority(1), "custom")
	if err != nil {
		fmt.Printf("got error: %v", err)

		return
	}

	logger.Write([]byte("bytes message"))
}

func withMultipleErrChecks() {
	logger, err := syslog.New(syslog.Priority(1), "custom")
	if err != nil && true {
		fmt.Printf("got error: %v", err)

		return
	}

	logger.Write([]byte("bytes message"))
}

func skippedErr() {
	logger, _ := syslog.New(syslog.Priority(1), "custom")
	logger.Write([]byte("bytes message"))
}

func noValues() {
	_, _ = with2Values()
}

func rightTest(t *testing.T) {
	v, err := with2Values()

	require.ErrorIs(t, nil, err)
	require.Equal(t, v, "")
}

func withSwitchExpr() {
	_, err := with2Values()

	switch err {
	case nil:
		print(err)
	}
}

func withSwitch() {
	_, err := with2Values()

	switch {
	case errors.Is(err, nil):
		print()
	case err != nil:
		print()
	default:
		print()
	}
}

type ErrStruct struct {
	field1 int
	field2 error
}

func errInStruct() ErrStruct {
	_, err := with2Values()

	e := ErrStruct{field1: 0, field2: err}

	return e
}

func SingleErrReturn() error {
	return nil
}

func singleErr() {
	err := SingleErrReturn()

	print()

	if err != nil {

	}
}
