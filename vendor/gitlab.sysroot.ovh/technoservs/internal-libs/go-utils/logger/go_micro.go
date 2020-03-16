package logger

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type GoMicroLogrus struct {
	l logrus.FieldLogger
}

func parentCaller() string {
	pc, _, _, ok := runtime.Caller(4)
	fn := runtime.FuncForPC(pc)
	if ok && fn != nil {
		return fn.Name()
	}

	return ""
}

func (gml GoMicroLogrus) Log(v ...interface{}) {
	pc := parentCaller()
	if strings.HasSuffix(pc, "Fatal") {
		gml.l.Fatal(v...)
	} else {
		gml.l.Info(v...)
	}
}

func (gml GoMicroLogrus) Logf(format string, v ...interface{}) {
	pc := parentCaller()
	if strings.HasSuffix(pc, "Fatalf") {
		gml.l.Fatalf(format, v...)
	} else {
		gml.l.Infof(format, v...)
	}
}
