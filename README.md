# dat

![Build & Test](https://github.com/Setheck/dat/workflows/Build%20&%20Test/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/setheck/dat)](https://goreportcard.com/report/github.com/setheck/dat)

```
dat is a simple tool for converting epochs,
when called without arguments dat returns the current epoch.
Likewise, if an epoch is not given the current epoch is assumed.

Usage:
  dat [epoch] [flags]

Flags:
  -a, --all            display the epoch in all formats
  -c, --copy           add output to clipboard
  -h, --help           help for dat
  -l, --local          display the epoch in the local timezone
  -m, --milliseconds   epochs in milliseconds
  -p, --paste          read input from clipboard
  -u, --utc            display the epoch in utc
  -v, --version        print version and exit


```

# install
See [Releases Page](https://github.com/Setheck/dat/releases) for the latest release prebuilt binaries.

Build and install yourself!
requires [golang](https://golang.org/doc/install)
```
git clone git@github.com:Setheck/dat.git
make install
```
 