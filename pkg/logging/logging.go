package logging

import (
    "fmt"
    "github.com/sirupsen/logrus"
    "io"
    "os"
    "path"
    "runtime"
)

var e *logrus.Entry

type Logger struct {
    *logrus.Entry
}

func GetLogger() *Logger {
    return &Logger{e}
}

type hook struct {
    Writers   []io.Writer
    LogLevels []logrus.Level
}

func (h *hook) Fire(entry *logrus.Entry) error {
    line, err := entry.String()
    if err != nil {
        return err
    }

    for _, w := range h.Writers {
        if _, err := w.Write([]byte(line)); err != nil {
            return err
        }
    }

    return nil
}

func (h *hook) Levels() []logrus.Level {
    return h.LogLevels
}

func init() {
    l := logrus.New()
    l.SetReportCaller(true)
    l.Formatter = &logrus.TextFormatter{
        CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
            filename := path.Base(frame.File)
            return fmt.Sprintf("%s()", frame.Func.Name()),
                fmt.Sprintf("%s:%d", filename, frame.Line)
        },
        DisableColors: true,
        FullTimestamp: true,
        DisableQuote:  true,
    }

    err := os.MkdirAll("logs", 0644)
    if err != nil {
        panic(err)
    }

    logsFile, err := os.OpenFile("logs/logs.log",
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
    if err != nil {
        panic(err)
    }

    l.SetOutput(io.Discard)

    l.AddHook(&hook{[]io.Writer{logsFile, os.Stdout}, logrus.AllLevels})

    l.SetLevel(logrus.TraceLevel)

    e = logrus.NewEntry(l)
}
