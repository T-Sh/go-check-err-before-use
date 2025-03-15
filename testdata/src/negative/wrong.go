package negative

import (
	"fmt"
	"log/syslog"
)

func wrong() {
	logger, err := syslog.New(syslog.Priority(1), "custom")

	logger.Write([]byte("bytes message"))

	if err != nil {
		fmt.Printf("got error: %v", err)

		return
	}
}
