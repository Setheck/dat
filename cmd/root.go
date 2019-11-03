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

// parseEpochTime looks at args and attempts to determine the incoming epoch
func (r *RootCommand) parseEpochTime(args []string, def time.Time) (time.Time, error) {
	if len(args) > 0 {
		epoch, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("not a valid epoch")
		}
		*r.all = true
		return time.Unix(epoch, 0), nil
	}
	return def, nil
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
		Run: func(cmd *cobra.Command, args []string) {
			if *r.ver {
				fmt.Printf("%s - version:%s build:%s\n", Application, Version, Build)
				return
			}
			var err error
			if *r.paste {
				args, err = argsFromClipboard()
				if err != nil {
					fmt.Println(err)
				}
			}
			tm, err := r.parseEpochTime(args, time.Now())
			if err != nil {
				fmt.Println(err)
				return
			}
			output := r.BuildOutput(tm)
			if *r.copy {
				if err := clipboard.WriteAll(output); err != nil {
					fmt.Println("failed to copy to clipboard:", err)
				}
			}
			fmt.Print(output)
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
		fmt.Println(err)
		os.Exit(1)
	}
}

// argsFromClipboard is a helper to overwrite the args with the clipboard
// falls back to given args,
func argsFromClipboard() ([]string, error) {
	if in, err := clipboard.ReadAll(); err != nil {
		return nil, fmt.Errorf("reading from clipboard failed: %v", err)
	} else {
		in := strings.TrimSpace(in)
		return []string{in}, nil
	}
}
