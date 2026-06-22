package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func printBannerMenu(out io.Writer) {
	fmt.Fprintln(out, "Choose the number of the desired banner:")
	fmt.Fprintln(out, "1. standard")
	fmt.Fprintln(out, "2. shadow")
	fmt.Fprintln(out, "3. thinkertoy")
}

func printColorPrompt(out io.Writer) {
	fmt.Fprintln(out, "Enter color, or press Enter for no color (/exit to quit)")
	fmt.Fprintln(out, "Examples: red, #ff0000, rgb,255,0,0, hsl,0,100%,50%")
	fmt.Fprint(out, "Color: ")
}

func printOutputPrompt(out io.Writer) {
	fmt.Fprint(out, "Write output to a file? (y/N, /exit to quit): ")
}

var interactiveReader *bufio.Reader
var interactiveReaderSource *os.File

func runInteractive(in io.Reader, out io.Writer) (string, string, string, string, bool, error) {
	opts, exit, err := collectInteractiveOptions(bufio.NewReader(in), out, false)
	return opts.Text, opts.BannerPath, opts.Color, opts.Substring, exit, err
}

func runInteractiveOptions(in io.Reader, out io.Writer) (Options, bool, error) {
	return collectInteractiveOptions(bufio.NewReader(in), out, true)
}

func RunInteractive() (string, string, string, string, bool, error) {
	if interactiveReader == nil || interactiveReaderSource != os.Stdin {
		interactiveReader = bufio.NewReader(os.Stdin)
		interactiveReaderSource = os.Stdin
	}

	opts, exit, err := collectInteractiveOptions(interactiveReader, os.Stdout, false)
	return opts.Text, opts.BannerPath, opts.Color, opts.Substring, exit, err
}

func RunInteractiveOptions() (Options, bool, error) {
	if interactiveReader == nil || interactiveReaderSource != os.Stdin {
		interactiveReader = bufio.NewReader(os.Stdin)
		interactiveReaderSource = os.Stdin
	}

	return collectInteractiveOptions(interactiveReader, os.Stdout, true)
}

func collectInteractiveOptions(reader *bufio.Reader, out io.Writer, askOutput bool) (Options, bool, error) {
	var text string

	for {
		fmt.Fprint(out, "Enter your text (or type '/exit' to quit): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return Options{}, true, err
		}

		text = strings.TrimRight(input, "\r\n")

		if text == "" {
			fmt.Fprintln(out, "No valid string input")
			continue
		}

		if text == "/exit" {
			return Options{}, true, nil
		}

		break
	}

	printBannerMenu(out)

	var bannerPath string

	for {
		fmt.Fprint(out, "Enter 1, 2 or 3 (/exit to quit): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return Options{}, true, err
		}

		input = strings.TrimSpace(input)

		switch input {
		case "1":
			bannerPath = "banners/standard.txt"
		case "2":
			bannerPath = "banners/shadow.txt"
		case "3":
			bannerPath = "banners/thinkertoy.txt"
		case "/exit":
			return Options{}, true, nil
		default:
			fmt.Fprintln(out, "Invalid banner choice")
			continue
		}

		break
	}

	printColorPrompt(out)

	input, err := reader.ReadString('\n')
	if err != nil {
		return Options{}, true, err
	}

	color := strings.TrimSpace(input)
	if color == "/exit" {
		return Options{}, true, nil
	}

	opts := Options{
		Text:       text,
		BannerPath: bannerPath,
		Color:      color,
	}

	if color == "" {
		return finishInteractiveOptions(reader, out, opts, askOutput)
	}

	fmt.Fprint(out, "Enter substring to color, or press Enter to color the whole string (/exit to quit): ")

	input, err = reader.ReadString('\n')
	if err != nil {
		return Options{}, true, err
	}

	substring := strings.TrimRight(input, "\r\n")
	if substring == "/exit" {
		return Options{}, true, nil
	}

	opts.Substring = substring
	return finishInteractiveOptions(reader, out, opts, askOutput)
}

func finishInteractiveOptions(reader *bufio.Reader, out io.Writer, opts Options, askOutput bool) (Options, bool, error) {
	if !askOutput {
		return opts, false, nil
	}

	for {
		printOutputPrompt(out)

		input, err := reader.ReadString('\n')
		if err != nil {
			return Options{}, true, err
		}

		choice := strings.ToLower(strings.TrimSpace(input))
		switch choice {
		case "":
			return opts, false, nil
		case "n", "no":
			return opts, false, nil
		case "y", "yes":
			return readInteractiveOutputFile(reader, out, opts)
		case "/exit":
			return Options{}, true, nil
		default:
			fmt.Fprintln(out, "Invalid output choice")
		}
	}
}

func readInteractiveOutputFile(reader *bufio.Reader, out io.Writer, opts Options) (Options, bool, error) {
	for {
		fmt.Fprint(out, "Enter output file name (/exit to quit): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return Options{}, true, err
		}

		outputFile := strings.TrimSpace(input)
		if outputFile == "/exit" {
			return Options{}, true, nil
		}
		if outputFile == "" {
			fmt.Fprintln(out, "No valid output file name")
			continue
		}

		opts.OutputFile = outputFile
		return opts, false, nil
	}
}
