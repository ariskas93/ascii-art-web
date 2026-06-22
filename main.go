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
		text, bannerPath, color, substring, exit, err := cli.PickMode(os.Args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if exit {
			return
		}

		bannerMap, err := banner.LoadBanner(bannerPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		lines := renderer.NormalizeInput(text)
		ascii, err := renderer.BuildASCII(lines, bannerMap, substring, color)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if interactive {
				continue
			}
			os.Exit(1)
		}

		err = output.PrintASCII(ascii)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !interactive {
			return
		}
	}
}
