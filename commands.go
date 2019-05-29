package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Commands struct {
	Config
	Log
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
	if cmd.State == None || cmd.Status.Ancient() {
		cmd.ResetStatus()
	} else {
		for cmd.expired() {
			cmd.autoAdvance1()
		}
	}
}

func (cmd *Commands) autoAdvance1() {
	if cmd.expired() {
		switch cmd.State {
		case Work:
			cmd.StartBreak(cmd.Started.Add(cmd.WorkSessionDuration()))
		case ShortBreak:
			cmd.StartWorkSession(cmd.Started.Add(cmd.ShortBreakDuration()))
		case LongBreak:
			cmd.StartWorkSession(cmd.Started.Add(cmd.LongBreakDuration()))
		case None:
			cmd.StartWorkSession(time.Now())
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
	case None:
		fmt.Printf("No active pomodoro\n")
	}
}

func (cmd *Commands) ResetStatus() {
	start := round(time.Now(), cmd.WorkSessionRound)

	cmd.SessionCount = 1
	cmd.State = Work
	cmd.Started = start

	SaveStatus(cmd.Status)
	cmd.AddLogEntry(start, "", cmd.State, "Reset")
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
	case "WorkSessionRound":
		cmd.WorkSessionRound = configInt(val)
	case "ShortBreakMinutes":
		cmd.ShortBreakMinutes = configInt(val)
	case "ShortBreakRound":
		cmd.ShortBreakRound = configInt(val)
	case "LongBreakMinutes":
		cmd.LongBreakMinutes = configInt(val)
	case "LongBreakRound":
		cmd.LongBreakRound = configInt(val)
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
	case "WorkSessionRound":
		fmt.Println(cmd.WorkSessionRound)
	case "ShortBreakMinutes":
		fmt.Println(cmd.ShortBreakMinutes)
	case "ShortBreakRound":
		fmt.Println(cmd.ShortBreakRound)
	case "LongBreakMinutes":
		fmt.Println(cmd.LongBreakMinutes)
	case "LongBreakRound":
		fmt.Println(cmd.LongBreakRound)
	case "LongBreakInterval":
		fmt.Println(cmd.LongBreakInterval)
	default:
		fmt.Fprintln(os.Stderr, "ERROR: Unrecognized config key:", key)
		os.Exit(1)
	}
}

func (cmd *Commands) StartBreak(t time.Time) {
	switch cmd.State {
	case Work:
		if cmd.SessionCount >= cmd.LongBreakInterval {
			cmd.StartLongBreak(t)
		} else {
			cmd.StartShortBreak(t)
		}
	case None:
		cmd.StartShortBreak(t)
	case ShortBreak, LongBreak:
		// Do nothing
	}
}

func (cmd *Commands) StartLongBreak(t time.Time) {
	t = round(t, cmd.LongBreakRound)

	switch cmd.State {
	case Work, ShortBreak, None:
		cmd.AddLogEntry(t, cmd.State, LongBreak, "")
		cmd.State = LongBreak
		cmd.Started = t
	case LongBreak:
		// Do nothing
	}

	SaveStatus(cmd.Status)
}

func (cmd *Commands) StartShortBreak(t time.Time) {
	t = round(t, cmd.ShortBreakRound)

	switch cmd.State {
	case Work, LongBreak, None:
		cmd.AddLogEntry(t, cmd.State, ShortBreak, "")
		cmd.State = ShortBreak
		cmd.Started = t
	case ShortBreak:
		// Do nothing
	}

	SaveStatus(cmd.Status)
}

func (cmd *Commands) StartWorkSession(t time.Time) {
	t = round(t, cmd.WorkSessionRound)

	switch cmd.State {
	case Work:
		// Do nothing
	case ShortBreak:
		cmd.AddLogEntry(t, cmd.State, Work, "")
		cmd.State = Work
		cmd.Started = t
		cmd.SessionCount += 1
	case LongBreak, None:
		cmd.AddLogEntry(t, cmd.State, Work, "")
		cmd.State = Work
		cmd.Started = t
		cmd.SessionCount = 1
	}

	SaveStatus(cmd.Status)
}
