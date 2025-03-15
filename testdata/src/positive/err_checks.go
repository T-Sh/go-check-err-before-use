package positive

import (
	"errors"
	"fmt"
	"log/syslog"
)

func withMultipleChecksInIf() {
	logger, err := syslog.New(syslog.Priority(1), "custom")
	if err != nil && true {
		fmt.Printf("got error: %v", err)

		return
	}

	logger.Write([]byte("bytes message"))
}

func checkErrInSecondExprInIf() {
	logger, err := syslog.New(syslog.Priority(1), "custom")
	if true && err != nil {
		fmt.Printf("got error: %v", err)

		return
	}

	logger.Write([]byte("bytes message"))
}

func checkErrInErrPkg() {
	logger, err := syslog.New(syslog.Priority(1), "custom")
	if errors.Is(err, errors.ErrUnsupported) {
		fmt.Printf("got error: %v", err)

		return
	}

	logger.Write([]byte("bytes message"))
}

func checkErrInFunc() {
	f := func(e error) bool {
		return true
	}
	logger, err := syslog.New(syslog.Priority(1), "custom")
	if f(err) {
		fmt.Printf("got error: %v", err)

		return
	}

	logger.Write([]byte("bytes message"))
}

func withInlineCheck() {
	if logger, err := syslog.New(syslog.Priority(1), "custom"); err != nil {
		logger.Write([]byte("bytes message"))
	}
}
