// Package logger provides a logger for Centrifugo server.
// This is an adapted code from Steve Francia's jWalterWeatherman
// library - see https://github.com/spf13/jWalterWeatherman
package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Level describes the chosen log level
type Level int

type NotePad struct {
	Handle io.Writer
	Level  Level
	Prefix string
	Logger *log.Logger
}

// checkLevel exists to prevent calling underlying logger methods when not needed.
func (n *NotePad) checkLevel() bool {
	if n.Level < outputThreshold && n.Level < logThreshold {
		return false
	}
	return true
}

func (n *NotePad) Print(v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Print(v...)
}

func (n *NotePad) Printf(format string, v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Printf(format, v...)
}

func (n *NotePad) Println(v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Println(v...)
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (n *NotePad) Fatal(v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Fatal(v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (n *NotePad) Fatalf(format string, v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Fatalf(format, v...)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (n *NotePad) Fatalln(v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Fatalln(v...)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (n *NotePad) Panic(v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Panic(v...)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (n *NotePad) Panicf(format string, v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Panicf(format, v...)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (n *NotePad) Panicln(v ...interface{}) {
	if ok := n.checkLevel(); !ok {
		return
	}
	n.Logger.Panicln(v...)
}

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelFatal
	LevelNone

	DefaultLogThreshold    = LevelInfo
	DefaultStdoutThreshold = LevelInfo
)

var (
	logger *log.Logger

	LogHandle  io.Writer = ioutil.Discard
	OutHandle  io.Writer = os.Stdout
	BothHandle io.Writer = io.MultiWriter(LogHandle, OutHandle)

	Flag int = log.Ldate | log.Ltime

	TRACE    *NotePad = &NotePad{Level: LevelTrace, Handle: os.Stdout, Logger: logger, Prefix: "[T]: "}
	DEBUG    *NotePad = &NotePad{Level: LevelDebug, Handle: os.Stdout, Logger: logger, Prefix: "[D]: "}
	INFO     *NotePad = &NotePad{Level: LevelInfo, Handle: os.Stdout, Logger: logger, Prefix: "[I]: "}
	WARN     *NotePad = &NotePad{Level: LevelWarn, Handle: os.Stdout, Logger: logger, Prefix: "[W]: "}
	ERROR    *NotePad = &NotePad{Level: LevelError, Handle: os.Stdout, Logger: logger, Prefix: "[E]: "}
	CRITICAL *NotePad = &NotePad{Level: LevelCritical, Handle: os.Stdout, Logger: logger, Prefix: "[C]: "}
	FATAL    *NotePad = &NotePad{Level: LevelFatal, Handle: os.Stdout, Logger: logger, Prefix: "[F]: "}

	NotePads []*NotePad = []*NotePad{TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL, FATAL}

	logThreshold    Level = DefaultLogThreshold
	outputThreshold Level = DefaultStdoutThreshold
)

var LevelMatches = map[string]Level{
	"TRACE":    LevelTrace,
	"DEBUG":    LevelDebug,
	"INFO":     LevelInfo,
	"WARN":     LevelWarn,
	"ERROR":    LevelError,
	"CRITICAL": LevelCritical,
	"FATAL":    LevelFatal,
	"NONE":     LevelNone,
}

func init() {
	initialize()
}

// initialize initializes loggers
func initialize() {
	BothHandle = io.MultiWriter(LogHandle, OutHandle)
	for _, n := range NotePads {
		if n.Level < outputThreshold && n.Level < logThreshold {
			n.Handle = ioutil.Discard
		} else if n.Level >= outputThreshold && n.Level >= logThreshold {
			n.Handle = BothHandle
		} else if n.Level >= outputThreshold && n.Level < logThreshold {
			n.Handle = OutHandle
		} else {
			n.Handle = LogHandle
		}
	}

	for _, n := range NotePads {
		n.Logger = log.New(n.Handle, n.Prefix, Flag)
	}
}

// Ensures that the level provided is within the bounds of available levels
func levelCheck(level Level) Level {
	switch {
	case level <= LevelTrace:
		return LevelTrace
	case level >= LevelFatal:
		return LevelFatal
	default:
		return level
	}
}

// Establishes a threshold where anything matching or above will be logged
func SetLogThreshold(level Level) {
	logThreshold = levelCheck(level)
	initialize()
}

// Establishes a threshold where anything matching or above will be output
func SetStdoutThreshold(level Level) {
	outputThreshold = levelCheck(level)
	initialize()
}

// Conveniently Sets the Log Handle to a io.writer created for the file behind the given filepath
// Will only append to this file
func SetLogFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		CRITICAL.Println("Failed to open log file:", path, err)
		return err
	}
	LogHandle = file
	initialize()
	return nil
}

func SetLogFlag(flag int) {
	Flag = flag
	initialize()
}
