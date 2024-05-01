# dat

![Build & Test](https://github.com/Setheck/dat/workflows/Build%20&%20Test/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/setheck/dat)](https://goreportcard.com/report/github.com/setheck/dat)

```
dat is a simple tool for converting epochs,
when called without arguments dat returns the current epoch.
Likewise, if an epoch is not given the current epoch is assumed.

Usage:
  dat [epoch] [flags]

Flags:
  -a, --all             display the epoch and formatted local and utc values of the epoch
  -c, --copy            copy output to the clipboard
  -d, --delta string    a duration in which to modify the epoch (ex:+2h3s) see https://golang.org/pkg/time/#ParseDuration
  -f, --format string   https://golang.org/pkg/time/ format for time output including constant names
  -h, --help            help for dat
  -l, --local           display the formatted epoch in the local timezone
  -m, --milliseconds    epochs in milliseconds
  -p, --paste           read input from the clipboard
  -t, --tf              attempt to parse input as a known time format
  -u, --utc             display the formatted epoch in the utc timezone
  -v, --version         print version and exit
  -z, --zone string     display a specific time zone by tz database name see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
```

# install
See [Releases Page](https://github.com/Setheck/dat/releases) for the latest release prebuilt binaries.

Build and install yourself!
requires [golang](https://golang.org/doc/install)
```bash
git clone git@github.com:Setheck/dat.git
cd dat && make install
```
 
