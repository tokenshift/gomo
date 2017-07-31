package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func AutoAdvance() {
	config := GetConfig()
	status := GetStatus()

	if status.Ancient() {
		ResetStatus()
	} else if status.Expired(config) {
		switch status.State {
		case Work:
			StartBreak()
		case ShortBreak:
			StartWorkSession()
		case LongBreak:
			StartWorkSession()
		}
	}
}

func DisplayStatus() {
	config := GetConfig()
	status := GetStatus()

	minutesLeft := status.MinutesLeft(config)

	if minutesLeft <= 0 {
		if status.State == Work {
			fmt.Println("Time to take a break!")
		} else {
			fmt.Println("Back to work!")
		}
	}

	switch status.State {
	case Work:
		fmt.Printf("Working (%.1f minutes remaining)\n", status.MinutesLeft(config))
	case ShortBreak:
		fmt.Printf("Short break (%.1f minutes remaining)\n", status.MinutesLeft(config))
	case LongBreak:
		fmt.Printf("Long break (%.1f minutes remaining)\n", status.MinutesLeft(config))
	}
}

func ResetStatus() {
	status := DefaultStatus()

	status.SessionCount = 0
	status.State = Work
	status.Started = time.Now()

	SaveStatus(status)
}

func configInt(val string) int {
	i, err := strconv.ParseUint(val, 0, 0)
	checkFatal(err)

	return int(i)
}

func SetConfig(key, val string) {
	config := GetConfig()

	switch key {
	case "WorkSessionMinutes":
		config.WorkSessionMinutes = configInt(val)
	case "ShortBreakMinutes":
		config.ShortBreakMinutes = configInt(val)
	case "LongBreakMinutes":
		config.LongBreakMinutes = configInt(val)
	case "LongBreakInterval":
		config.LongBreakInterval = configInt(val)
	default:
		fmt.Fprintln(os.Stderr, "ERROR: Unrecognized config key:", key)
		os.Exit(1)
	}

	SaveConfig(config)
}

func ShowAllConfig() {
	config := GetConfig()
	fmt.Println(config)
}

func ShowConfig(key string) {
	config := GetConfig()

	switch key {
	case "WorkSessionMinutes":
		fmt.Println(config.WorkSessionMinutes)
	case "ShortBreakMinutes":
		fmt.Println(config.ShortBreakMinutes)
	case "LongBreakMinutes":
		fmt.Println(config.LongBreakMinutes)
	case "LongBreakInterval":
		fmt.Println(config.LongBreakInterval)
	default:
		fmt.Fprintln(os.Stderr, "ERROR: Unrecognized config key:", key)
		os.Exit(1)
	}
}

func StartBreak() {
	status := GetStatus()
	config := GetConfig()

	switch status.State {
	case Work:
		if status.SessionCount >= config.LongBreakInterval {
			StartLongBreak()
		} else {
			StartShortBreak()
		}
	case ShortBreak, LongBreak:
		// Do nothing
	}
}

func StartLongBreak() {
	status := GetStatus()

	switch status.State {
	case Work, ShortBreak:
		status.State = LongBreak
		status.Started = time.Now()
	case LongBreak:
		// Do nothing
	}

	SaveStatus(status)
}

func StartShortBreak() {
	status := GetStatus()

	switch status.State {
	case Work, LongBreak:
		status.State = ShortBreak
		status.Started = time.Now()
	case ShortBreak:
		// Do nothing
	}

	SaveStatus(status)
}

func StartWorkSession() {
	status := GetStatus()

	switch status.State {
	case Work:
		// Do nothing
	case ShortBreak:
		status.State = Work
		status.Started = time.Now()
		status.SessionCount += 1
	case LongBreak:
		status.State = Work
		status.SessionCount = 0
	}

	SaveStatus(status)
}
