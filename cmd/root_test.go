package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var truePtr = true
var falsePtr = false

func ExampleVersion() {
	rc := RootCommand{
		ver: &truePtr,
	}
	rc.Init()
	if err := rc.cmd.RunE(nil, nil); err != nil {
		log.Fatal(err) // this shouldn't ever happen!
	}
	// Output: dat - version:v0.0.0 build:2019-11-02T01:23:46-0700
}

// TODO: in progress
//func TestRootCommand(t *testing.T) {
//	clipboard = new(mocks.Clipper)
//	tests := []struct {
//		name  string
//		args  []string
//		copy  bool
//		paste bool
//		want  string
//	}{
//		{"", []string{""}, false, false, ""},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			rc := RootCommand{
//				ver: &falsePtr,
//			}
//			rc.Init()
//			if err := rc.cmd.RunE(nil, nil); err != nil {
//				t.Fatal(err)
//			}
//		})
//	}
//}

func TestRootCommand_BuildOutput(t *testing.T) {
	tm := time.Now()
	tmstr := strconv.FormatInt(tm.Unix(), 10)
	tests := []struct {
		name  string
		time  time.Time
		all   bool
		local bool
		utc   bool
		want  string
	}{
		{"no flags", tm, false, false, false, fmt.Sprintln(tmstr)},
		{"utc", tm, false, false, true, fmt.Sprintln(tm.UTC().Format(DateFormat))},
		{"local", tm, false, true, false, fmt.Sprintln(tm.Local().Format(DateFormat))},
		{"utc and local", tm, false, true, true,
			fmt.Sprintf("local: %s\n  utc: %s\n",
				tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
		{"all", tm, true, false, false,
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n",
				tm.Unix(), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rc := RootCommand{
				ver:   &falsePtr,
				local: &test.local,
				utc:   &test.utc,
				all:   &test.all,
			}
			got := rc.BuildOutput(test.time)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestParseEpochTime(t *testing.T) {
	tm := time.Unix(1572762509, 0)
	tmstr := strconv.FormatInt(tm.Unix(), 10)
	tests := []struct {
		name  string
		str   string
		want  time.Time
		error bool
	}{
		{"can't parse", "qqqqqq", time.Time{}, true},
		{"parsed", tmstr, tm, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseEpochTime(test.str)
			if test.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name string
		size int
		want string
	}{
		{"", 0, ""},
		{"", -1, ""},
		{"no trunc", 10, "no trunc"},
		{"ab", 1, "a"},
		{"abc", 2, "ab"},
		{"abc", 3, "abc"},
		{"abcd", 3, "..."},
		{"happy trees", 8, "happy..."},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TruncateString(test.name, test.size)
			assert.Equal(t, test.want, got)
		})
	}
}
