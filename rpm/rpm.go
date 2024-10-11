package rpm

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v13/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Package contains package info
type Package struct {
	Name       string
	Dist       string
	Scriptlets string
	Payload    []*Object
	IsSrc      bool
}

// Object contains info about payload object
type Object struct {
	Path     string
	User     string
	Group    string
	Mode     os.FileMode
	IsConfig bool
	IsDir    bool
	IsLink   bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ReadRPM reads info from package
func ReadRPM(file string) (*Package, error) {
	var err error

	pkg := &Package{}

	pkg.Name, pkg.Dist, pkg.IsSrc, err = extractPackageInfo(file)

	if err != nil {
		return nil, err
	}

	pkg.Payload, err = extractPayloadInfo(file)

	if err != nil {
		return nil, err
	}

	pkg.Scriptlets, err = extractScriptlets(file)

	if err != nil {
		return nil, err
	}

	return pkg, nil
}

// IsPackage returns true if given file is an rpm package
func IsPackage(file string) bool {
	_, err := execRPMCommand("-qp", file)
	return err == nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// String returns string representation of package
func (p *Package) String() string {
	return fmt.Sprintf(
		"{Name: %s | Dist: %s | Payload: %d | Scriplets: %d}",
		p.Name, p.Dist, len(p.Payload), len(p.Scriptlets),
	)
}

// String returns string representation of payload object
func (o *Object) String() string {
	user := o.User
	group := o.Group
	isConfig := "N"
	isDir := "N"
	isLink := "N"

	if user == "" {
		user = "-"
	}

	if group == "" {
		group = "-"
	}

	if o.IsConfig {
		isConfig = "Y"
	}

	if o.IsDir {
		isDir = "Y"
	}

	if o.IsLink {
		isLink = "Y"
	}

	return fmt.Sprintf(
		"{Path: %s | Mode: %s | User: %s | Group: %s | Config: %s | Dir: %s | Link: %s}",
		o.Path, o.Mode, user, group, isConfig, isDir, isLink,
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// extractPayloadInfo extracts info about package payload
func extractPayloadInfo(file string) ([]*Object, error) {
	dumpData, err := execRPMCommand("-qp", "--dump", file)

	if err != nil {
		return nil, err
	}

	return parseDumpData(dumpData)
}

// extractScriptlets extracts raw scriptlets data
func extractScriptlets(file string) (string, error) {
	return execRPMCommand("-qp", "--scripts", file)
}

// extractPackageInfo extracts package name, dist and src package flag
func extractPackageInfo(file string) (string, string, bool, error) {
	data, err := execRPMCommand("-qp", "--qf", "%{name} %{release} %{sourcepackage}", file)

	if err != nil {
		return "", "", false, err
	}

	name := strutil.ReadField(data, 0, false, ' ')
	dist := extractDist(strutil.ReadField(data, 1, false, ' '))
	isSrc := strutil.ReadField(data, 2, false, ' ') == "1"

	return name, dist, isSrc, nil
}

// parseDumpData parses dump data
func parseDumpData(data string) ([]*Object, error) {
	r := strings.NewReader(data)
	s := bufio.NewScanner(r)

	var payload []*Object

	for s.Scan() {
		payload = append(payload, parsePayloadInfo(s.Text()))
	}

	return payload, nil
}

// parsePayloadInfo parses payload object info
func parsePayloadInfo(data string) *Object {
	modeStr := strutil.Tail(strutil.ReadField(data, 4, false, ' '), 4)
	modeUint, _ := strconv.ParseUint(modeStr, 8, 32)
	emptyHash := strings.Trim(strutil.ReadField(data, 3, false, ' '), "0") == ""
	link := strutil.ReadField(data, 10, false, ' ')

	return &Object{
		Path:     strutil.ReadField(data, 0, false, ' '),
		User:     strutil.ReadField(data, 5, false, ' '),
		Group:    strutil.ReadField(data, 6, false, ' '),
		Mode:     os.FileMode(modeUint),
		IsConfig: strutil.ReadField(data, 7, false, ' ') == "1",
		IsDir:    link == "X" && emptyHash,
		IsLink:   link != "X" && emptyHash,
	}
}

// execRPMCommand executes rpm command with given options
func execRPMCommand(options ...string) (string, error) {
	output, err := exec.Command("rpm", options...).Output()
	return string(output), err
}

// extractDist extracts dist data from release info
func extractDist(data string) string {
	dotIndex := strings.LastIndex(data, ".")
	return strutil.Substring(data, dotIndex+1, 9999)
}
