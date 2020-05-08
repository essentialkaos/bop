package generator

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"

	PATH "pkg.re/essentialkaos/ek.v12/path"
	"pkg.re/essentialkaos/ek.v12/sliceutil"

	"github.com/essentialkaos/bop/data"
	"github.com/essentialkaos/bop/rpm"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Generate generates bibop test data
func Generate(name string, services []string, info *data.Info) (string, string) {
	data := genHeader(name, info)
	data += genDependencies(info)
	data += genOptions(info)
	data += genVariables(info, services)
	data += genEnvCheck(info)
	data += genServicesCheck(info, services)
	data += genSharedLibsCheck(info)
	data += genStaticLibsCheck(info)
	data += genHeadersCheck(info)
	data += genPkgConfigCheck(info)
	data += genPython2ModuleCheck(info)
	data += genPython3ModuleCheck(info)
	data += genPythonWheelsCheck(info)

	fmt.Println(data)

	return genOutputName(name, info), data
}

// ////////////////////////////////////////////////////////////////////////////////// //

// genOutputName generates output file name
func genOutputName(name string, info *data.Info) string {
	osVersion := getOSVersion(info.Dist)

	if osVersion != -1 {
		switch {
		case len(info.Services) != 0,
			len(info.Python2Modules) != 0,
			len(info.Python3Modules) != 0:
			return fmt.Sprintf("%s-c%d.bibop", name, osVersion)
		}
	}

	return fmt.Sprintf("%s.recipe", name)
}

// genHeader generates header
func genHeader(name string, info *data.Info) string {
	var data string

	osVersion := getOSVersion(info.Dist)

	if osVersion > 0 && len(info.Services) != 0 {
		data += fmt.Sprintf("# Bibop recipe for %s for CentOS %d\n", name, osVersion)
	} else {
		data += fmt.Sprintf("# Bibop recipe for %s\n", name)
	}

	data += "# See more: https://kaos.sh/bibop\n\n"

	return data
}

// genDependencies generates dependencies definition
func genDependencies(info *data.Info) string {
	return fmt.Sprintf("pkg %s\n\n", strings.Join(info.Pkgs, " "))
}

// genOptions generates options
func genOptions(info *data.Info) string {
	var data string

	if len(info.Services) == 0 {
		data += "fast-finish yes\n\n"
	} else {
		data += "require-root yes\n\n"
	}

	return data
}

// genVariables generates variables
func genVariables(info *data.Info, services []string) string {
	if getOSVersion(info.Dist) == 6 {
		return ""
	}

	if len(services) == 0 && len(info.Services) != 0 {
		return "var delay 3\n\n"
	}

	for _, service := range info.Services {
		if sliceutil.Contains(services, service) {
			return "var delay 3\n\n"
		}
	}

	return ""
}

// genEnvCheck generates environment checks
func genEnvCheck(info *data.Info) string {
	if isSimpleService(info) {
		return genBasicEnvCheck(info)
	}

	data := genAppsCheck(info)
	data += genConfigsCheck(info)
	data += genUsersAndGroupsCheck(info)
	data += genServicesPresenceCheck(info)

	return data
}

// genServicesCheck generates checks for services
func genServicesCheck(info *data.Info, services []string) string {
	if len(info.Services) == 0 {
		return ""
	}

	var data string

	osVersion := getOSVersion(info.Dist)

	for i := 0; i < 3; i++ {
		for _, service := range info.Services {
			if len(services) > 0 && !sliceutil.Contains(services, service) {
				continue
			}

			switch i {
			case 0:
				data += genServiceStartCheck(service, osVersion)
			case 1:
				data += genServiceStatusCheck(service, osVersion)
			case 2:
				data += genServiceStopCheck(service, osVersion)
			}

			data += "\n"
		}
	}

	return data
}

// genSharedLibsCheck generates checks for shared libs
func genSharedLibsCheck(info *data.Info) string {
	if len(info.SharedLibs) == 0 {
		return ""
	}

	data := `command "-" "Check shared libs"` + "\n"

	for _, lib := range info.SharedLibs {
		data += fmt.Sprintf("  lib-loaded %s\n", lib)
	}

	return data + "\n"
}

// genStaticLibsCheck generates checks for static libs
func genStaticLibsCheck(info *data.Info) string {
	if len(info.StaticLibs) == 0 {
		return ""
	}

	data := `command "-" "Check static libs"` + "\n"

	for _, lib := range info.StaticLibs {
		data += fmt.Sprintf("  exist %s\n", lib.Path)
		data += fmt.Sprintf("  mode %s %o\n\n", lib.Path, lib.Mode)
	}

	return data
}

// genHeadersCheck generates checks for libs headers
func genHeadersCheck(info *data.Info) string {
	if len(info.Headers) == 0 {
		return ""
	}

	data := `command "-" "Check headers"` + "\n"

	for _, header := range info.Headers {
		data += fmt.Sprintf("  lib-header %s\n", header)
	}

	return data + "\n"
}

