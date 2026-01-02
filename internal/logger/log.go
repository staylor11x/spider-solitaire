package logger

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
)

var enabled atomic.Bool

func init() {
	enabled.Store(os.Getenv("SPIDER_DEBUG") == "1")
}

func SetEnabled(v bool) { enabled.Store(v) }
func Enabled() bool     { return enabled.Load() }

func logf(level string, format string, args ...any) {
	if !enabled.Load() {
		return
	}
	log.Printf("%s: %s", level, fmt.Sprintf(format, args...))
}

func Debug(format string, args ...any) { logf("DEBUG", format, args...) }
func Info(format string, args ...any)  { logf("INFO", format, args...) }
func Warn(format string, args ...any)  { logf("WARN", format, args...) }
func Error(format string, args ...any) { logf("ERROR", format, args...) }
