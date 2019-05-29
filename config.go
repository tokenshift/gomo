package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	WorkSessionMinutes int
	WorkSessionRound   int
	ShortBreakMinutes  int
	ShortBreakRound    int
	LongBreakMinutes   int
	LongBreakRound     int
	LongBreakInterval  int
}

const ConfigPath = "~/.gomo/config"

func configPath() string {
	dir, err := homedir.Expand(ConfigPath)
	checkFatal(err)

	return dir
}

func DefaultConfig() Config {
	return Config{
		WorkSessionMinutes: 25,
		WorkSessionRound:   5,
		ShortBreakMinutes:  5,
		ShortBreakRound:    1,
		LongBreakMinutes:   15,
		LongBreakRound:     1,
		LongBreakInterval:  4,
	}
}

func GetConfig() Config {
	var config Config

	_, err := toml.DecodeFile(configPath(), &config)
	if err != nil {
		if os.IsNotExist(err) {
			config = DefaultConfig()
			SaveConfig(config)
		} else {
			checkFatal(err)
		}
	}

	return config
}

func SaveConfig(config Config) {
	err := os.MkdirAll(filepath.Dir(configPath()), os.ModePerm)
	checkFatal(err)

	f, err := os.Create(configPath())
	checkFatal(err)
	defer f.Close()

	encoder := toml.NewEncoder(f)
	err = encoder.Encode(config)
	checkFatal(err)
}

func (config Config) Write(out io.Writer) {
	encoder := toml.NewEncoder(out)
	encoder.Encode(config)
}

func toDuration(minutes int) time.Duration {
	dur, err := time.ParseDuration(fmt.Sprintf("%dm", minutes))
	if err != nil {
		panic(err)
	}

	return dur
}

func (config Config) WorkSessionDuration() time.Duration {
	return toDuration(config.WorkSessionMinutes)
}

func (config Config) LongBreakDuration() time.Duration {
	return toDuration(config.LongBreakMinutes)
}

func (config Config) ShortBreakDuration() time.Duration {
	return toDuration(config.ShortBreakMinutes)
}
