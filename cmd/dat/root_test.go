package main

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/Setheck/dat/pkg/clipper"
	"github.com/Setheck/dat/pkg/mocks"
)

const (
	tzLosAngeles = "America/Los_Angeles"
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

	// format
	assert.NotNil(t, fset.ShorthandLookup("f"))
	assert.NotNil(t, fset.Lookup("format"))

	// delta
	assert.NotNil(t, fset.ShorthandLookup("d"))
	assert.NotNil(t, fset.Lookup("delta"))

	// zone
	assert.NotNil(t, fset.ShorthandLookup("z"))
	assert.NotNil(t, fset.Lookup("zone"))
}

func TestRootCommand_Options(t *testing.T) {
	var truePtr = true
	var falsePtr = false

	tests := []struct {
		name string
		rc   *RootCommand
		want options
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
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{Version: truePtr}},
		{"copy flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &truePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{Copy: truePtr}},
		{"paste flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &truePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{Paste: truePtr}},
		{"all flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &truePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{All: truePtr}},
		{"local flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &truePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{Local: truePtr}},
		{"utc flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &truePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{UTC: truePtr}},
		{"m flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &truePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{Milliseconds: truePtr}},
		{"f flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, time.RFC3339),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{Format: time.RFC3339}},
		{"d flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, "360h10m"),
				zone:         StfPtr(t, ""),
				tf:           &falsePtr,
			},
			options{Delta: "360h10m"}},
		{"z flag",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, "America/Los_Angeles"),
				tf:           &falsePtr,
			},
			options{Zone: "America/Los_Angeles"}},
		{"detect format (tf)",
			&RootCommand{
				ver:          &falsePtr,
				local:        &falsePtr,
				utc:          &falsePtr,
				all:          &falsePtr,
				copy:         &falsePtr,
				paste:        &falsePtr,
				milliseconds: &falsePtr,
				format:       StfPtr(t, ""),
				delta:        StfPtr(t, ""),
				zone:         StfPtr(t, ""),
				tf:           &truePtr,
			},
			options{Tf: true}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.rc.options()
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
	saveClipboard := clipper.ClipboardHelper
	saveStdOut := stdOut
	saveBuildOutput := buildOutput
	saveTimeNow := timeNow
	defer func() {
		clipper.ClipboardHelper = saveClipboard
		stdOut = saveStdOut
		buildOutput = saveBuildOutput
		timeNow = saveTimeNow
	}()
	testOutput := "fake output data"
	buildOutput = func(tm time.Time, opts options) string {
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
		options      options
		want         string
		clipboardErr error
		epochErr     error
	}{
		{"version", nil, options{Version: true}, fmt.Sprintf("%s\napp:     dat\nversion: v0.0.0\nbuilt:   2019-11-02T01:23:46-0700\n", banner), nil, nil},
		{"no args", nil, options{}, testOutput, nil, nil},
		{"with input", []string{goodEpoch}, options{}, testOutput, nil, nil},
		{"millisecond input", nil, options{Milliseconds: true}, testOutput, nil, nil},
		{"input bad epoch", []string{"asdf"}, options{}, testOutput, nil, assert.AnError},
		{"read from clipboard", nil, options{Paste: true}, testOutput, nil, nil},
		{"read from clipboard error", nil, options{Paste: true}, testOutput, assert.AnError, nil},
		{"copy to clipboard", nil, options{Copy: true}, testOutput, nil, nil},
		{"copy to clipboard error", nil, options{Copy: true}, testOutput, assert.AnError, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputBuffer := new(bytes.Buffer)
			stdOut = outputBuffer

			mockClipper := new(mocks.Clipper)
			mockClipper.On("ReadAll").Return(goodEpoch, test.clipboardErr)
			mockClipper.On("WriteAll", testOutput).Return(test.clipboardErr)
			clipper.ClipboardHelper = mockClipper

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
	tmStr := strconv.FormatInt(tm.Unix(), 10)
	tmStrPlus100h := strconv.FormatInt(tm.Add(100*time.Hour).Unix(), 10)
	tmStrMinus100h := strconv.FormatInt(tm.Add(-100*time.Hour).Unix(), 10)
	tmStrMillis := strconv.FormatInt(tm.UnixNano()/int64(time.Millisecond), 10)
	laZone, err := time.LoadLocation(tzLosAngeles)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		time    time.Time
		options options
		want    string
	}{
		{"no flags", tm, options{},
			fmt.Sprintln(tmStr)},
		{"delta +100h", tm, options{Delta: "100h"},
			fmt.Sprintln(tmStrPlus100h)},
		{"delta -100h", tm, options{Delta: "-100h"},
			fmt.Sprintln(tmStrMinus100h)},
		{"milliseconds", tm, options{Milliseconds: true},
			fmt.Sprintln(tmStrMillis)},
		{"format", tm, options{Format: "rfc3339"},
			fmt.Sprintln(tm.Format(time.RFC3339))},
		{"zone", tm, options{Zone: tzLosAngeles},
			fmt.Sprintln(tm.In(laZone).Format(DateFormat))},
		{"utc", tm, options{UTC: true},
			fmt.Sprintln(tm.UTC().Format(DateFormat))},
		{"utc and zone", tm, options{UTC: true, Zone: tzLosAngeles},
			fmt.Sprintf("  utc: %s\n zone: %s\n", tm.UTC().Format(DateFormat), tm.In(laZone).Format(DateFormat))},
		{"local", tm, options{Local: true},
			fmt.Sprintln(tm.Local().Format(DateFormat))},
		{"local and zone", tm, options{Local: true, Zone: tzLosAngeles},
			fmt.Sprintf("local: %s\n zone: %s\n", tm.Local().Format(DateFormat), tm.In(laZone).Format(DateFormat))},
		{"utc and local", tm, options{Local: true, UTC: true},
			fmt.Sprintf("local: %s\n  utc: %s\n", tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
		{"all with detection", tm, options{All: true, detectedFormat: "rfc3339"},
			fmt.Sprintf("detected: %s\nepoch: %d\nlocal: %s\n  utc: %s\n", "rfc3339", tm.Unix(), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
		{"all", tm, options{All: true},
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n", tm.Unix(), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
		{"all with zone", tm, options{All: true, Zone: tzLosAngeles},
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n zone: %s\n", tm.Unix(), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat), tm.In(laZone).Format(DateFormat))},
		{"ms all", tm, options{Milliseconds: true, All: true},
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n", tm.UnixNano()/int64(time.Millisecond), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat))},
		{"ms all with zone", tm, options{Milliseconds: true, All: true, Zone: tzLosAngeles},
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n zone: %s\n", tm.UnixNano()/int64(time.Millisecond), tm.Local().Format(DateFormat), tm.UTC().Format(DateFormat), tm.In(laZone).Format(DateFormat))},
		{"ms all with format", tm, options{Milliseconds: true, All: true, Format: "rfc3339"},
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n", tm.UnixNano()/int64(time.Millisecond), tm.Local().Format(time.RFC3339), tm.UTC().Format(time.RFC3339))},
		{"ms all with format and zone", tm, options{Milliseconds: true, All: true, Zone: tzLosAngeles, Format: "rfc3339"},
			fmt.Sprintf("epoch: %d\nlocal: %s\n  utc: %s\n zone: %s\n", tm.UnixNano()/int64(time.Millisecond), tm.Local().Format(time.RFC3339), tm.UTC().Format(time.RFC3339), tm.In(laZone).Format(time.RFC3339))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := BuildOutput(test.time, test.options)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestParseEpochTime(t *testing.T) {
	timeEpoch := int64(1572762509)
	timeEpochMillis := int64(1572762509000)
	tmStr := strconv.FormatInt(timeEpoch, 10)
	tmStrMillis := strconv.FormatInt(timeEpochMillis, 10)
	tests := []struct {
		name           string
		str            string
		isMilliseconds bool
		want           time.Time
		error          bool
	}{
		{"can't parse", "qqqqqq", false, time.Time{}, true},
		{"parsed", tmStr, false, time.Unix(timeEpoch, 0), false},
		{"parsed", tmStrMillis, true, time.Unix(0, timeEpochMillis*int64(time.Millisecond)), false},
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

func StfPtr(t *testing.T, s string) *string {
	t.Helper()
	return &s
}

func TestFormatOutput(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name   string
		time   time.Time
		format string
		want   string
	}{
		{"ANSIC", testTime, "ansic", testTime.Format(time.ANSIC)},
		{"UnixDate", testTime, "UnixDate", testTime.Format(time.UnixDate)},
		{"RubyDate", testTime, "RubyDate", testTime.Format(time.RubyDate)},
		{"RFC822", testTime, "RFC822", testTime.Format(time.RFC822)},
		{"RFC822Z", testTime, "RFC822Z", testTime.Format(time.RFC822Z)},
		{"RFC850", testTime, "RFC850", testTime.Format(time.RFC850)},
		{"RFC1123", testTime, "RFC1123", testTime.Format(time.RFC1123)},
		{"RFC1123Z", testTime, "RFC1123Z", testTime.Format(time.RFC1123Z)},
		{"RFC3339", testTime, "RFC3339", testTime.Format(time.RFC3339)},
		{"RFC3339Nano", testTime, "RFC3339Nano", testTime.Format(time.RFC3339Nano)},
		{"Kitchen", testTime, "Kitchen", testTime.Format(time.Kitchen)},
		{"Stamp", testTime, "Stamp", testTime.Format(time.Stamp)},
		{"StampMilli", testTime, "StampMilli", testTime.Format(time.StampMilli)},
		{"StampMicro", testTime, "StampMicro", testTime.Format(time.StampMicro)},
		{"StampNano", testTime, "StampNano", testTime.Format(time.StampNano)},
		{"other format", testTime, "Jan 15:05:04 MST -700 2006", testTime.Format("Jan 15:05:04 MST -700 2006")},
		{"unknown format", testTime, "NoTaGoOdFoRmAt", testTime.Format(DateFormat)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := FormatOutput(test.time, test.format)
			assert.Equal(t, got, test.want)
		})
	}
}

func TestAddDelta(t *testing.T) {
	tests := []struct {
		name         string
		time         time.Time
		delta        string
		expectedTime time.Time
	}{
		{"no delta, no addition", time.Unix(1625211007, 0), "", time.Unix(1625211007, 0)},
		{"adding 300h2m", time.Unix(1625211007, 0), "300h2m", time.Unix(1626291127, 0)},
		{"subtracting 300h2m", time.Unix(1625211007, 0), "-300h2m", time.Unix(1624130887, 0)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := AddDelta(test.time, test.delta)
			assert.Equal(t, test.expectedTime, got)
		})
	}
}
