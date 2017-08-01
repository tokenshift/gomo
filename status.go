package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

type Status struct {
	SessionCount int
	State        State
	Started      time.Time
}

type State string

const (
	Work       State = "WORK"
	ShortBreak       = "SHORT_BREAK"
	LongBreak        = "LONG_BREAK"
)

const StatusPath = "~/.gomo/status"

func statusPath() string {
	dir, err := homedir.Expand(StatusPath)
	checkFatal(err)

	return dir
}

func DefaultStatus() Status {
	return Status{
		SessionCount: 0,
		State:        Work,
		Started:      time.Now(),
	}
}

func GetStatus() Status {
	var status Status

	_, err := toml.DecodeFile(statusPath(), &status)
	if err != nil {
		if os.IsNotExist(err) {
			status = DefaultStatus()
			SaveStatus(status)
		} else {
			checkFatal(err)
		}
	}

	return status
}

func SaveStatus(status Status) {
	err := os.MkdirAll(filepath.Dir(statusPath()), os.ModePerm)
	checkFatal(err)

	f, err := os.Create(statusPath())
	checkFatal(err)
	defer f.Close()

	encoder := toml.NewEncoder(f)
	err = encoder.Encode(status)
	checkFatal(err)
}

func (status Status) MinutesElapsed() float64 {
	return time.Now().Sub(status.Started).Minutes()
}

func (status Status) Ancient() bool {
	return status.MinutesElapsed() > 720
}
