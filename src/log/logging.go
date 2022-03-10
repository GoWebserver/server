package log

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"server/src"
)

type logOptions struct {
	suffix    string
	colorCode string
}

var prefix = map[string]logOptions{
	"Debug":  {suffix: "*", colorCode: "\u001B[38;2;100;80;100m"},
	"Normal": {suffix: ">", colorCode: "\u001B[38;2;255;255;255m"},
	"Error":  {suffix: "!", colorCode: "\u001b[38;2;255;0;0m"},
}

func Err(err error, message ...interface{}) {
	log(prefix["Error"], 1, message...)
	if err != nil {
		log(prefix["Error"], 1, err.Error())
	}
	if src.GetConfig().Debug {
		debug.PrintStack()
	}
}

func Debug(message ...interface{}) {
	if src.GetConfig().Debug {
		log(prefix["Debug"], 1, message...)
	}
}

func Log(message ...interface{}) {
	log(prefix["Normal"], 1, message...)
}

type LogWriter struct {
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	log(prefix["Normal"], 2, string(p))
	return len(p), nil
}

func log(logOption logOptions, skip uint8, message ...interface{}) {
	now := time.Now()

	var location string

	if src.GetConfig().Logging.LogFile {
		_, file, line, ok := runtime.Caller(int(1 + skip))
		if !ok {
			file = "???"
			line = 0
		}

		file = file[strings.LastIndexByte(file, '/')+1:] // get relative path
		var locationStretch = strconv.Itoa(int(src.GetConfig().Logging.StretchFile))
		location = fmt.Sprintf("%-"+locationStretch+"s", fmt.Sprintf("%s:%d", file, line))
	}

	var printPrefix string
	if src.GetConfig().Logging.LogPrefix {
		var prefixStretch = strconv.Itoa(int(src.GetConfig().Logging.StretchPrefix))
		printPrefix = fmt.Sprintf("%-"+prefixStretch+"s %s", prefix, logOption.suffix)
	}

	var printStr string
	for _, mess := range message {
		printStr += fmt.Sprintf("%v", mess) + " "
	}

	_, err := os.Stdout.Write([]byte(fmt.Sprintf(
		"%s%s %s|%s %s \u001b[0m\n",
		logOption.colorCode,
		now.Format("2006.01.02 15:04:05.000"),
		location, printPrefix, printStr,
	)))
	if err != nil {
		return
	}
}
