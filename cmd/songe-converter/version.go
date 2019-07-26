package main

import (
	"fmt"
	"runtime"
)

var (
	sha1ver   = "unknown"
	gitTag    string
	buildTime string

	printVer bool
)

type versionRow struct {
	label string
	value string
}

func printVersionInfo() {
	versionRows := make([]versionRow, 0)
	addRow := func(l, v string) {
		row := versionRow{label: l, value: v}
		versionRows = append(versionRows, row)
	}

	var version string
	if gitTag == "" {
		version = "dev"
	} else {
		version = gitTag
	}

	addRow("Version", version)
	addRow("Git Hash", sha1ver)
	addRow("Go Version", runtime.Version())

	if buildTime != "" {
		addRow("Build Time", buildTime)
	}

	var widest int
	for _, r := range versionRows {
		width := len(r.label) + 2
		if width > widest {
			widest = width
		}
	}

	for _, r := range versionRows {
		fmt.Printf("%*s %s\n", widest*-1, r.label, r.value)
	}
}
