package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// Given a file or directory, finds which filesystem it is on,
// by parsing /proc/self/mountinfo and comparing dev_t
// of the file to that in the mountinfo field.
func GetFilesystemType(path string) (string, error) {
	var st syscall.Stat_t

	err := syscall.Stat(path, &st)
	if err != nil {
		return "", err
	}

	return getFSTypeByDev(st.Dev)
}

// convert minor:major string from /proc/self/mountinfo into dev_t
func parseDev(s string) uint64 {
	var major uint32
	var minor uint32

	n, _ := fmt.Sscanf(s, "%d:%d", &major, &minor)

	if n != 2 {
		return 0
	}

	return uint64(major<<8 + minor)
}

func getFSTypeByDev(dev uint64) (string, error) {
	mi, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}
	defer mi.Close()

	sc := bufio.NewScanner(mi)
	for sc.Scan() {
		line := strings.Split(sc.Text(), " ")
		if len(line) < 10 {
			return "", fmt.Errorf("Short line in /proc/self/mountinfo: %v\n", line)
		}
		dstr := line[2] // major:minor: value of st_dev for files on filesystem
		fs := line[8]   // filesystem type:  name of filesystem of the form "type[.subtype]"
		d := parseDev(dstr)
		if d == 0 {
			return "", fmt.Errorf("Can't parse device %s", dstr)
		}
		if d == dev {
			return fs, nil
		}
	}

	return "", nil
}
