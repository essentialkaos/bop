package extractor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	PATH "github.com/essentialkaos/ek/v13/path"
	"github.com/essentialkaos/ek/v13/sliceutil"
	"github.com/essentialkaos/ek/v13/strutil"

	"github.com/essentialkaos/bop/data"
	"github.com/essentialkaos/bop/rpm"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var pkgConfigGlobs = []string{
	"/usr/lib/pkgconfig/*.pc",
	"/usr/lib64/pkgconfig/*.pc",
}

var staticLibsGlobs = []string{
	"/usr/lib/*.a",
	"/usr/lib64/*.a",
}

var sharedLibsGlobs = []string{
	"/usr/lib/*.so.*",
	"/usr/lib64/*.so.*",
}

var systemdUnitGlobs = []string{
	"/usr/lib/systemd/system/*.service",
	"/usr/lib/systemd/user/*.service",
}

var includeDir = "/usr/include"

// ////////////////////////////////////////////////////////////////////////////////// //

// ProcessPackages reads rpm files and extracts info from them
func ProcessPackages(files []string) (*data.Info, error) {
	pkgs, err := readPackagesData(files)
	if err != nil {
		return nil, err
	}

	if isPackagesWithMixedDist(pkgs) {
		return nil, fmt.Errorf("Packages for different versions of OS can not be used for test generation")
	}

	if err != nil {
		return nil, err
	}

	return extractPackagesInfo(pkgs), nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// readPackagesData read info packages info from rpm files
func readPackagesData(files []string) ([]*rpm.Package, error) {
	var pkgs []*rpm.Package

	for _, file := range files {
		pkg, err := rpm.ReadRPM(file)

		if err != nil {
			return nil, err
		}

		if pkg.IsSrc {
			continue
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

// extractPackagesInfo extracts info from packages
func extractPackagesInfo(pkgs []*rpm.Package) *data.Info {
	info := &data.Info{
		Users:  make(map[string]*data.User),
		Groups: make(map[string]*data.Group),
	}

	for _, pkg := range pkgs {
		addPackageInfo(info, pkg)
	}

	return info
}

// addPackageInfo extracts info from package
func addPackageInfo(info *data.Info, pkg *rpm.Package) {
	info.Pkgs = append(info.Pkgs, pkg.Name)
	info.Dist = pkg.Dist

	addAppsInfo(info, pkg)
	addConfigsInfo(info, pkg)
	addCompletions(info, pkg)
	addLibsInfo(info, pkg)
	addHeadersInfo(info, pkg)
	addPkgConfigsInfo(info, pkg)
	addOwnersInfo(info, pkg)
	addServicesInfo(info, pkg)
	addPython2ModulesInfo(info, pkg)
	addPython3ModulesInfo(info, pkg)
	addPythonWheels(info, pkg)

	sort.Strings(info.Pkgs)
	sort.Strings(info.Apps)
	sort.Strings(info.PkgConfigs)
	sort.Strings(info.SharedLibs)
	sort.Strings(info.Headers)
	sort.Strings(info.Services)
	sort.Strings(info.Python2Modules)
	sort.Strings(info.Python3Modules)

	info.Services = sliceutil.Deduplicate(info.Services)

	if pkg.Scriptlets != "" {
		extractUsersData(pkg.Scriptlets, info.Users)
		extractGroupsData(pkg.Scriptlets, info.Groups)
	}
}

// addAppsInfo extracts info about applications from package info
func addAppsInfo(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		switch {
		case obj.IsDir:
			continue
		case obj.Mode|0111 == 0:
			continue
		}

		switch {
		case strings.HasPrefix(obj.Path, "/usr/bin/"),
			strings.HasPrefix(obj.Path, "/usr/sbin/"),
			strings.HasPrefix(obj.Path, "/bin/"),
			strings.HasPrefix(obj.Path, "/sbin/"):
			info.Apps = append(info.Apps, path.Base(obj.Path))
		}
	}
}

// addConfigsInfo extracts info about configuration files from package info
func addConfigsInfo(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		if obj.IsConfig {
			info.Configs = append(info.Configs, obj)
		}
	}
}

// addCompletions extracts info about shell completions
func addCompletions(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		if strutil.HasPrefixAny(
			obj.Path,
			"/usr/share/bash-completion/completions",
			"/usr/share/fish/vendor_completions.d",
			"/usr/share/zsh/site-functions",
		) {
			info.Completions = append(info.Completions, obj.Path)
		}
	}
}

// addPkgConfigsInfo extracts info about package configuration files
// from package info
func addPkgConfigsInfo(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		if matchAnyGlob(obj.Path, pkgConfigGlobs) {
			cfgName := strutil.Exclude(PATH.Base(obj.Path), ".pc")
			info.PkgConfigs = append(info.PkgConfigs, cfgName)
		}
	}
}

