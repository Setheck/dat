package main

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/pflag"

	"github.com/Setheck/dat/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	rc := NewRootCommand()
	assert.NotNil(t, rc)
	assert.NotNil(t, rc.cmd)

}

func TestRootCommand_ParseFlags(t *testing.T) {
	fset := &pflag.FlagSet{}
	mockCommand := new(mocks.CobraCommand)
	mockCommand.On("Flags").Return(fset)
	rc := &RootCommand{cmd: mockCommand}
	rc.ParseFlags()
	mockCommand.AssertExpectations(t)

	// version
	assert.NotNil(t, fset.ShorthandLookup("v"))
	assert.NotNil(t, fset.Lookup("version"))

	// all
	assert.NotNil(t, fset.ShorthandLookup("a"))
	assert.NotNil(t, fset.Lookup("all"))

	// local
	assert.NotNil(t, fset.ShorthandLookup("l"))
	assert.NotNil(t, fset.Lookup("local"))

	// utc
	assert.NotNil(t, fset.ShorthandLookup("u"))
	assert.NotNil(t, fset.Lookup("utc"))

	// copy
	assert.NotNil(t, fset.ShorthandLookup("c"))
	assert.NotNil(t, fset.Lookup("copy"))

	// paste
	assert.NotNil(t, fset.ShorthandLookup("p"))
	assert.NotNil(t, fset.Lookup("paste"))
}

func TestRootCommand_Options(t *testing.T) {
	var truePtr = true
	var falsePtr = false

	tests := []struct {
		name string
		rc   *RootCommand
		want Options
	}{
		{"version flag",
			&RootCommand{
				ver:          &truePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
			},
			Options{Version: truePtr}},
		{"copy flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &truePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
			},
			Options{Copy: truePtr}},
		{"paste flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &truePtr,
				milliseconds: &falsePtr,
			},
			Options{Paste: truePtr}},
		{"all flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &truePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
			},
			Options{All: truePtr}},
		{"local flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &truePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
			},
			Options{Local: truePtr}},
		{"utc flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &truePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
			},
			Options{UTC: truePtr}},
		{"m flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &truePtr,
			},
			Options{Milliseconds: truePtr}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.rc.Options()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestRootCommand_Execute(t *testing.T) {
	cmdMock := new(mocks.CobraCommand)
	cmdMock.On("Execute").Return(nil)

	rc := &RootCommand{cmd: cmdMock}
	err := rc.Execute()
	assert.Nil(t, err)
	cmdMock.AssertExpectations(t)
}

func TestRunInit(t *testing.T) {
	saveClipboard := ClipboardHelper
	saveStdOut := StdOut
	saveBuildOutput := buildOutput
	saveTimeNow := timeNow
	defer func() {
		ClipboardHelper = saveClipboard
		StdOut = saveStdOut
		buildOutput = saveBuildOutput
		timeNow = saveTimeNow
	}()
	testOutput := "fake output data"
	buildOutput = func(tm time.Time, opts Options) string {
		return testOutput
	}
	testTime := time.Unix(1601167426, 0)
	timeNow = func() time.Time {
		return testTime
	}

	goodEpoch := "1601167426"

	tests := []struct {
		name         string
		args         []string
		options      Options
		want         string
		clipboardErr error
		epochErr     error
	}{
		{"version", nil, Options{Version: true}, "dat - version:v0.0.0 build:2019-11-02T01:23:46-0700\n", nil, nil},
		{"no args", nil, Options{}, testOutput, nil, nil},
		{"with input", []string{goodEpoch}, Options{}, testOutput, nil, nil},
		{"millisecond input", nil, Options{Milliseconds: true}, testOutput, nil, nil},
		{"input bad epoch", []string{"asdf"}, Options{}, testOutput, nil, assert.AnError},
		{"read from clipboard", nil, Options{Paste: true}, testOutput, nil, nil},
		{"read from clipboard error", nil, Options{Paste: true}, testOutput, assert.AnError, nil},
		{"copy to clipboard", nil, Options{Copy: true}, testOutput, nil, nil},
		{"copy to clipboard error", nil, Options{Copy: true}, testOutput, assert.AnError, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputBuffer := new(bytes.Buffer)
			StdOut = outputBuffer

			mockClipper := new(mocks.Clipper)
			mockClipper.On("ReadAll").Return(goodEpoch, test.clipboardErr)
			mockClipper.On("WriteAll", testOutput).Return(test.clipboardErr)
			ClipboardHelper = mockClipper

			err := RunE(test.options, test.args)
			if test.clipboardErr != nil || test.epochErr != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.want, outputBuffer.String())
			}
		})
	}
}

func TestRootCommand_BuildOutput(t *testing.T) {
	tm := time.Now()
	tmstr := strconv.FormatInt(tm.Unix(), 10)
	tmstr_ms := strconv.FormatInt(tm.UnixNano()/int64(time.Millisecond), 10)
	tests := []struct {
		name         string
		time         time.Time
		all          bool
		local        bool
		utc          bool
		milliseconds bool
		want         string
	}{
		{"no flags", tm, false, false, false, false, fmt.Sprintln(tmstr)},
		{"milliseconds", tm, false, false, false, true, fmt.Sprintln(tmstr_ms)},
		{"utc", tm, false, false, true, false, fmt.Sprintln(tm.UTC().Format(DateFormat))},
		{"local", tm, false, true, false, false, fmt.Sprintln(tm.Local().Format(DateFormat))},
		{"utc and local", tm, false, true, true, false,
			fmt.Sprintf("local: %s\n  utc: %s\n",
				tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
		{"all", tm, true, false, false, false,
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n",
				tm.Unix(), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
		{"ms all", tm, true, false, false, true,
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n",
				tm.UnixNano()/int64(time.Millisecond), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts := Options{
				All:          test.all,
				Local:        test.local,
				UTC:          test.utc,
				Milliseconds: test.milliseconds,
			}
			got := BuildOutput(test.time, opts)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestParseEpochTime(t *testing.T) {
	timeEpoch := int64(1572762509)
	timeEpochMs := int64(1572762509000)
	tmstr := strconv.FormatInt(timeEpoch, 10)
	tmstrMs := strconv.FormatInt(timeEpochMs, 10)
	tests := []struct {
		name           string
		str            string
		isMilliseconds bool
		want           time.Time
		error          bool
	}{
		{"can't parse", "qqqqqq", false, time.Time{}, true},
		{"parsed", tmstr, false, time.Unix(timeEpoch, 0), false},
		{"parsed", tmstrMs, true, time.Unix(0, timeEpochMs*int64(time.Millisecond)), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseEpochTime(test.str, test.isMilliseconds)
			if test.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, test.want.Equal(got), "want:%d got:%d", test.want.UnixNano(), got.UnixNano())
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
