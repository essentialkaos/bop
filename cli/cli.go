package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v12/env"
	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/pluralize"
	"pkg.re/essentialkaos/ek.v12/timeutil"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/update"

	"github.com/essentialkaos/bop/extractor"
	"github.com/essentialkaos/bop/generator"
	"github.com/essentialkaos/bop/rpm"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "bop"
	VER  = "1.0.1"
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

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_OUTPUT:   {},
	OPT_SERVICE:  {Mergeble: true},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func Init() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) < 2 {
		showUsage()
		return
	}

	name := args[0]
	files := args[1:]

	checkSystem()
	checkFiles(files)
	processFiles(name, files)
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

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	genUsage().Render()
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "name", "package…")

	info.AddOption(OPT_OUTPUT, "Output file", "file")
	info.AddOption(OPT_SERVICE, "List of services for checking {c}(mergable){!}", "service")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample("redis redis*.rpm", "Generate tests for Redis package")

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "bop"))
	case "fish":
		fmt.Printf(fish.Generate(info, "bop"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "bop"))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout prints info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/bop", update.GitHubChecker},
	}

	about.Render()
}
