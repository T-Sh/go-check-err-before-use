package positive

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
