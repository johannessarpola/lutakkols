package logger

import (
	"testing"
)

type testLog struct {
	format string
	args   []interface{}
}

func (l *testLog) Clear() {
	l.format = ""
	l.args = nil
}
func (l *testLog) Debug(v ...interface{}) {
	l.args = v
}

func (l *testLog) Debugf(format string, v ...interface{}) {
	l.format = format
	l.args = v
}

func (l *testLog) Infof(format string, v ...interface{}) {
	l.format = format
	l.args = v
}

func (l *testLog) Info(v ...interface{}) {
	l.args = v
}

func (l *testLog) Print(v ...interface{}) {
	l.args = v
}

func (l *testLog) Errorf(format string, v ...interface{}) {
	l.format = format
	l.args = v
}

func (l *testLog) Warnf(format string, v ...interface{}) {
	l.format = format
	l.args = v
}

func (l *testLog) Warn(v ...interface{}) {
	l.args = v
}

func (l *testLog) Error(v ...interface{}) {
	l.args = v
}

func TestLog(t *testing.T) {

	tl := testLog{}
	SetLogger(&tl)

	Log.Infof("format %s", "arg")
	if tl.format != "format %s" {
		t.Errorf("format not set")
	}
	if tl.args == nil {
		t.Errorf("args not set")
	}

	tl.Clear()

	Log.Info("arg2")
	if tl.args == nil {
		t.Errorf("args not set")
	}

}
