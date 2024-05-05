package logger

import (
	"fmt"
	"log"
	"strings"

	"github.com/TwiN/go-color"
)

var logLevelFncMap = map[string]func(string){
	"debug": func(msg string) {
		log.Println(msg)
	},
	"error": func(msg string) {
		log.Fatalln(msg)
	},
	"info": func(msg string) {
		log.Println(msg)
	},
	"warn": func(msg string) {
		log.Println(msg)
	},
}

func printNonError(level string, msg string) {
	prefix := fmt.Sprintf("[%s]", strings.ToUpper(level))

	if level == "debug" {
		prefix = color.With(color.Gray, prefix)
	} else if level == "info" {
		prefix = color.With(color.Blue, prefix)
	} else if level == "warn" {
		prefix = color.With(color.Yellow, prefix)
	}

	logLevelFncMap[level](fmt.Sprintf("%s %s", prefix, msg))
}

func printError(err error) {
	log.Fatalln(color.With(color.Red, "[ERROR]"), err)
}

// prints a debug message to the console
func Debug(msg string) {
	printNonError("debug", msg)
}

// prints an error message to the console
// and exits the application
func Error(err error) {
	printError(err)
}

// prints an info message to the console
func Info(msg string) {
	printNonError("info", msg)
}

// prints a warning message to the console
func Warn(msg string) {
	printNonError("warn", msg)
}
