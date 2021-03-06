// SPDX-License-Identifier: MIT
//
// Copyright © 2019 Kent Gibson <warthog618@gmail.com>.

// +build linux

// A clone of libgpiod gpioinfo.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:                   "info [flags] [chip]...",
	Short:                 "Info about chip lines",
	Long:                  `Print information about all lines of the specified GPIO chip(s) (or all gpiochips if none are specified).`,
	Run:                   info,
	DisableFlagsInUseLine: true,
}

func info(cmd *cobra.Command, args []string) {
	rc := 0
	cc := []string(nil)
	cc = append(cc, args...)
	if len(cc) == 0 {
		cc = gpiod.Chips()
	}
	for _, path := range cc {
		c, err := gpiod.NewChip(path)
		if err != nil {
			logErr(cmd, err)
			rc = 1
			continue
		}
		fmt.Printf("%s - %d lines:\n", c.Name, c.Lines())
		for o := 0; o < c.Lines(); o++ {
			li, err := c.LineInfo(o)
			if err != nil {
				logErr(cmd, err)
				rc = 1
				continue
			}
			printLineInfo(li)
		}
		c.Close()
	}
	os.Exit(rc)
}

func printLineInfo(li gpiod.LineInfo) {
	if len(li.Name) == 0 {
		li.Name = "unnamed"
	}
	if li.Requested {
		if len(li.Consumer) == 0 {
			li.Consumer = "kernel"
		}
		if strings.Contains(li.Consumer, " ") {
			li.Consumer = "\"" + li.Consumer + "\""
		}
	} else {
		li.Consumer = "unused"
	}
	dirn := "input"
	if li.IsOut {
		dirn = "output"
	}
	active := "active-high"
	if li.ActiveLow {
		active = "active-low"
	}
	flags := []string(nil)
	if li.Requested {
		flags = append(flags, "used")
	}
	if li.OpenDrain {
		flags = append(flags, "open-drain")
	}
	if li.OpenSource {
		flags = append(flags, "open-source")
	}
	flstr := ""
	if len(flags) > 0 {
		flstr = "[" + strings.Join(flags, " ") + "]"
	}
	fmt.Printf("\tline %3d:%12s%12s%8s%13s%s\n",
		li.Offset, li.Name, li.Consumer, dirn, active, flstr)
}