// genPkgConfigCheck generates checks for pkg-config
func genPkgConfigCheck(info *data.Info) string {
	if len(info.PkgConfigs) == 0 {
		return ""
	}

	data := `command "-" "Check pkg-config"` + "\n"

	for _, cfg := range info.PkgConfigs {
		data += fmt.Sprintf("  lib-config %s\n", cfg)
	}

	return data + "\n"
}

// genPython2ModuleCheck generates checks for Python 2 modules
func genPython2ModuleCheck(info *data.Info) string {
	if len(info.Python2Modules) == 0 {
		return ""
	}

	data := `command "-" "Check Python 2 installation"` + "\n"

	if len(info.Python2Dirs) > 0 {
		for _, dir := range info.Python2Dirs {
			data += fmt.Sprintf("  exist %s\n", getPythonModuleFilePath(dir.Path))
			data += fmt.Sprintf("  dir %s\n\n", getPythonModuleFilePath(dir.Path))
		}
	}

	if len(info.Python2Files) > 0 {
		for _, file := range info.Python2Files {
			data += fmt.Sprintf("  exist %s\n", getPythonModuleFilePath(file.Path))
		}

		data += "\n"
	}

	for _, module := range info.Python2Modules {
		data += fmt.Sprintf("  python-module %s\n", module)
	}

	return data + "\n"
}

// genPython3ModuleCheck generates checks for Python 3 modules
func genPython3ModuleCheck(info *data.Info) string {
	if len(info.Python3Modules) == 0 {
		return ""
	}

	data := `command "-" "Check Python 3 installation"` + "\n"

	if len(info.Python3Dirs) > 0 {
		for _, dir := range info.Python3Dirs {
			data += fmt.Sprintf("  exist %s\n", getPythonModuleFilePath(dir.Path))
			data += fmt.Sprintf("  dir %s\n\n", getPythonModuleFilePath(dir.Path))
		}
	}

	if len(info.Python3Files) > 0 {
		for _, file := range info.Python3Files {
			data += fmt.Sprintf("  exist %s\n", getPythonModuleFilePath(file.Path))
		}

		data += "\n"
	}

	for _, module := range info.Python3Modules {
		data += fmt.Sprintf("  python3-module %s\n", module)
	}

	return data + "\n"
}

// genPythonWheelsCheck generates checks for Python wheels
func genPythonWheelsCheck(info *data.Info) string {
	if len(info.PythonWheels) == 0 {
		return ""
	}

	data := `command "-" "Check Python wheels"` + "\n"

	for _, wheel := range info.PythonWheels {
		data += fmt.Sprintf("  exist %s\n", wheel.Path)
		data += fmt.Sprintf("  perms %s %o\n\n", wheel.Path, wheel.Mode)
	}

	return data
}

// genBasicEnvCheck generates env checks for very simple package
func genBasicEnvCheck(info *data.Info) string {
	if len(info.Apps) == 0 && len(info.Services) == 0 && len(info.Configs) == 0 {
		return ""
	}

	data := `command "-" "Check environment"` + "\n"

	if len(info.Apps) > 0 {
		for _, app := range info.Apps {
			data += fmt.Sprintf("  app %s\n", app)
		}

		data += "\n"
	}

	if len(info.Services) > 0 {
		for _, service := range info.Services {
			data += fmt.Sprintf("  service-present %s\n", service)
		}

		data += "\n"
	}

	if len(info.Configs) > 0 {
		for _, config := range info.Configs {
			data += genConfigCheck(config)
			data += "\n"
		}

		data += "\n"
	}

	return data
}

// genAppsCheck generates checks for applications
func genAppsCheck(info *data.Info) string {
	if len(info.Apps) == 0 {
		return ""
	}

	data := `command "-" "Check apps"` + "\n"

	for _, app := range info.Apps {
		data += fmt.Sprintf("  app %s\n", app)
	}

	return data + "\n"
}

// genConfigsCheck generates checks for configuration files and directories
func genConfigsCheck(info *data.Info) string {
	if len(info.Configs) == 0 {
		return ""
	}

	data := `command "-" "Check configuration files and directories"` + "\n"

	for _, config := range info.Configs {
		data += genConfigCheck(config)
		data += "\n"
	}

	return data
}

// genServicesPresenceCheck generates checks for services presence
func genServicesPresenceCheck(info *data.Info) string {
	if len(info.Services) == 0 {
		return ""
	}

	data := `command "-" "Check services presence"` + "\n"

	for _, service := range info.Services {
		data += fmt.Sprintf("  service-present %s\n", service)
	}

	return data + "\n"
}

// genServiceStartCheck generates checks for service start
func genServiceStartCheck(service string, osVersion int) string {
	var data string

	if osVersion < 7 {
		data = fmt.Sprintf("command \"service %s start\" \"Start %s daemon\"\n", service, service)
		data += "  exit 0\n"
	} else {
		data = fmt.Sprintf("command \"systemctl start %s\" \"Start %s daemon\"\n", service, service)
		data += "  wait {delay}\n"
	}

	data += fmt.Sprintf("  service-works %s\n", service)

	return data
}