// addLibsInfo extracts info about libs from package info
func addLibsInfo(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		if matchAnyGlob(obj.Path, sharedLibsGlobs) && !obj.IsLink {
			info.SharedLibs = append(info.SharedLibs, formatLibName(obj.Path))
		}

		if matchAnyGlob(obj.Path, staticLibsGlobs) && !obj.IsLink {
			info.StaticLibs = append(info.StaticLibs, obj)
		}
	}
}

// addHeadersInfo extracts info about libs headers
func addHeadersInfo(info *data.Info, pkg *rpm.Package) {
	headers := make(map[string]bool)

	for _, obj := range pkg.Payload {
		if !strings.HasPrefix(obj.Path, includeDir) {
			continue
		}

		headerDir := PATH.DirN(strutil.Exclude(obj.Path, includeDir+"/"), 1)

		headers[headerDir] = true
	}

	if len(headers) == 0 {
		return
	}

	for headerDir := range headers {
		info.Headers = append(info.Headers, headerDir)
	}
}

// addOwnersInfo extracts info about users from package info
func addOwnersInfo(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		if obj.User != "root" {
			info.Users[obj.User] = &data.User{Name: obj.User}
		}

		if obj.Group != "root" {
			info.Groups[obj.Group] = &data.Group{Name: obj.Group}
		}
	}
}

// addServicesInfo extracts info about service from package info
func addServicesInfo(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		if matchAnyGlob(obj.Path, systemdUnitGlobs) {
			service := strutil.Exclude(path.Base(obj.Path), ".service")
			info.Services = append(info.Services, service)
		}

		if strings.HasPrefix(obj.Path, "/etc/rc.d/init.d/") {
			info.Services = append(info.Services, path.Base(obj.Path))
		}
	}
}

// addPython2ModulesInfo adds info about Python 2 modules
func addPython2ModulesInfo(info *data.Info, pkg *rpm.Package) {
	modules := make(map[string]bool)

	for _, obj := range pkg.Payload {
		dir, ok := isPythonModuleObject(obj.Path, "2")

		if !ok || strings.HasSuffix(obj.Path, ".egg-info") {
			continue
		}

		if obj.IsDir && isValidPythonModuleDir(obj.Path) {
			modules[extractPythonModuleName(obj.Path, dir)] = true
			info.Python2Dirs = append(info.Python2Dirs, obj)
			continue
		}

		if strings.HasSuffix(obj.Path, "__init__.py") {
			info.Python2Files = append(info.Python2Files, obj)
			continue
		}
	}

	info.Python2Modules = append(info.Python2Modules, mapToSlice(modules)...)
}

// addPython3ModulesInfo adds info about Python 3 modules
func addPython3ModulesInfo(info *data.Info, pkg *rpm.Package) {
	modules := make(map[string]bool)

	for _, obj := range pkg.Payload {
		dir, ok := isPythonModuleObject(obj.Path, "3")

		if !ok || strings.HasSuffix(obj.Path, ".egg-info") {
			continue
		}

		if obj.IsDir && isValidPythonModuleDir(obj.Path) {
			modules[extractPythonModuleName(obj.Path, dir)] = true
			info.Python3Dirs = append(info.Python3Dirs, obj)
			continue
		}

		if strings.HasSuffix(obj.Path, "__init__.py") {
			info.Python3Files = append(info.Python3Files, obj)
			continue
		}
	}

	info.Python3Modules = append(info.Python3Modules, mapToSlice(modules)...)
}

// addPythonWheels adds info about Python wheels
func addPythonWheels(info *data.Info, pkg *rpm.Package) {
	for _, obj := range pkg.Payload {
		if strings.HasSuffix(obj.Path, ".whl") {
			info.PythonWheels = append(info.PythonWheels, obj)
		}
	}
}

// isPackagesWithMixedDist returns true if given package set contains packages for
// different OS versions
func isPackagesWithMixedDist(pkgs []*rpm.Package) bool {
	var dist string

	for _, pkg := range pkgs {
		if dist != "" && dist != pkg.Dist {
			return true
		}

		dist = pkg.Dist
	}

	return false
}

// formatLibName formats lib name to glob
func formatLibName(file string) string {
	basename := path.Base(file)
	soIndex := strings.Index(basename, ".so.")

	return strutil.Substring(basename, 0, soIndex) + ".so.*"
}

// matchAnyGlob returns true if given string matches for any of given patterns
func matchAnyGlob(name string, patterns []string) bool {
	for _, pattern := range patterns {
		match, _ := filepath.Match(pattern, name)

		if match {
			return true
		}
	}

	return false
}

