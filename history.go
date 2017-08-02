package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

type History struct{}

const HistoryPath = "~/.gomo/history"

func historyPath() string {
	path, err := homedir.Expand(HistoryPath)
	checkFatal(err)

	return path
}

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func (h History) AddLogEntry(logEntry string) {
	err := os.MkdirAll(filepath.Dir(historyPath()), os.ModePerm)
	checkFatal(err)

	f, err := os.OpenFile(historyPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkFatal(err)
	defer f.Close()

	fmt.Fprintln(f, timestamp(), logEntry)
}

func (h History) DisplayHistory() {
	f, err := os.Open(historyPath())
	if err == nil {
		defer f.Close()
		io.Copy(os.Stdout, f)
	} else if !os.IsNotExist(err) {
		checkFatal(err)
	}
}
