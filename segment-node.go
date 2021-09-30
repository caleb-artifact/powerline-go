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
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Engines map[string]interface{} `json:"engines"`
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
	pkgPath, ok := findPackageJSON(p.cfg.RecursivePackageSearch)
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
	engines := ""

	if pkg.Engines != nil && len(pkg.Engines) > 0 {

		for engine, val := range pkg.Engines {
			switch version := val.(type) {
			case string:
				if len(engines) > 0 {
					engines += ", "
				}
				engines += strings.TrimSpace(engine) + " " + strings.TrimSpace(version)
			}
		}

		engines = " (" + engines + ")"
	}

	str := ""
	if version == "" && name == "" {
		str = "!" + engines
	} else if version == "" || name == "" {
		str = name + version + engines
	} else {
		str = name + "@" + version + engines
	}

	return str
}

func segmentNodeVersion(p *powerline) []pwl.Segment {
	nodeVersion := getNodeVersion()
	if nodeVersion == "" {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "node-version",
		Content:    "\u2B22 " + nodeVersion,
		Foreground: p.theme.NodeVersionFg,
		Background: p.theme.NodeVersionBg,
	}}
}

func segmentNode(p *powerline) []pwl.Segment {
	pkgVersion := getPackageVersionString(p)
	if pkgVersion == "" {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "node-segment",
		Content:    pkgVersion + " \u2B22",
		Foreground: p.theme.NodeFg,
		Background: p.theme.NodeBg,
	}}
}
