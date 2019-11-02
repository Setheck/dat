package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.0.0"
	Build   = "2019-11-02T01:23:46-0700"
)

const format = "01/02/2006 15:04:05 -0700"

type RootCommand struct {
	cmd *cobra.Command

	ver   *bool
	local *bool
	utc   *bool
	all   *bool
}

func (r *RootCommand) ParseFlags() {
	flgs := r.cmd.Flags()
	r.ver = flgs.BoolP("version", "v", false, "print version and exit")
	r.all = flgs.BoolP("all", "a", false, "display the epoch in all formats")
	r.local = flgs.BoolP("local", "l", false, "display the epoch in the local timezone")
	r.utc = flgs.BoolP("utc", "u", false, "display the epoch in utc")
}

func (r *RootCommand) InitCmd() {
	r.cmd = &cobra.Command{
		Use: "dat [epoch]",
		Long: `dat is a simple tool for converting epochs,
when called without arguments dat returns the current epoch.
Likewise, if an epoch is not given the current epoch is assumed.
If given an epoch, all formats (epoch, local, utc) will be output.`,
		Run: func(cmd *cobra.Command, args []string) {
			if *r.ver {
				fmt.Printf("dat - version:%s  build:%s\n", Version, Build)
				return
			}
			tm := time.Now()
			if len(args) > 0 {
				epoch, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					fmt.Println("not a valid epoch")
					return
				}
				if epoch > 9999999999 {
					sec := epoch / 1000
					nsec := epoch % 1000
					tm = time.Unix(sec, nsec)
				} else {
					tm = time.Unix(epoch, 0)
				}
				*r.all = true
			}
			switch {
			case *r.all:
				fmt.Println("epoch:", tm.Unix())
				fallthrough
			case *r.local && *r.utc:
				fmt.Println("local:", tm.Local().Format(format))
				fmt.Println("  utc:", tm.UTC().Format(format))
			default:
				out := strconv.FormatInt(tm.Unix(), 10)
				if *r.local {
					out = tm.Local().Format(format)
				} else if *r.utc {
					out = tm.UTC().Format(format)
				}
				fmt.Println(out)
			}
		},
	}
}

func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}

func Execute() {
	rootCmd := RootCommand{}
	rootCmd.InitCmd()
	rootCmd.ParseFlags()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
