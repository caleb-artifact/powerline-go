package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

const pkgfile = "package.json"

type packageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func getNodeVersion() string {
	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(out), "\n")
}

func findPackageJSON(deepSearch bool) (path string, ok bool) {
	pkgPath := "./" + pkgfile
	stat, err := os.Stat(pkgPath)
	if err == nil && !stat.IsDir() {
		return pkgPath, true
	} else if !deepSearch || (!os.IsNotExist(err) && !os.IsPermission(err)) {
		return "", false
	}

	prevDir, err := os.Stat(".")
	if err != nil {
		return "", false
	}
	for parent := ".."; ; parent = "../" + parent {
		if len(parent) > 1024 {
			return "", false
		}
		stat, err = os.Stat(parent)
		if err != nil {
			return "", false
		} else if os.SameFile(stat, prevDir) {
			return "", false
		}
		prevDir = stat

		stat, err = os.Stat(parent + "/" + pkgfile)
		if err == nil && !stat.IsDir() {
			return parent + "/" + pkgfile, true
		} else if err != nil && !os.IsNotExist(err) && !os.IsPermission(err) {
			return "", false
		}
	}
}

func getPackageVersionString(p *powerline) string {
	pkgPath, ok := findPackageJSON(p.cfg.NodeDeepPackageSearch)
	if !ok {
		return ""
	}
	pkg := packageJSON{}
	raw, err := ioutil.ReadFile(pkgPath)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(raw, &pkg)
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(pkg.Version)
	name := strings.TrimSpace(pkg.Name)

	if version == "" && name == "" {
		return "!"
	} else if version == "" || name == "" {
		return name + version
	} else {
		return name + "@" + version
	}
}

func segmentNode(p *powerline) []pwl.Segment {
	nodeVersion := getNodeVersion()
	pkgVersion := getPackageVersionString(p)

	segments := []pwl.Segment{}

	if nodeVersion != "" {
		segments = append(segments, pwl.Segment{
			Name:       "node",
			Content:    "\u2B22 " + nodeVersion,
			Foreground: p.theme.NodeVersionFg,
			Background: p.theme.NodeVersionBg,
		})
	}

	if pkgVersion != "" {
		segments = append(segments, pwl.Segment{
			Name:       "node-segment",
			Content:    pkgVersion + " \u2B22",
			Foreground: p.theme.NodeFg,
			Background: p.theme.NodeBg,
		})
	}

	return segments
}
