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
	_, _ = returns2Values()
}

func rightTest(t *testing.T) {
	v, err := returns2Values()

	require.ErrorIs(t, nil, err)
	require.Equal(t, v, "")
}

func withSwitchExpr() {
	_, err := returns2Values()

	switch err {
	case nil:
		print(err)
	}
}

func withSwitch() {
	_, err := returns2Values()

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
	_, err := returns2Values()

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

func withOk() {
	_, ok := returnsBool()
	if ok {
		return
	}
}

func errWithOk() {
	var errCh = make(chan error)
	close(errCh)

	err, ok := <-errCh
	if !ok { // linter triggers on this
		return
	}

	panic(err.Error())
}
