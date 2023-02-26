<p align="center"><a href="#readme"><img src="https://gh.kaos.st/bop.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/bop/ci"><img src="https://kaos.sh/w/bop/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/bop"><img src="https://kaos.sh/r/bop.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/bop"><img src="https://kaos.sh/b/e0f30749-1508-45dd-8d1c-c074b200f101.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/w/bop/codeql"><img src="https://kaos.sh/w/bop/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#docker-support">Docker support</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

</br>

`bop` is a utility for generating [bibop](https://kaos.sh/bibop) tests for RPM packages.

### Installation

#### From source

To build the `bop` from scratch, make sure you have a working Go 1.18+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go install github.com/essentialkaos/bop
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/bop/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) bop
```

### Docker support

Official `bop` images available on [Docker Hub](https://kaos.sh/d/bop) and [GitHub Container Registry](https://kaos.sh/p/bop). Install the latest version of Docker, then:

```bash
curl -#L -o bop-docker https://kaos.sh/bop/bop-docker
chmod +x bop-docker
sudo mv bop-docker /usr/bin/
bop-docker test-name package.rpm
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```
sudo bop --completion=bash 1> /etc/bash_completion.d/bop
```


ZSH:
```
sudo bop --completion=zsh 1> /usr/share/zsh/site-functions/bop
```


Fish:
```
sudo bop --completion=fish 1> /usr/share/fish/vendor_completions.d/bop.fish
```

### Man documentation

You can generate man page using next command:

```bash
bop --generate-man | sudo gzip > /usr/share/man/man1/bop.1.gz
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

  bop htop htop*.rpm
  Generate simple tests for package

  bop redis redis*.rpm -s redis
  Generate tests with service check

  bop -o zl.recipe zlib zlib*.rpm minizip*.rpm
  Generate tests with custom name
```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/bop/ci.svg?branch=master)](https://kaos.sh/w/bop/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/bop/ci.svg?branch=develop)](https://kaos.sh/w/bop/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
