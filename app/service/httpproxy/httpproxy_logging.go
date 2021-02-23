package httpproxy

import (
	"bytes"
	"github.com/gogf/gf/os/glog"
)

// errorLogger is the error logging logger for underlying net/http.Server.
type errorLogger struct {
	logger *glog.Logger
}

// Write implements the io.Writer interface.
func (l *errorLogger) Write(p []byte) (n int, err error) {
	l.logger.Skip(1).Error(string(bytes.TrimRight(p, "\r\n")))
	return len(p), nil
}
