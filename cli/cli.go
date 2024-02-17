package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/env"
	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/pluralize"
	"github.com/essentialkaos/ek/v12/terminal/tty"
	"github.com/essentialkaos/ek/v12/timeutil"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/bop/cli/support"
	"github.com/essentialkaos/bop/extractor"
	"github.com/essentialkaos/bop/generator"
	"github.com/essentialkaos/bop/rpm"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "bop"
	VER  = "1.3.0"
	DESC = "Utility for generating formal bibop tests for RPM packages"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_OUTPUT   = "o:output"
	OPT_SERVICE  = "s:service"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_OUTPUT:   {},
	OPT_SERVICE:  {Mergeble: true},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.BOOL},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

var colorTagApp, colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print()
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Print(APP, VER, gitRev, gomod)
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) < 2:
		genUsage().Print()
		os.Exit(0)
	}

	name := args.Get(0).String()
	files := args.Strings()[1:]

	checkSystem()
	checkFiles(files)
	processFiles(name, files)
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{#9966CC}", "{#9966CC}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#140}", "{#140}"
	default:
		colorTagApp, colorTagVer = "{*}{m}", "{m}"
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// checkSystem checks system
func checkSystem() {
	if env.Which("rpm") == "" {
		printErrorAndExit("rpm utility is mandatory for this application")
	}
}

// processFiles runs files processing
func processFiles(name string, files []string) {
	fmtc.Printf(
		"Generating {#85}bibop{!} tests for {*}%s{!} based on given %s…\n",
		name, pluralize.P("%s (%d)", len(files), "package", "packages"),
	)

	start := time.Now()

	info, err := extractor.ProcessPackages(files)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	services := parseServiceList(options.GetS(OPT_SERVICE))
	output, data := generator.Generate(name, services, info)

	if options.Has(OPT_OUTPUT) {
		output = options.GetS(OPT_OUTPUT)
	}

	err = ioutil.WriteFile(output, []byte(data), 0644)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	fmtc.Printf(
		"{*}Recipe saved as {#85}%s{!} {s-}(processing took %s){!}\n",
		output, timeutil.PrettyDuration(time.Since(start)),
	)
}

// checkFiles checks input files
func checkFiles(files []string) {
	var hasErrors bool

	for _, file := range files {
		switch {
		case !fsutil.IsExist(file):
			printError("%s does not exist", file)
			hasErrors = true
		case !fsutil.IsReadable(file):
			printError("%s is not readable", file)
			hasErrors = true
		case !rpm.IsPackage(file):
			printError("%s is not an rpm package", file)
			hasErrors = true
		}
	}

	if hasErrors {
		os.Exit(1)
	}
}

// parseServiceList parses service list option data
func parseServiceList(data string) []string {
	if data == "" {
		return nil
	}

	if strings.Contains(data, ",") {
		data = strings.Replace(data, ",", " ", -1)
	}

	return strings.Fields(data)
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit print error message and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "bop"))
	case "fish":
		fmt.Printf(fish.Generate(info, "bop"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "bop"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "name", "package…")

	info.AppNameColorTag = colorTagApp

	info.AddOption(OPT_OUTPUT, "Output file", "file")
	info.AddOption(OPT_SERVICE, "List of services for checking {c}(mergeable){!}", "service")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample("htop htop*.rpm", "Generate simple tests for package")
	info.AddExample("redis redis*.rpm -s redis", "Generate tests with service check")
	info.AddExample("-o zl.recipe zlib zlib*.rpm minizip*.rpm", "Generate tests with custom name")

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/bop", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
