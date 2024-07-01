package logger

import "fmt"

type StdOutLogger struct{}

func (l *StdOutLogger) Debug(v ...interface{}) {
	fmt.Print(v...)
}

func (l *StdOutLogger) Debugf(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func (l *StdOutLogger) Infof(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func (l *StdOutLogger) Info(v ...interface{}) {
	fmt.Println(v...)
}

func (l *StdOutLogger) Errorf(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func (l *StdOutLogger) Error(v ...interface{}) {
	fmt.Println(v...)
}

func (l *StdOutLogger) Warnf(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func (l *StdOutLogger) Warn(v ...interface{}) {
	fmt.Println(v...)
}
