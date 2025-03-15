package src

import (
	"fmt"
	"log/syslog"
)

func insideFor() {
	for range 3 {
		logger, err := syslog.New(syslog.Priority(1), "custom")

		logger.Write([]byte("bytes message"))

		if err != nil {
			fmt.Printf("got error: %v", err)

			return
		}
	}
}

func insideIf() {
	if true {
		logger, err := syslog.New(syslog.Priority(1), "custom")

		logger.Write([]byte("bytes message"))

		if err != nil {
			fmt.Printf("got error: %v", err)

			return
		}
	}
}

func wrong() {
	logger, err := syslog.New(syslog.Priority(1), "custom")

	logger.Write([]byte("bytes message"))

	if err != nil {
		fmt.Printf("got error: %v", err)

		return
	}
}

func noValues() {
	_, _ = with2Values()
}

func errName() {
	res, errVal := with2Values()

	print(res)

	if errVal != nil {
	}
}

func wrongCustom2Func() {
	res, err := with2Values()

	print(res)

	if err != nil {
	}
}

func wrongCustom3Func() {
	res, str, err := with3Values()

	print(res, str)

	if err != nil {
	}
}

func with2Values() (int, error) {
	return 0, nil
}

func with3Values() (int, string, error) {
	return 0, "", nil
}
