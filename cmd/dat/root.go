package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/Setheck/dat/pkg/build"
	"github.com/Setheck/dat/pkg/clipper"
)

const DateFormat = "01/02/2006 15:04:05 -0700"

var supportedTimeFormats = map[string]string{
	// format: name
	time.ANSIC:       "ANSIC",
	time.UnixDate:    "UnixDate",
	time.RubyDate:    "RubyDate",
	time.RFC822:      "RFC822",
	time.RFC822Z:     "RFC822Z",
	time.RFC850:      "RFC850",
	time.RFC1123:     "RFC1123",
	time.RFC1123Z:    "RFC1123Z",
	time.RFC3339:     "RFC3339",
	time.RFC3339Nano: "RFC3339Nano",
	time.Kitchen:     "Kitchen",
	time.Stamp:       "Stamp",
	time.StampMilli:  "StampMilli",
	time.StampMicro:  "StampMicro",
	time.StampNano:   "StampNano",
}

// CobraCommand interface for *cobra.Command
type CobraCommand interface {
	Execute() error
	Flags() *pflag.FlagSet
}

var _ CobraCommand = &cobra.Command{}

// RootCommand root cobra command
type RootCommand struct {
	cmd CobraCommand

	ver          *bool
	local        *bool
	utc          *bool
	all          *bool
	copy         *bool
	paste        *bool
	milliseconds *bool
	format       *string
	delta        *string
	zone         *string
	tf           *bool
}

// options
type options struct {
	Version      bool
	Copy         bool
	Paste        bool
	All          bool
	Local        bool
	UTC          bool
	Milliseconds bool
	Format       string
	Delta        string
	Zone         string
	Tf           bool

	detectedFormat string
}

// NewRootCommand creates a new instance of a RootCommand
func NewRootCommand() *RootCommand {
	rc := &RootCommand{}
	rc.cmd = &cobra.Command{
		Use: fmt.Sprint(build.Application, " [epoch]"),
		Long: fmt.Sprint(build.Application, ` is a simple tool for converting epochs,
when called without arguments dat returns the current epoch.
Likewise, if an epoch is not given the current epoch is assumed.`),
		SilenceUsage: true, // prevent usage on error
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunE(rc.options(), args)
		},
	}
	return rc
}

// ParseFlags parse and assign flags
func (r *RootCommand) ParseFlags() {
	flgs := r.cmd.Flags()
	r.ver = flgs.BoolP("version", "v", false, "print version and exit")
	r.all = flgs.BoolP("all", "a", false, "display the epoch and formatted local and utc values of the epoch")
	r.local = flgs.BoolP("local", "l", false, "display the formatted epoch in the local timezone")
	r.utc = flgs.BoolP("utc", "u", false, "display the formatted epoch in the utc timezone")
	r.copy = flgs.BoolP("copy", "c", false, "copy output to the clipboard")
	r.paste = flgs.BoolP("paste", "p", false, "read input from the clipboard")
	r.milliseconds = flgs.BoolP("milliseconds", "m", false, "epochs in milliseconds")
	r.format = flgs.StringP("format", "f", "", "https://golang.org/pkg/time/ format for time output including constant names")
	r.delta = flgs.StringP("delta", "d", "", "a duration in which to modify the epoch (ex:+2h3s) see https://golang.org/pkg/time/#ParseDuration")
	r.zone = flgs.StringP("zone", "z", "", "display a specific time zone by tz database name see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones")
	r.tf = flgs.BoolP("tf", "t", false, "attempt to parse input as a known time format")
}

// options retrieves command input options
func (r *RootCommand) options() options {
	return options{
		Version:      *r.ver,
		Copy:         *r.copy,
		Paste:        *r.paste,
		All:          *r.all,
		Local:        *r.local,
		UTC:          *r.utc,
		Milliseconds: *r.milliseconds,
		Format:       *r.format,
		Delta:        *r.delta,
		Zone:         *r.zone,
		Tf:           *r.tf,
	}
}

// Execute run the command
func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}

// test points
var stdOut io.Writer = os.Stdout
var buildOutput = BuildOutput

var timeNow = time.Now

var banner = strings.ReplaceAll(`      _       _   
     | |     | |  
   __| | __ _| |_ 
  / _q |/ _q | __|
 | (_| | (_| | |_
  \__,_|\__,_|\__|`, "q", "`")

