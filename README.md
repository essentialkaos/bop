<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/bop/ci"><img src="https://kaos.sh/w/bop/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/bop"><img src="https://kaos.sh/r/bop.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/w/bop/codeql"><img src="https://kaos.sh/w/bop/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

</br>

`bop` is a utility for generating [bibop](https://kaos.sh/bibop) tests for RPM packages.

### Installation

#### From source

To build the `bop` from scratch, make sure you have a working [Go 1.22+](https://github.com/essentialkaos/.github/blob/master/GO-VERSION-SUPPORT.md) workspace (_[instructions](https://go.dev/doc/install)_), then:

```bash
go install github.com/essentialkaos/bop@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/bop/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) bop
```

#### Container image

Official `bop` images available on [GitHub Container Registry](https://kaos.sh/p/bop) and [Docker Hub](https://kaos.sh/d/bop). Install the latest version of [Podman](https://podman.io/getting-started/installation.html) or [Docker](https://docs.docker.com/engine/install/), then:

```bash
curl -#L -o bop-container https://kaos.sh/bop/bop-container
chmod +x bop-container
sudo mv bop-container /usr/bin/bop
bop test-name package.rpm
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo bop --completion=bash 1> /etc/bash_completion.d/bop
```


ZSH:
```bash
sudo bop --completion=zsh 1> /usr/share/zsh/site-functions/bop
```


Fish:
```bash
sudo bop --completion=fish 1> /usr/share/fish/vendor_completions.d/bop.fish
```

### Man documentation

You can generate man page using next command:

```bash
bop --generate-man | sudo gzip > /usr/share/man/man1/bop.1.gz
```

### Usage

<img src=".github/images/usage.svg" />

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
