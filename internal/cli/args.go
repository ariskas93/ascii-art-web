package cli

import (
	"errors"
	"strings"
)

const UsageMessage = `Usage: go run . [STRING] [BANNER]

EX: go run . something standard`

const ColorUsageMessage = `Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"`

const OutputUsageMessage = `Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard`

const defaultBannerPath = "banners/standard.txt"

type Options struct {
	Text       string
	BannerPath string
	Color      string
	Substring  string
	OutputFile string
}

type positionalArg struct {
	value string
	index int
}

func RunArgumentModeOptions(args []string) (Options, error) {
	return parseFlexibleArguments(args, true)
}

func RunArgumentMode(args []string) (string, string, string, string, error) {
	opts, err := parseFlexibleArguments(args, false)
	if err != nil {
		return "", "", "", "", err
	}

	return opts.Text, opts.BannerPath, opts.Color, opts.Substring, nil
}

func parseFlexibleArguments(args []string, allowOutput bool) (Options, error) {
	if len(args) < 2 {
		return Options{}, errors.New(UsageMessage)
	}

	opts := Options{BannerPath: defaultBannerPath}
	operands, colorIndex, err := collectOptions(args, allowOutput, &opts)
	if err != nil {
		return Options{}, err
	}

	if err := applyOperands(&opts, operands, colorIndex); err != nil {
		return Options{}, err
	}

	return opts, nil
}

func collectOptions(args []string, allowOutput bool, opts *Options) ([]positionalArg, int, error) {
	operands := []positionalArg{}
	colorIndex := -1
	seenColor := false
	seenOutput := false

	for i := 1; i < len(args); i++ {
		arg := args[i]

		switch {
		case isColorSyntax(arg):
			if seenColor {
				return nil, -1, errors.New(ColorUsageMessage)
			}
			color, err := parseColorFlag(arg)
			if err != nil {
				return nil, -1, err
			}
			opts.Color = color
			colorIndex = i
			seenColor = true
		case isOutputSyntax(arg):
			if !allowOutput || seenOutput {
				return nil, -1, errors.New(OutputUsageMessage)
			}
			outputFile, err := parseOutputFlag(arg)
			if err != nil {
				return nil, -1, err
			}
			opts.OutputFile = outputFile
			seenOutput = true
		case strings.HasPrefix(arg, "--"):
			return nil, -1, errors.New(UsageMessage)
		default:
			operands = append(operands, positionalArg{
				value: arg,
				index: i,
			})
		}
	}

	return operands, colorIndex, nil
}

func applyOperands(opts *Options, operands []positionalArg, colorIndex int) error {
	if opts.Color != "" {
		return applyColorOperands(opts, operands, colorIndex)
	}

	text, bannerPath, err := parseTextBannerOperands(operands)
	if err != nil {
		if len(operands) == 0 && opts.OutputFile != "" {
			return errors.New(OutputUsageMessage)
		}
		return err
	}

	opts.Text = text
	opts.BannerPath = bannerPath
	return nil
}

func applyColorOperands(opts *Options, operands []positionalArg, colorIndex int) error {
	if len(operands) == 0 {
		return errors.New(ColorUsageMessage)
	}

	text, bannerPath, err := parseTextBannerOperands(operands)
	if err == nil {
		opts.Text = text
		opts.BannerPath = bannerPath
		return nil
	}

	candidateIndex := findSubstringCandidate(operands, colorIndex)
	if candidateIndex >= 0 {
		remaining := removeOperand(operands, candidateIndex)
		if len(remaining) == 0 {
			return errors.New(ColorUsageMessage)
		}

		text, bannerPath, err = parseTextBannerOperands(remaining)
		if err == nil {
			opts.Text = text
			opts.BannerPath = bannerPath
			opts.Substring = operands[candidateIndex].value
			return nil
		}

		return errors.New(ColorUsageMessage)
	}

	if len(operands) > 2 {
		return errors.New(ColorUsageMessage)
	}

	return err
}

func parseTextBannerOperands(operands []positionalArg) (string, string, error) {
	if len(operands) == 0 || len(operands) > 2 {
		return "", "", errors.New(UsageMessage)
	}

	text := operands[0].value
	if text == "" {
		return "", "", errors.New(UsageMessage)
	}

	bannerPath := defaultBannerPath
	if len(operands) == 2 {
		var err error
		bannerPath, err = normalizeBannerArg(operands[1].value)
		if err != nil {
			return "", "", err
		}
	}

	return text, bannerPath, nil
}

func findSubstringCandidate(operands []positionalArg, colorIndex int) int {
	for i, arg := range operands {
		if arg.index == colorIndex+1 {
			return i
		}
	}
	return -1
}

func removeOperand(operands []positionalArg, index int) []positionalArg {
	remaining := make([]positionalArg, 0, len(operands)-1)
	remaining = append(remaining, operands[:index]...)
	remaining = append(remaining, operands[index+1:]...)
	return remaining
}

func parseColorFlag(flag string) (string, error) {
	if !strings.HasPrefix(flag, "--color=") {
		return "", errors.New(ColorUsageMessage)
	}
	color := strings.TrimPrefix(flag, "--color=")
	if color == "" {
		return "", errors.New(ColorUsageMessage)
	}
	return color, nil
}

func parseOutputFlag(flag string) (string, error) {
	if !strings.HasPrefix(flag, "--output=") {
		return "", errors.New(OutputUsageMessage)
	}
	outputFile := strings.TrimPrefix(flag, "--output=")
	if outputFile == "" {
		return "", errors.New(OutputUsageMessage)
	}
	return outputFile, nil
}

func isOutputSyntax(arg string) bool {
	return arg == "--output" ||
		strings.HasPrefix(arg, "--output=") ||
		strings.HasPrefix(arg, "-output=") ||
		strings.HasPrefix(arg, "--outputs=")
}

func isColorSyntax(arg string) bool {
	return arg == "--color" ||
		strings.HasPrefix(arg, "--color=") ||
		strings.HasPrefix(arg, "-color=") ||
		strings.HasPrefix(arg, "--colour=")
}

func normalizeBannerArg(banner string) (string, error) {
	switch banner {
	case "standard", "standard.txt":
		return "banners/standard.txt", nil
	case "shadow", "shadow.txt":
		return "banners/shadow.txt", nil
	case "thinkertoy", "thinkertoy.txt":
		return "banners/thinkertoy.txt", nil
	default:
		return "", errors.New(UsageMessage)
	}
}
