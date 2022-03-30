package common

import (
	"github.com/antiWalker/golib/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var InfoLogger *zap.SugaredLogger
var DebugLogger *zap.SugaredLogger
var ErrorLogger *zap.SugaredLogger
var WarnLogger *zap.SugaredLogger

func init() {
	ip, _ := ExternalIP()
	hostname, _ := os.Hostname()
	for _, arg := range os.Args {
		if strings.Contains(arg, "logDir") {
			logDir := strings.Split(arg, "=")[1]
			*global.LogsDir = logDir
			break
		}
	}
	InfoLogger = NewLogger(*global.LogsDir+"/info.log", zapcore.InfoLevel, 500, 30, 7, true, hostname, ip.String()).Sugar()
	DebugLogger = NewLogger(*global.LogsDir+"/debug.log", zapcore.InfoLevel, 500, 30, 7, true, hostname, ip.String()).Sugar()
	ErrorLogger = NewLogger(*global.LogsDir+"/error.log", zapcore.InfoLevel, 500, 30, 7, true, hostname, ip.String()).Sugar()
	WarnLogger = NewLogger(*global.LogsDir+"/warn.log", zapcore.InfoLevel, 500, 30, 7, true, hostname, ip.String()).Sugar()
}
