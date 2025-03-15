package positive

import (
	"log/syslog"
)

func main() {
	logger, _ := syslog.New(syslog.Priority(1), "custom")
	logger.Write([]byte("bytes message"))
}
