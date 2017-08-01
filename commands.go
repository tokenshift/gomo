package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Commands struct {
	Config
	Status
}

func (cmd Commands) minutesLeft() float64 {
	switch cmd.State {
	case Work:
		return float64(cmd.WorkSessionMinutes) - cmd.MinutesElapsed()
	case ShortBreak:
		return float64(cmd.ShortBreakMinutes) - cmd.MinutesElapsed()
	case LongBreak:
		return float64(cmd.LongBreakMinutes) - cmd.MinutesElapsed()
	}

	return 0.0
}

func (cmd Commands) expired() bool {
	return cmd.minutesLeft() <= 0.0
}

func (cmd *Commands) AutoAdvance() {
	if cmd.Status.Ancient() {
		cmd.ResetStatus()
	} else if cmd.expired() {
		switch cmd.State {
		case Work:
			cmd.StartBreak()
		case ShortBreak:
			cmd.StartWorkSession()
		case LongBreak:
			cmd.StartWorkSession()
		}
	}
}

func (cmd Commands) DisplayStatus() {
	minutesLeft := cmd.minutesLeft()

	if minutesLeft <= 0 {
		if cmd.State == Work {
			fmt.Println("Time to take a break!")
		} else {
			fmt.Println("Back to work!")
		}
	}

	switch cmd.State {
	case Work:
		fmt.Printf("Working (%.1f minutes remaining)\n", minutesLeft)
	case ShortBreak:
		fmt.Printf("Short break (%.1f minutes remaining)\n", minutesLeft)
	case LongBreak:
		fmt.Printf("Long break (%.1f minutes remaining)\n", minutesLeft)
	}
}

func (cmd *Commands) ResetStatus() {
	cmd.SessionCount = 0
	cmd.State = Work
	cmd.Started = time.Now()

	SaveStatus(cmd.Status)
}

func configInt(val string) int {
	i, err := strconv.ParseUint(val, 0, 0)
	checkFatal(err)

	return int(i)
}

func (cmd *Commands) SetConfig(key, val string) {
	switch key {
	case "WorkSessionMinutes":
		cmd.WorkSessionMinutes = configInt(val)
	case "ShortBreakMinutes":
		cmd.ShortBreakMinutes = configInt(val)
	case "LongBreakMinutes":
		cmd.LongBreakMinutes = configInt(val)
	case "LongBreakInterval":
		cmd.LongBreakInterval = configInt(val)
	default:
		fmt.Fprintln(os.Stderr, "ERROR: Unrecognized config key:", key)
		os.Exit(1)
	}

	SaveConfig(cmd.Config)
}

func (cmd Commands) ShowAllConfig() {
	cmd.Config.Write(os.Stdout)
}

func (cmd Commands) ShowConfig(key string) {
	switch key {
	case "WorkSessionMinutes":
		fmt.Println(cmd.WorkSessionMinutes)
	case "ShortBreakMinutes":
		fmt.Println(cmd.ShortBreakMinutes)
	case "LongBreakMinutes":
		fmt.Println(cmd.LongBreakMinutes)
	case "LongBreakInterval":
		fmt.Println(cmd.LongBreakInterval)
	default:
		fmt.Fprintln(os.Stderr, "ERROR: Unrecognized config key:", key)
		os.Exit(1)
	}
}

func (cmd *Commands) StartBreak() {
	switch cmd.State {
	case Work:
		if cmd.SessionCount >= cmd.LongBreakInterval {
			cmd.StartLongBreak()
		} else {
			cmd.StartShortBreak()
		}
	case ShortBreak, LongBreak:
		// Do nothing
	}
}

func (cmd *Commands) StartLongBreak() {
	switch cmd.State {
	case Work, ShortBreak:
		cmd.State = LongBreak
		cmd.Started = time.Now()
	case LongBreak:
		// Do nothing
	}

	SaveStatus(cmd.Status)
}

func (cmd *Commands) StartShortBreak() {
	switch cmd.State {
	case Work, LongBreak:
		cmd.State = ShortBreak
		cmd.Started = time.Now()
	case ShortBreak:
		// Do nothing
	}

	SaveStatus(cmd.Status)
}

func (cmd *Commands) StartWorkSession() {
	switch cmd.State {
	case Work:
		// Do nothing
	case ShortBreak:
		cmd.State = Work
		cmd.Started = time.Now()
		cmd.SessionCount += 1
	case LongBreak:
		cmd.State = Work
		cmd.Started = time.Now()
		cmd.SessionCount = 0
	}

	SaveStatus(cmd.Status)
}