// genServiceStatusCheck generates checks for service status check
func genServiceStatusCheck(service string, osVersion int) string {
	var data string

	if osVersion < 7 {
		data = fmt.Sprintf("command \"service %s status\" \"Check status of %s daemon\"\n", service, service)
		data += "  exit 0\n"
	} else {
		data = fmt.Sprintf("command \"systemctl status %s\" \"Check status of %s daemon\"\n", service, service)
		data += "  expect \"active (running)\"\n"
	}

	return data
}

// genServiceStopCheck generates checks for service stop
func genServiceStopCheck(service string, osVersion int) string {
	var data string

	if osVersion < 7 {
		data = fmt.Sprintf("command \"systemctl stop %s\" \"Stop %s daemon\"\n", service, service)
		data += "  exit 0\n"
	} else {
		data = fmt.Sprintf("command \"service %s stop\" \"Stop %s daemon\"\n", service, service)
		data += "  wait {delay}\n"
	}

	data += fmt.Sprintf("  !service-works %s\n", service)

	return data

}

// genUsersAndGroupsCheck generates checks for users and groups
func genUsersAndGroupsCheck(info *data.Info) string {
	if len(info.Users) == 0 && len(info.Groups) == 0 {
		return ""
	}

	data := `command "-" "Check users and groups"` + "\n"

	if len(info.Users) > 0 {
		for _, user := range info.Users {
			data += genUserCheck(user)
		}

		data += "\n"
	}

	if len(info.Groups) > 0 {
		for _, group := range info.Groups {
			data += genGroupCheck(group)
		}

		data += "\n"
	}

	return data
}

// genConfigCheck generates checks for configuration file or directory
func genConfigCheck(config *rpm.Object) string {
	data := fmt.Sprintf("  exist %s\n", config.Path)

	if config.IsDir {
		data = fmt.Sprintf("  dir %s\n", config.Path)

		if config.Mode != 0755 {
			data += fmt.Sprintf("  mode %s %o\n", config.Path, config.Mode)
		}
	} else {
		if config.Mode != 0644 {
			data += fmt.Sprintf("  mode %s %o\n", config.Path, config.Mode)
		}
	}

	if config.User != "" && config.User != "root" {
		data += fmt.Sprintf(
			"  owner %s %s:%s\n",
			config.Path, config.User, config.Group,
		)
	}

	return data
}

// genUserCheck generates checks for given user
func genUserCheck(user *data.User) string {
	data := fmt.Sprintf("  user-exist %s\n", user.Name)

	if user.UID != "" {
		data += fmt.Sprintf("  user-id %s %s\n", user.Name, user.UID)
	}

	if user.GID != "" {
		data += fmt.Sprintf("  user-gid %s %s\n", user.Name, user.GID)
	}

	if user.Group != "" {
		fmt.Sprintf("  user-group %s %s\n", user.Name, user.Group)
	}

	if user.Home != "" {
		fmt.Sprintf("  user-home %s %s\n", user.Name, user.Home)
	}

	if user.Shell != "" {
		fmt.Sprintf("  user-shell %s %s\n", user.Name, user.Shell)
	}

	return data
}

// genGroupCheck generates checks for given group
func genGroupCheck(group *data.Group) string {
	data := fmt.Sprintf("  group-exist %s\n", group.Name)

	if group.GID != "" {
		data += fmt.Sprintf("  group-id %s %s\n", group.Name, group.GID)
	}

	return data
}

// getOSVersion returns OS version number
func getOSVersion(dist string) int {
	switch {
	case strings.HasPrefix(dist, "el6"):
		return 6
	case strings.HasPrefix(dist, "el7"):
		return 7
	case strings.HasPrefix(dist, "el8"):
		return 8
	}

	return -1
}

// getPythonModuleFilePath replaces part of path to variable
func getPythonModuleFilePath(path string) string {
	pathDir := PATH.DirN(path, 4)

	switch {
	case strings.HasPrefix(pathDir, "/usr/lib/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON2_SITELIB}", -1)
	case strings.HasPrefix(pathDir, "/usr/lib64/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON2_SITEARCH}", -1)
	case strings.HasPrefix(pathDir, "/usr/local/lib/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON2_SITELIB_LOCAL}", -1)
	case strings.HasPrefix(pathDir, "/usr/local/lib64/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON2_SITEARCH_LOCAL}", -1)
	case strings.HasPrefix(pathDir, "/usr/lib/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON3_SITELIB}", -1)
	case strings.HasPrefix(pathDir, "/usr/lib64/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON3_SITEARCH}", -1)
	case strings.HasPrefix(pathDir, "/usr/local/lib/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON3_SITELIB_LOCAL}", -1)
	case strings.HasPrefix(pathDir, "/usr/local/lib64/python3"):
		path = strings.Replace(path, pathDir, "{PYTHON3_SITEARCH_LOCAL}", -1)
	}

	return path
}

// isSimpleService checks if given info contains data for very simple package
func isSimpleService(info *data.Info) bool {
	switch {
	case len(info.Apps) > 3,
		len(info.Configs) > 3,
		len(info.Services) > 1,
		len(info.Users) > 1,
		len(info.Groups) > 1:
		return false
	}

	return true
}
