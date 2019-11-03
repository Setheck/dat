package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	Application = "dat"
	Version     = "v0.0.0"
	Build       = "2019-11-02T01:23:46-0700"
)

const DateFormat = "01/02/2006 15:04:05 -0700"

type RootCommand struct {
	cmd *cobra.Command

	ver   *bool
	local *bool
	utc   *bool
	all   *bool
	copy  *bool
	paste *bool
}

func (r *RootCommand) ParseFlags() {
	flgs := r.cmd.Flags()
	r.ver = flgs.BoolP("version", "v", false, "print version and exit")
	r.all = flgs.BoolP("all", "a", false, "display the epoch in all formats")
	r.local = flgs.BoolP("local", "l", false, "display the epoch in the local timezone")
	r.utc = flgs.BoolP("utc", "u", false, "display the epoch in utc")
	r.copy = flgs.BoolP("copy", "c", false, "add output to clipboard")
	r.paste = flgs.BoolP("paste", "p", false, "read input from clipboard")
}

func (r *RootCommand) BuildOutput(tm time.Time) string {
	output := ""
	switch {
	case *r.all:
		output += fmt.Sprintln("epoch:", tm.Unix())
		fallthrough
	case *r.local && *r.utc:
		output += fmt.Sprintln("local:", tm.Local().Format(DateFormat))
		output += fmt.Sprintln("  utc:", tm.UTC().Format(DateFormat))
	default:
		out := strconv.FormatInt(tm.Unix(), 10)
		if *r.local {
			out = tm.Local().Format(DateFormat)
		} else if *r.utc {
			out = tm.UTC().Format(DateFormat)
		}
		output = out
	}
	return output
}

func (r *RootCommand) Init() {
	r.cmd = &cobra.Command{
		Use: fmt.Sprint(Application, " [epoch]"),
		Long: fmt.Sprint(Application, ` is a simple tool for converting epochs,
when called without arguments dat returns the current epoch.
Likewise, if an epoch is not given the current epoch is assumed.
If given an epoch, all formats (epoch, local, utc) will be output.`),
		SilenceUsage: true, // prevent usage on error
		RunE: func(cmd *cobra.Command, args []string) error {
			if *r.ver {
				fmt.Printf("%s - version:%s build:%s\n", Application, Version, Build)
				return nil
			}

			// default to now
			epochstr := strconv.FormatInt(time.Now().Unix(), 10)
			if len(args) > 0 {
				epochstr = args[0]
			}

			var err error
			if *r.paste {
				epochstr, err = ReadFromClipboard()
				if err != nil {
					return err
				}
			}

			tm, err := ParseEpochTime(epochstr)
			if err != nil {
				return err
			}

			output := r.BuildOutput(tm)
			if *r.copy {
				if err := WriteToClipboard(output); err != nil {
					return err
				}
			}

			fmt.Print(output)
			return nil
		},
	}
}

func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}

func Execute() {
	rootCmd := RootCommand{}
	rootCmd.Init()
	rootCmd.ParseFlags()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ParseEpochTime tries to parse the string as an int, then converts to a time.Time
func ParseEpochTime(str string) (time.Time, error) {
	if epoch, err := strconv.ParseInt(str, 10, 64); err != nil {
		return time.Time{}, fmt.Errorf("%q is not a valid epoch", TruncateString(str, 20))
	} else {
		return time.Unix(epoch, 0), nil
	}
}

// ReadFromClipboard is a helper to overwrite the args with the clipboard
// falls back to given args,
func ReadFromClipboard() (string, error) {
	if in, err := clipboard.ReadAll(); err != nil {
		return "", fmt.Errorf("reading from clipboard failed: %v", err)
	} else {
		in := strings.TrimSpace(in)
		return in, nil
	}
}

func WriteToClipboard(s string) error {
	if err := clipboard.WriteAll(s); err != nil {
		fmt.Println("failed to write to clipboard:", err)
	}
	return nil
}

func TruncateString(str string, size int) string {
	if len(str) > size {
		if size > 3 {
			size -= 3
		}
		return str[0:size] + "..."
	}
	return str
}
