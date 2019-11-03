package cmd

import (
	"fmt"
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
	rc.cmd.Run(nil, []string{"-v"})
	// Output: dat - version:v0.0.0 build:2019-11-02T01:23:46-0700
}

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
		{"no flags", tm, false, false, false, tmstr},
		{"utc", tm, false, false, true, tm.UTC().Format(DateFormat)},
		{"local", tm, false, true, false, tm.Local().Format(DateFormat)},
		{"local", tm, true, false, false,
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
	def := time.Now()
	tm := time.Unix(1572762509, 0)
	tmstr := strconv.FormatInt(tm.Unix(), 10)
	tests := []struct {
		name  string
		args  []string
		want  time.Time
		error bool
	}{
		{"fallback to default", []string{}, def, false},
		{"can't parse", []string{"qqqqqq"}, time.Time{}, true},
		{"parsed", []string{tmstr}, tm, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rc := RootCommand{
				ver:   &falsePtr,
				local: &falsePtr,
				utc:   &falsePtr,
				all:   &falsePtr,
			}
			got, err := rc.parseEpochTime(test.args, def)
			if test.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}
