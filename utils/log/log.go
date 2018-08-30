package log

import (
	"os"

	"github.com/op/go-logging"
)

func init() {
	consoleBackend := logging.NewLogBackend(os.Stdout, "", 0)
	consoleFormat := "%{color}[%{module}][%{level:.5s}] %{time:15:04:05.000} %{shortfile} %{message} %{color:reset}"
	consoleFormatter := logging.MustStringFormatter(consoleFormat)
	consoleFormatterBackend := logging.NewBackendFormatter(consoleBackend, consoleFormatter)
	logging.SetBackend(consoleFormatterBackend)
}

func GetLogger(module string) *logging.Logger {
	return logging.MustGetLogger(module)
}
