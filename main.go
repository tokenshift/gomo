package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New("gomo", "Command-line pomodoro time tracking.")

	breakCmd  = app.Command("break", "Start a break.")
	breakLong = breakCmd.Flag("long", "Start a long break.").Short('l').Bool()

	configCmd = app.Command("config", "Get/update config values.")
	configKey = configCmd.Arg("key", "Name of the config value to get/set.").String()
	configVal = configCmd.Arg("value", "The new value to set.").String()

	restartCmd = app.Command("restart", "Start a new work session and restart all counters.")

	statusCmd  = app.Command("status", "Display the current work status.")
	statusAuto = statusCmd.Flag("auto", "Automatically start the next break/session if time has expired.").Short('a').Bool()

	workCmd = app.Command("work", "Begin a work session.")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case breakCmd.FullCommand():
		if *breakLong {
			StartLongBreak()
		} else {
			StartBreak()
		}
	case configCmd.FullCommand():
		if configKey == nil {
			ShowAllConfig()
		} else if configVal == nil {
			ShowConfig(*configKey)
		} else {
			SetConfig(*configKey, *configVal)
		}
	case restartCmd.FullCommand():
		ResetStatus()
		StartWorkSession()
	case statusCmd.FullCommand():
		if *statusAuto {
			AutoAdvance()
		}
		DisplayStatus()
	case workCmd.FullCommand():
		StartWorkSession()
	}
}

func checkFatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}
