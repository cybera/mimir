# README

## Setup

Make sure you have Go 1.11+ installed and you'll probably want `$GOPATH/bin` in your `$PATH` for convenience. Note that you should clone the repository *outside* of your `$GOPATH` to enable module support.

```bash
$ go get -u github.com/gobuffalo/packr/v2/packr2
$ git clone https://github.com/cybera/ccds
```

That's it, the first time you run or build it, all dependencies will be installed automatically.

## Usage

You can either run it directly:

```bash
$ go run main.go
```

Or build it first:

```bash
packr2 build
./ccds
```

Note the use of `packr2` rather than `go build`. Packr is a tool to bundle static assets in Go binaries and it wraps the usual build and install commands.