package main

import (
	"database/sql"
	"fmt"
	"io"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
)

type Log struct{}

const LogPath = "~/.gomo/history.sqlite3"

func (l Log) AddLogEntry(oldState, newState State, message string) {
	db := getLogDb()

	stmt, err := db.Prepare(`INSERT INTO log_entries VALUES (?, ?, ?, ?)`)
	checkFatal(err)

	_, err = stmt.Exec(time.Now(), string(oldState), string(newState), message)
	checkFatal(err)
}

func (l Log) DisplayLog(today bool) {
	db := getLogDb()

	var rows *sql.Rows
	var err error

	if today {
		rows, err = db.Query(`
			SELECT * FROM log_entries
			WHERE date(timestamp) = date('now')
			ORDER BY timestamp DESC`)
	} else {
		rows, err = db.Query(`
			SELECT * FROM log_entries
			ORDER BY timestamp DESC`)
	}

	checkFatal(err)
	defer rows.Close()

	WithPager(func(pipe io.WriteCloser) {
		defer pipe.Close()

		var timestamp time.Time
		var oldState, newState, message string

		for rows.Next() {
			err = rows.Scan(&timestamp, &oldState, &newState, &message)
			checkFatal(err)
			if message == "" {
				fmt.Fprintf(pipe,
					"%s %s -> %s\n",
					timestamp.Format(time.RFC3339),
					oldState,
					newState)
			} else {
				fmt.Fprintf(pipe,
					"%s %s -> %s: %s\n",
					timestamp.Format(time.RFC3339),
					oldState,
					newState,
					message)
			}
		}
	})
}

func applySchema(db *sql.DB) {
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS log_entries (
		timestamp TIMESTAMP,
		old_state TEXT,
		new_state TEXT,
		message TEXT)`)
	checkFatal(err)

	_, err = stmt.Exec()
	checkFatal(err)
}

func getLogDb() *sql.DB {
	db, err := sql.Open("sqlite3", logPath())
	checkFatal(err)

	applySchema(db)

	return db
}

func logPath() string {
	path, err := homedir.Expand(LogPath)
	checkFatal(err)

	return path
}
