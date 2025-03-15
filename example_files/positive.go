package src

import (
	"fmt"
	"log/syslog"
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
