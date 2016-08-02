package oom

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Report struct {
	ProcessName string
	ProcessID   uint
	KernelLogs  []byte
}

func ParseReportFromKernelLogs(kl []byte) (*Report, error) {
	report := &Report{
		KernelLogs: kl,
	}

	if name, err := findProcessName(kl); err != nil {
		return nil, err
	} else {
		report.ProcessName = name
	}

	if pid, err := findProcessID(kl); err != nil {
		return nil, err
	} else {
		report.ProcessID = pid
	}

	return report, nil
}

func findProcessID(kl []byte) (uint, error) {
	match := regexp.MustCompile("(?m)PID: ([0-9]+)").FindSubmatch(kl)
	if match == nil {
		return 0, errors.New("Could not find process ID")
	}
	if len(match) != 2 {
		return 0, errors.New("Got too many matches for process ID")
	}
	pid, err := strconv.Atoi(string(match[1]))
	if err != nil {
		return 0, fmt.Errorf("Cannot convert found process ID: %s", err)
	}
	return uint(pid), nil
}

func findProcessName(kl []byte) (string, error) {
	match := regexp.MustCompile("(?m)([a-z0-9_-]+) invoked oom-killer").FindSubmatch(kl)
	if match == nil {
		return "", errors.New("Could not find process name")
	}
	if len(match) != 2 {
		return "", errors.New("Got too many matches for process name")
	}
	return string(match[1]), nil
}