// RunE is the command run function
func RunE(opts options, args []string) (err error) {
	if opts.Version {
		fmt.Fprintln(stdOut, banner)
		fmt.Fprintln(stdOut, "app:    ", build.Application)
		fmt.Fprintln(stdOut, "version:", build.Version)
		fmt.Fprintln(stdOut, "built:  ", build.Built)
		return nil
	}

	// default to now
	var epocInt int64
	epoch := timeNow()
	if opts.Milliseconds {
		epocInt = epoch.UnixMilli()
	} else {
		epocInt = epoch.Unix()
	}
	epochstr := strconv.FormatInt(epocInt, 10)

	// take value passed in
	if len(args) > 0 {
		epochstr = args[0]
	}

	// paste mode reads from the clipboard
	if opts.Paste {
		var err error
		epochstr, err = clipper.ClipboardHelper.ReadAll()
		if err != nil {
			return err
		}
	}

	var (
		tm     time.Time
		format string
	)
	if opts.Tf {
		tm, format, err = ParseTime(epochstr)
		if err != nil {
			return err
		}
		if name, ok := supportedTimeFormats[format]; ok {
			opts.detectedFormat = name
		} else {
			fmt.Println("failed to detect format.")
		}
	} else {
		// validate and convert to time
		tm, err = ParseEpochTime(epochstr, opts.Milliseconds)
		if err != nil {
			return err
		}
	}

	output := buildOutput(tm, opts)
	if opts.Copy {
		if err := clipper.ClipboardHelper.WriteAll(strings.TrimSpace(output)); err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(stdOut, output)
	return err
}

// BuildOutput returns the output of the time for the given options
func BuildOutput(tm time.Time, opts options) string {
	output := ""

	// add delta if applicable.
	if opts.Delta != "" {
		tm = AddDelta(tm, opts.Delta)
	}

	var intTime int64
	if opts.Milliseconds {
		intTime = tm.UnixMilli()
	} else {
		intTime = tm.Unix()
	}

	outFormat := DateFormat
	if opts.Format != "" {
		outFormat = opts.Format
	}

	var formattedZone string
	if opts.Zone != "" {
		loc, err := time.LoadLocation(opts.Zone)
		if err == nil {
			formattedZone = FormatOutput(tm.In(loc), outFormat)
		}
	}

	switch {
	case opts.All:
		if opts.detectedFormat != "" {
			output += fmt.Sprintln("detected:", opts.detectedFormat)
		}
		output += fmt.Sprintln("epoch:", intTime)
		fallthrough
	case opts.Local && opts.UTC:
		output += fmt.Sprintln("local:", FormatOutput(tm.Local(), outFormat))
		output += fmt.Sprintln("  utc:", FormatOutput(tm.UTC(), outFormat))
		if formattedZone != "" {
			output += fmt.Sprintln(" zone:", formattedZone)
		}

	case opts.Local && formattedZone != "":
		output += fmt.Sprintln("local:", FormatOutput(tm.Local(), outFormat))
		output += fmt.Sprintln(" zone:", formattedZone)

	case opts.UTC && formattedZone != "":
		output += fmt.Sprintln("  utc:", FormatOutput(tm.UTC(), outFormat))
		output += fmt.Sprintln(" zone:", formattedZone)

	default:
		out := strconv.FormatInt(intTime, 10)
		if opts.Local {
			out = FormatOutput(tm.Local(), outFormat)
		} else if opts.UTC {
			out = FormatOutput(tm.UTC(), outFormat)
		} else if opts.Format != "" {
			out = FormatOutput(tm, outFormat)
		} else if formattedZone != "" {
			out = formattedZone
		}
		output = fmt.Sprintln(out)
	}

	return output
}

// AddDelta simply adds the given duration if it is valid to the given time, ignores otherwise.
func AddDelta(tm time.Time, delta string) time.Time {
	dur, err := time.ParseDuration(delta)
	if err == nil {
		tm = tm.Add(dur)
	}
	return tm
}

// FormatOutput parses the provided time against the provided format string.
// replacing named constants with the expected format.
func FormatOutput(tm time.Time, outFmtS string) string {
	outFmt := outFmtS
	lowerFmt := strings.ToLower(outFmtS)
	for format, name := range supportedTimeFormats {
		if strings.ToLower(name) == lowerFmt {
			outFmt = format
			break
		}
	}
	formattedTime := tm.Format(outFmt)
	if formattedTime == outFmt {
		outFmt = DateFormat
	}

	return tm.Format(outFmt)
}

// ParseEpochTime tries to parse the string as an int, then converts to a time.Time
func ParseEpochTime(str string, milliseconds bool) (time.Time, error) {
	epoch, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("%q is not a valid epoch", TruncateString(str, 20))
	}

	if milliseconds {
		return time.Unix(0, epoch*int64(time.Millisecond)), nil
	}

	return time.Unix(epoch, 0), nil
}

func ParseTime(str string) (time.Time, string, error) {
	for format := range supportedTimeFormats {
		tm, err := time.Parse(format, str)
		if err == nil {
			return tm, format, nil
		}
	}
	return time.Unix(0, 0), "", fmt.Errorf("invalid time format")
}

// TruncateString reduces the size of str to the given size.
// if str is truncated, it will include a trailing ellipsis.
func TruncateString(str string, size int) string {
	if size < 0 {
		size = 0
	}
	length := len(str)
	postfix := ""
	if length > size {
		if length > 3 {
			postfix = "..."
			size -= 3
		}
		return str[0:size] + postfix
	}
	return str
}
