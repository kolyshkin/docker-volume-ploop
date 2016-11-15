package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
)

// Err is a structure used to return errors from vstorage execution
type Err struct {
	c int
	s string
}

// Error returns a string representation of a vstorage error
func (e *Err) Error() string {
	return fmt.Sprintf("vstorage error %d: %s", e.c, e.s)
}

func vstorageRunCmd(stdout io.Writer, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("vstorage", args...)
	cmd.Stdout = stdout
	cmd.Stderr = &stderr

	logrus.Debugf("Run: %s\n", strings.Join([]string{cmd.Path, strings.Join(cmd.Args[1:], " ")}, " "))

	err := cmd.Run()
	if err == nil {
		return nil
	}

	// Command returned an error, get the first line of stderr
	errStr, _ := stderr.ReadString('\n')

	// Get the exit code (Unix-specific)
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			errCode := status.ExitStatus()
			return &Err{c: errCode, s: errStr}
		}
	}
	// unknown exit code
	return &Err{c: -1, s: errStr}
}

func vstorage(args ...string) error {
	return vstorageRunCmd(nil, args...)
}

func vstorageOut(args ...string) (string, error) {
	var stdout bytes.Buffer

	ret := vstorageRunCmd(&stdout, args...)
	out := stdout.String()
	// logrus.Debugf("%s", out)

	return out, ret
}

func vstorageSetTier(path string, tier int8) error {
	var err error

	// TODO: ignore if path is not on vstorage
	if tier >= 0 {
		arg := fmt.Sprintf("tier=%d", tier)
		err = vstorage("set-attr", "-R", path, arg)
	}

	return err
}

// Check if a file/directory is actually on vstorage
func isOnVstorage(path string) bool {
	fs, err := GetFilesystemType(path)
	if err != nil {
		logrus.Errorf("Can't figure %s fs: %v", err)
		return false
	}

	return fs == "fuse.vstorage"
}
