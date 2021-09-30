package main

import (
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentNPM(p *powerline) []pwl.Segment {
	out, err := exec.Command("npm", "--version").Output()
	if err != nil {
		return []pwl.Segment{}
	}
	npmVersion := strings.TrimSuffix(string(out), "\n")

	if npmVersion == "" {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "npm-version",
		Content:    "npm " + npmVersion,
		Foreground: p.theme.NPMFg,
		Background: p.theme.NPMBg,
	}}
}
