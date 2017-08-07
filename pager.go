package main

import (
	"io"
	"os"
	"os/exec"
)

const DefaultPager = "less"

type PagedFunc func(io.WriteCloser)

func SelectPager() string {
	pager := os.Getenv("PAGER")
	if pager == "" {
		return DefaultPager
	} else {
		return pager
	}
}

func WithPager(f PagedFunc) {
	pager := SelectPager()
	pager, err := exec.LookPath(pager)
	if err != nil {
		f(os.Stdout)
	} else {
		cmd := exec.Command(pager)

		in, out := io.Pipe()

		cmd.Stdin = in
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Start()
		go func() {
			f(out)
		}()
		cmd.Wait()
	}
}
