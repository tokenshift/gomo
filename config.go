package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	WorkSessionMinutes int
	ShortBreakMinutes  int
	LongBreakMinutes   int
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
		ShortBreakMinutes:  5,
		LongBreakMinutes:   15,
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
