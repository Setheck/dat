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
)

var (
	Application = "dat"
	Version     = "v0.0.0"
	Built       = "2019-11-02T01:23:46-0700"
)

const DateFormat = "01/02/2006 15:04:05 -0700"

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
}

// NewRootCommand creates a new instance of a RootCommand
func NewRootCommand() *RootCommand {
	rc := &RootCommand{}
	rc.cmd = &cobra.Command{
		Use: fmt.Sprint(Application, " [epoch]"),
		Long: fmt.Sprint(Application, ` is a simple tool for converting epochs,
when called without arguments dat returns the current epoch.
Likewise, if an epoch is not given the current epoch is assumed.`),
		SilenceUsage: true, // prevent usage on error
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunE(rc.Options(), args)
		},
	}
	return rc
}

// ParseFlags parse and assign flags
func (r *RootCommand) ParseFlags() {
	flgs := r.cmd.Flags()
	r.ver = flgs.BoolP("version", "v", false, "print version and exit")
	r.all = flgs.BoolP("all", "a", false, "display the epoch in all formats")
	r.local = flgs.BoolP("local", "l", false, "display the epoch in the local timezone")
	r.utc = flgs.BoolP("utc", "u", false, "display the epoch in utc")
	r.copy = flgs.BoolP("copy", "c", false, "add output to clipboard")
	r.paste = flgs.BoolP("paste", "p", false, "read input from clipboard")
	r.milliseconds = flgs.BoolP("milliseconds", "m", false, "epochs in milliseconds")
	r.format = flgs.StringP("format", "f", "", "https://golang.org/pkg/time/ format for time output including constant names.")
}

// Options retrieves command input options
func (r *RootCommand) Options() Options {
	return Options{
		Version:      *r.ver,
		Copy:         *r.copy,
		Paste:        *r.paste,
		All:          *r.all,
		Local:        *r.local,
		UTC:          *r.utc,
		Milliseconds: *r.milliseconds,
		Format:       *r.format,
	}
}

// Execute run the command
func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}

// test points
var StdOut io.Writer = os.Stdout
var buildOutput = BuildOutput
var timeNow = time.Now

var Banner = strings.ReplaceAll(`      _       _   
     | |     | |  
   __| | __ _| |_ 
  / _q |/ _q | __|
 | (_| | (_| | |_
  \__,_|\__,_|\__|`, "q", "`")

// RunE is the command run function
func RunE(opts Options, args []string) error {
	if opts.Version {
		fmt.Fprintln(StdOut, Banner)
		fmt.Fprintln(StdOut, "app:    ", Application)
		fmt.Fprintln(StdOut, "version:", Version)
		fmt.Fprintln(StdOut, "built:  ", Built)
		return nil
	}

	// default to now
	var epocInt int64
	if opts.Milliseconds {
		epocInt = timeNow().UnixNano() / int64(time.Millisecond)
	} else {
		epocInt = timeNow().Unix()
	}
	epochstr := strconv.FormatInt(epocInt, 10)

	// take value passed in
	if len(args) > 0 {
		epochstr = args[0]
	}

	// paste mode reads from the clipboard
	if opts.Paste {
		var err error
		epochstr, err = ClipboardHelper.ReadAll()
		if err != nil {
			return err
		}
	}

	// validate and convert to time
	tm, err := ParseEpochTime(epochstr, opts.Milliseconds)
	if err != nil {
		return err
	}

	output := buildOutput(tm, opts)
	if opts.Copy {
		if err := ClipboardHelper.WriteAll(output); err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(StdOut, output)
	return err
}

// Options
type Options struct {
	Version      bool
	Copy         bool
	Paste        bool
	All          bool
	Local        bool
	UTC          bool
	Milliseconds bool
	Format       string
}

// BuildOutput returns the output of the time for the given options
func BuildOutput(tm time.Time, opts Options) string {
	output := ""

	var intTime int64
	if opts.Milliseconds {
		intTime = tm.UnixNano() / int64(time.Millisecond)
	} else {
		intTime = tm.Unix()
	}

	outFormat := DateFormat
	if opts.Format != "" {
		outFormat = opts.Format
	}

	switch {
	case opts.All:
		output += fmt.Sprintln("epoch:", intTime)
		fallthrough
	case opts.Local && opts.UTC:
		output += fmt.Sprintln("local:", FormatOutput(tm.Local(), outFormat))
		output += fmt.Sprintln("  utc:", FormatOutput(tm.UTC(), outFormat))
	default:
		out := strconv.FormatInt(intTime, 10)
		if opts.Local {
			out = FormatOutput(tm.Local(), outFormat)
		} else if opts.UTC {
			out = FormatOutput(tm.UTC(), outFormat)
		}
		output = fmt.Sprintln(out)
	}
	return output
}

func FormatOutput(tm time.Time, outFmtS string) string {
	outFmt := outFmtS
	switch strings.ToLower(outFmtS) {
	case "ansic":
		outFmt = time.ANSIC
	case "unixdate":
		outFmt = time.UnixDate
	case "rubydate":
		outFmt = time.RubyDate
	case "rfc822":
		outFmt = time.RFC822
	case "rfc822z":
		outFmt = time.RFC822Z
	case "rfc850":
		outFmt = time.RFC850
	case "rfc1123":
		outFmt = time.RFC1123
	case "rfc1123z":
		outFmt = time.RFC1123Z
	case "rfc3339":
		outFmt = time.RFC3339
	case "rfc3339nano":
		outFmt = time.RFC3339Nano
	case "kitchen":
		outFmt = time.Kitchen
	case "stamp":
		outFmt = time.Stamp
	case "stampmilli":
		outFmt = time.StampMilli
	case "stampmicro":
		outFmt = time.StampMicro
	case "stampnano":
		outFmt = time.StampNano
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
