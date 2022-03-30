package data

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/bop/rpm"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Info contains info about all packages
type Info struct {
	Dist       string
	Pkgs       []string
	Apps       []string
	Configs    []*rpm.Object
	SharedLibs []string
	StaticLibs []*rpm.Object
	Headers    []string
	PkgConfigs []string
	Users      UserMap
	Groups     GroupMap
	Services   []string

	Python2Dirs    []*rpm.Object
	Python2Files   []*rpm.Object
	Python2Modules []string
	Python3Dirs    []*rpm.Object
	Python3Files   []*rpm.Object
	Python3Modules []string
	PythonWheels   []*rpm.Object
}

// UserMap is map user name → user info
type UserMap map[string]*User

// User contains info about user
type User struct {
	Name  string
	UID   string
	GID   string
	Group string
	Home  string
	Shell string
}

// GroupMap is map group name → user info
type GroupMap map[string]*Group

// Group contains info about group
type Group struct {
	Name string
	GID  string
}

// ////////////////////////////////////////////////////////////////////////////////// //
