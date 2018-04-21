# BoNES
BoNES is a NES emulation library for go, as well as a cli for NES related
utilities. For in depth documentation on the cli, run `bones -h` for general
usage or `bones [command] -h` for detailed information.

## Installation
BoNES is written in go, and requires golang to be installed on your computer.
I recommend installing go and compiling as described below (really quick and
easy), but if you are in a hurry or just can't be bothered, you can get the
pre-compiled binaries of the latest release in the
[releases page](https://github.com/m4ntis/bones/releases).

NOTE: If you can't find the binaries for your platform, you will need to
compile it yourself.

### Build deps
Building bones on a debian based distro requires the following packages to be
installed:
- `libgl1-mesa-dev`
- `xorg-dev`

### Building and installing

For information about installing go, you can visit
[Golang's download page](https://golang.org/dl).

After installing the lastest version of go, you can execute the following
command to get and compile from source:

```sh
go get -u github.com/m4ntis/bones/bones
```

Make sure that you add Go's bin directory to your `PATH` environment variable
as explained the language's installation instructions.

## Caveats
BoNES still currently implements only basic hardware features and basic
rendering functionality, meaning that harder to emulate games such as the ones
listed [here](https://wiki.nesdev.com/w/index.php/Tricky-to-emulate_games).
BoNES has been currently tested and found to work propperly only with a Donkey
Kong rom.

## Build Status
[![Build Status](https://travis-ci.org/m4ntis/bones.svg?branch=master)](https://travis-ci.org/m4ntis/bones)
[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/m4ntis/bones)
