package main

import (
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentYarn(p *powerline) []pwl.Segment {
	out, err := exec.Command("yarn", "--version").Output()
	if err != nil {
		return []pwl.Segment{}
	}
	version := strings.TrimSuffix(string(out), "\n")

	if version == "" {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "yarn-version",
		Content:    "\U000130E0 " + version,
		Foreground: p.theme.YarnFg,
		Background: p.theme.YarnBg,
	}}
}
