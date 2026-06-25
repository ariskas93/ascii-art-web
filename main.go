package main

import (
	"fmt"
	"os"

	"ascii-art/internal/banner"
	"ascii-art/internal/cli"
	"ascii-art/internal/output"
	"ascii-art/internal/renderer"
)

func main() {
	interactive := len(os.Args) == 1

	for {
		opts, exit, err := cli.PickModeOptions(os.Args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if exit {
			return
		}

		bannerMap, err := banner.LoadBanner(opts.BannerPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		lines := renderer.NormalizeInput(opts.Text)
		ascii, err := renderer.BuildASCII(lines, bannerMap, opts.Substring, opts.Color)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if interactive {
				continue
			}
			os.Exit(1)
		}

		if opts.OutputFile != "" {
			err = output.WriteASCIIFile(opts.OutputFile, ascii)
		} else {
			err = output.PrintASCII(ascii)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !interactive {
			return
		}
	}
}
