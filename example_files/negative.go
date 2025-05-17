package src

import (
	"fmt"
	"log/syslog"
	"testing"

	"github.com/stretchr/testify/require"
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

func errName() {
	res, errVal := returns2Values()

	print(res)

	if errVal != nil {
	}
}

func wrongCustom2Func() {
	res, err := returns2Values()

	print(res)

	if err != nil {
	}
}

func wrongCustom3Func() {
	res, str, err := returns3Values()

	print(res, str)

	if err != nil {
	}
}

func wrongWithOk() {
	_, ok := returnsBool()
	print()

	if !ok {
		return
	}
}

func wrongTest(t *testing.T) {
	v, err := returns2Values()

	require.Equal(t, v, "")
	require.ErrorIs(t, nil, err)
}
