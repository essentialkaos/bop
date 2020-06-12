<p align="center"><a href="#readme"><img src="https://gh.kaos.st/bop.svg"/></a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/bop"><img src="https://travis-ci.com/essentialkaos/bop.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-bop-master"><img alt="codebeat badge" src="https://codebeat.co/badges/e0f30749-1508-45dd-8d1c-c074b200f101" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/bop"><img src="https://goreportcard.com/badge/github.com/essentialkaos/bop" alt="GoReportCard" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg" alt="License" /></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#docker-support">Docker support</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

</br>

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

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/bop/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) bop
```

### Docker support

You can use [Docker containers](https://hub.docker.com/r/essentialkaos/bop) for testing your packages. Install latest version of Docker, then:

```bash
curl -#L -o bop-docker https://kaos.sh/bop/bop-docker
chmod +x bop-docker
[sudo] mv bop-docker /usr/bin/
bop-docker test-name package.rpm
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```
[sudo] bop --completion=bash 1> /etc/bash_completion.d/bop
```


ZSH:
```
[sudo] bop --completion=zsh 1> /usr/share/zsh/site-functions/bop
```


Fish:
```
[sudo] bop --completion=fish 1> /usr/share/fish/vendor_completions.d/bop.fish
```

### Usage

```
Usage: bop {options} name package…

Options

  --output, -o file        Output file
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
