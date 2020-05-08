<p align="center"><a href="#readme"><img src="https://gh.kaos.st/bop.svg"/></a></p>

<p align="center"><a href="#screenshots">Screenshots</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#thanks">Thanks</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/bop"><img src="https://travis-ci.com/essentialkaos/bop.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/bop"><img src="https://goreportcard.com/badge/github.com/essentialkaos/bop" alt="GoReportCard" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg" alt="License" /></a>
</p>

`bop` is a utility for generating [bibop](https://kaos.sh/bibop) tests for RPM packages.

### Installation

#### From source

Before the initial install, allow git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `bop` from scratch, make sure you have a working Go 1.13+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/bop
```

If you want to update `bop` to latest stable release, do:

```
go get -u github.com/essentialkaos/bop
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/bop/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) bop
```

### Usage

```
Usage: bop {options} name package…

Options

  --service, -s service    List of services for checking (mergable)
  --no-color, -nc          Disable colors in output
  --help, -h               Show this help message
  --version, -v            Show version

Examples

  bop redis redis*.rpm
  Generate tests for Redis package

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.com/essentialkaos/bop.svg?branch=master)](https://travis-ci.com/essentialkaos/bop) |
| `develop` | [![Build Status](https://travis-ci.com/essentialkaos/bop.svg?branch=develop)](https://travis-ci.com/essentialkaos/bop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
