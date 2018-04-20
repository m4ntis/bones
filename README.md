# BoNES
BoNES is a NES emulator library for go, as well as a cli for NES related tool.
For in depth documentation on the cli, see this [README](bones/README.md).

## Installation
BoNES is written in go, and requires golang to be installed on your computer.
I recommend installing go and compiling as described below (really quick and
easy), but if you are in a hurry or just can't be bothered, you can get the
already compiled binaries from the latest release in the
[releases page](https://github.com/m4ntis/bones/releases).

For information about installing go, you can visit
[Golang's download page](https://golang.org/dl).

After installing the lastest version of go, you can execute the following
command to get and compile from source:

```sh
go get -u github.com/m4ntis/bones/bones
```

Make sure that you add Go's bin directory to your `PATH` environment variable
as explained the language's installation instructions.

## Build Status
[![Build Status](https://travis-ci.org/m4ntis/bones.svg?branch=master)](https://travis-ci.org/m4ntis/bones)
[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/m4ntis/bones)
