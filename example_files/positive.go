package src

import (
	"errors"
	"fmt"
	"log/syslog"
	"os"
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

func rightWithVar() {
	var logger, err = syslog.New(syslog.Priority(1), "custom")
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

func errInStruct() ErrStruct {
	_, err := returns2Values()

	e := ErrStruct{field1: 0, field2: err}

	return e
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

func errWithReturn() error {
	_, err := returns2Values()

	return err
}

func withSwitchInsideCase(state int) (int, error) {
	_, err := os.OpenFile("test", 0, os.ModePerm)

	switch state {
	case 3:
		return -1, err
	}

	if err != nil {
		return -2, err
	}
	return 0, nil
}

func withIfInside(state int) (int, error) {
	_, err := os.OpenFile("test", 0, os.ModePerm)

	if state == 0 {
		print()
		return state, err
	}

	if err != nil {
		return -2, err
	}
	return 0, nil
}
