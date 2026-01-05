package utils

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02T15:04:05-07:00")
	level := strings.ToUpper(entry.Level.String())

	status := "-"
	if v, ok := entry.Data["status"]; ok {
		status = fmt.Sprintf("%v", v)
	}

	latency := "-"
	if v, ok := entry.Data["latency"]; ok {
		latency = fmt.Sprintf("%v", v)
	}

	ip := "-"
	if v, ok := entry.Data["ip"]; ok {
		ip = fmt.Sprintf("%v", v)
	}

	method := "-"
	if v, ok := entry.Data["method"]; ok {
		method = fmt.Sprintf("%v", v)
	}

	path := "-"
	if v, ok := entry.Data["path"]; ok {
		path = fmt.Sprintf("%v", v)
	}

	msg := fmt.Sprintf(
		"%s | %-5s | %-6v | %-8v | %-15v | %-6v | %v\n",
		timestamp,
		level,
		status,
		latency,
		ip,
		method,
		path,
	)

	return []byte(msg), nil
}