// extractUsersData extracts lines with useradd commands from scriptles
func extractUsersData(data string, users map[string]*data.User) {
	lines := extractLines(data, "useradd")

	for _, line := range lines {
		user := parseUserAddCommand(line)

		if user != nil {
			users[user.Name] = user
		}
	}
}

// extractGroupsData extracts lines with groupadd commands from scriptles
func extractGroupsData(data string, groups map[string]*data.Group) {
	lines := extractLines(data, "groupadd")

	for _, line := range lines {
		group := parseGroupAddCommand(line)

		if group != nil {
			groups[group.Name] = group
		}
	}
}

// extractLines extracts lines with given command
func extractLines(data, command string) []string {
	r := strings.NewReader(data)
	s := bufio.NewScanner(r)

	var fullLine string
	var result []string

	for s.Scan() {
		line := s.Text()

		if !strings.Contains(line, command+" ") && fullLine == "" {
			continue
		}

		if strings.HasSuffix(strings.Trim(line, " "), "\\") {
			fullLine += strings.TrimRight(strings.Trim(line, " "), "\\")
			continue
		}

		fullLine += strings.Trim(line, " ")

		result = append(result, fullLine)

		fullLine = ""
	}

	return result
}

// parseUserAddCommand parses useradd command
func parseUserAddCommand(command string) *data.User {
	result := &data.User{}
	gi := strings.Index(command, "useradd ")

	if gi == -1 {
		return nil
	}

	command = command[gi+8:]

	var isComment bool

	for i := 0; i < 20; i++ {
		option := strutil.ReadField(command, i, true, ' ')

		if strings.Contains(option, "\"") || strings.Contains(option, "'") {
			isComment = !isComment
			continue
		}

		if isComment {
			continue
		}

		switch option {
		case "-D", "--defaults", "-m", "--create-home", "-l", "--no-log-init",
			"-M", "--no-create-home", "-N", "--no-user-group", "-o", "--non-unique",
			"-r", "--system", "-U", "--user-group", "-c", "--comment":
			continue // ignore option
		case "-e", "--expiredate", "-f", "--inactive",
			"-k", "--skel", "-K", "--key", "-p", "--password", "-R", "--root",
			"-P", "--prefix", "-Z", "--selinux-user":
			i++
			continue // ignore option and value
		case "-d", "--home-dir":
			result.Home = strutil.ReadField(command, i+1, true, ' ')
			i++
		case "-g", "--gid":
			result.GID = strutil.ReadField(command, i+1, true, ' ')
			i++
		case "-u", "--uid":
			result.UID = strutil.ReadField(command, i+1, true, ' ')
			i++
		case "-s", "--shell":
			result.Shell = strutil.ReadField(command, i+1, true, ' ')
			i++
		case "-G", "--groups":
			result.Group = strutil.ReadField(command, i+1, true, ' ')
			i++
		default:
			result.Name = option
			return result
		}
	}

	return result
}

// parseGroupAddCommand parses groupadd command
func parseGroupAddCommand(command string) *data.Group {
	result := &data.Group{}
	gi := strings.Index(command, "groupadd ")

	if gi == -1 {
		return nil
	}

	command = command[gi+9:]

	for i := 0; i < 10; i++ {
		option := strutil.ReadField(command, i, true, ' ')

		switch option {
		case "-f", "--force", "-o", "-non-unique", "-r", "--system":
			continue // ignore option
		case "-K", "--key", "-p", "--password", "-R", "--root", "-P", "--prefix":
			i++
			continue // ignore option and value
		case "-g", "--gid":
			result.GID = strutil.ReadField(command, i+1, true, ' ')
			i++
		default:
			result.Name = option
			return result
		}
	}

	return result
}

// isPythonModuleObject checks if given object is a part of Python module
func isPythonModuleObject(path, version string) (string, bool) {
	switch {
	case strings.HasPrefix(path, "/usr/lib64/python"+version),
		strings.HasPrefix(path, "/usr/lib/python"+version):
		return PATH.DirN(path, 3), true
	case strings.HasPrefix(path, "/usr/local/lib64/python"+version),
		strings.HasPrefix(path, "/usr/local/lib/python"+version):
		return PATH.DirN(path, 4), true
	}

	return "", false
}

// extractPythonModuleName extracts module name from path
func extractPythonModuleName(path, dir string) string {
	path = strutil.Exclude(path, dir+"/site-packages/")
	slashIndex := strings.IndexRune(path, '/')

	if slashIndex == -1 {
		return path
	}

	return path[:slashIndex]
}

// isValidPythonModuleDir returns true if dir should be used for tests
func isValidPythonModuleDir(dir string) bool {
	dirName := PATH.Base(dir)
	return !strings.HasPrefix(dirName, "__")
}

// mapToSlice converts map to slice
func mapToSlice(m map[string]bool) []string {
	var result []string

	for name := range m {
		result = append(result, name)
	}

	return result
}
