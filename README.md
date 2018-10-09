# Gomodoro

Command-line pomodoro application.

## Install

```
go get github.com/tokenshift/gomo
go install github.com/tokenshift/gomo
```

Or download a pre-packaged binary at https://github.com/tokenshift/gomo/releases.

## Use

Start a work session:

```
$ gomo work
```

Start a break:

```
$ gomo break
```

Display pomodoro status:

```
$ gomo status
Working (19.3 minutes remaining)
```

For command-line help:

```
$ gomo --help
```

Or add the following to your PS1 to get (and update) the current pomodoro
status automatically at every prompt:

```
export PS1='...$(gomo status --auto)...'
```

## Config

Gomodoro stores status and configuration information under `~/.gomo`. You can
edit your configuration by using the `gomo config` command, or by editing
`~/.gomo/config`. Configuration values include:

* `WorkSessionMinutes`  
  How long each work session is (default: 25).
* `ShortBreakMinutes`  
  How long a short break is (default: 5).
* `LongBreakMinutes`  
  How long a long break is (default: 15).
* `LongBreakInterval`  
  How many work sessions between each long break (default: 4).
