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
	Build       = "2019-11-02T01:23:46-0700"
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

	ver   *bool
	local *bool
	utc   *bool
	all   *bool
	copy  *bool
	paste *bool
}

// NewRootCommand creates a new instance of a RootCommand
func NewRootCommand() *RootCommand {
	rc := &RootCommand{}
	rc.cmd = &cobra.Command{
		Use: fmt.Sprint(Application, " [epoch]"),
		Long: fmt.Sprint(Application, ` is a simple tool for converting epochs,
when called without arguments dat returns the current epoch.
Likewise, if an epoch is not given the current epoch is assumed.
If given an epoch, all formats (epoch, local, utc) will be output.`),
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
}

// Options retrieves command input options
func (r *RootCommand) Options() Options {
	return Options{
		Version: *r.ver,
		Copy:    *r.copy,
		Paste:   *r.paste,
		All:     *r.all,
		Local:   *r.local,
		UTC:     *r.utc,
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

// RunE is the command run function
func RunE(opts Options, args []string) error {
	if opts.Version {
		_, err := fmt.Fprintf(StdOut, "%s - version:%s build:%s\n", Application, Version, Build)
		return err
	}

	// default to now
	epochstr := strconv.FormatInt(timeNow().Unix(), 10)

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
	tm, err := ParseEpochTime(epochstr)
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
	Version bool
	Copy    bool
	Paste   bool
	All     bool
	Local   bool
	UTC     bool
}

// BuildOutput returns the output of the time for the given options
func BuildOutput(tm time.Time, opts Options) string {
	output := ""
	switch {
	case opts.All:
		output += fmt.Sprintln("epoch:", tm.Unix())
		fallthrough
	case opts.Local && opts.UTC:
		output += fmt.Sprintln("local:", tm.Local().Format(DateFormat))
		output += fmt.Sprintln("  utc:", tm.UTC().Format(DateFormat))
	default:
		out := strconv.FormatInt(tm.Unix(), 10)
		if opts.Local {
			out = tm.Local().Format(DateFormat)
		} else if opts.UTC {
			out = tm.UTC().Format(DateFormat)
		}
		output = fmt.Sprintln(out)
	}
	return output
}

// ParseEpochTime tries to parse the string as an int, then converts to a time.Time
func ParseEpochTime(str string) (time.Time, error) {
	if epoch, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64); err != nil {
		return time.Time{}, fmt.Errorf("%q is not a valid epoch", TruncateString(str, 20))
	} else {
		return time.Unix(epoch, 0), nil
	}
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
