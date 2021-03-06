# Elvish: Expressive Programming Language + Versatile Interactive Shell

[![CI status](https://github.com/elves/elvish/workflows/CI/badge.svg)](https://github.com/elves/elvish/actions?query=workflow%3ACI)
[![FreeBSD test status](https://img.shields.io/cirrus/github/elves/elvish?logo=Cirrus%20CI&label=freebsd&task=Test%20on%20FreeBSD)](https://cirrus-ci.com/github/elves/elvish/master)
[![gccgo test status](https://img.shields.io/cirrus/github/elves/elvish?logo=Cirrus%20CI&label=gccgo&task=Test%20on%20gccgo)](https://cirrus-ci.com/github/elves/elvish/master)
[![Test Coverage](https://img.shields.io/codecov/c/github/elves/elvish/master.svg?logo=Codecov&label=coverage)](https://app.codecov.io/gh/elves/elvish/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/elves/elvish)](https://goreportcard.com/report/src.elv.sh)
[![GoDoc](https://img.shields.io/badge/godoc-api-blue.svg)](https://godoc.elv.sh)
[![Twitter](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/ElvishShell)

Elvish is an expressive programming language and a versatile interactive shell,
combined into one seamless package. It runs on Linux, BSDs, macOS and Windows.

Despite its pre-1.0 status, it is already suitable for most daily interactive
use.

**Visit the official website https://elv.sh for prebuilt binaries, blog posts,
documentation and other resources.**

User groups (all connected thanks to [Matrix](https://matrix.org)):
[![Gitter](https://img.shields.io/badge/gitter-elves/elvish-blue.svg?logo=gitter-white)](https://gitter.im/elves/elvish)
[![Telegram Group](https://img.shields.io/badge/telegram-@elvish-blue.svg)](https://telegram.me/elvish)
[![#elvish on freenode](https://img.shields.io/badge/freenode-%23elvish-blue.svg)](https://webchat.freenode.net/?channels=elvish)
[![#users:elves.sh](https://img.shields.io/badge/matrix-%23users:elv.sh-blue.svg)](https://matrix.to/#/#users:elves.sh)

## Building Elvish

Most users do not need to build Elvish from source. Prebuilt binaries for the
latest commit are provided for
[Linux amd64](https://dl.elv.sh/linux-amd64/elvish-HEAD.tar.gz),
[macOS amd64](https://dl.elv.sh/darwin-amd64/elvish-HEAD.tar.gz),
[Windows amd64](https://dl.elv.sh/windows-amd64/elvish-HEAD.zip), and
[many other platforms](https://elv.sh/get).

To build Elvish from source, you need

-   A supported OS: Linux, {Free,Net,Open}BSD, macOS, or Windows.

    **NOTE**: Windows support is experimental, and only Windows 10 is supported.

-   Go >= 1.15.

To build Elvish from source, follow these steps:

```sh
# 1. Start from any directory you want to store Elvish's source code
# 2. Clone the Git repository
git clone https://github.com/elves/elvish
# 3. Change into the repository
cd elvish
# 4. Build and install Elvish
make get
```

This will install Elvish to `~/go/bin`; you might want to add it to your `PATH`.

To install it elsewhere, override `GOBIN` in the `make` command:

```sh
make get GOBIN=$PWD # Install to the repo root (use $pwd if running in Elvish)
make get GOBIN=/usr/local/bin # Install to /usr/local/bin
```

Note that `GOBIN` must be an absolute path.

## Packaging Elvish

See [PACKAGING.md](PACKAGING.md) for notes for packagers.

## Contributing to Elvish

See [CONTRIBUTING.md](CONTRIBUTING.md) for notes for contributors.
