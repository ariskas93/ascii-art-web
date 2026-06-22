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

func RunArgumentModeOptions(args []string) (Options, error) {
	if len(args) < 2 {
		return Options{}, errors.New(UsageMessage)
	}

	outputFile := ""
	argsToParse := args

	if hasOutputSyntax(args[1:]) {
		var err error
		outputFile, argsToParse, err = parseOutputMode(args)
		if err != nil {
			return Options{}, err
		}
	}

	text, bannerPath, color, substring, err := RunArgumentMode(argsToParse)
	if err != nil {
		return Options{}, err
	}

	return Options{
		Text:       text,
		BannerPath: bannerPath,
		Color:      color,
		Substring:  substring,
		OutputFile: outputFile,
	}, nil
}

func RunArgumentMode(args []string) (string, string, string, string, error) {
	if len(args) < 2 {
		return "", "", "", "", errors.New(UsageMessage)
	}

	if len(args) > 5 {
		if hasColorSyntax(args[1:]) {
			return "", "", "", "", errors.New(ColorUsageMessage)
		}
		return "", "", "", "", errors.New(UsageMessage)
	}

	if countColorSyntax(args[1:]) > 1 {
		return "", "", "", "", errors.New(ColorUsageMessage)
	}

	if hasColorSyntax(args[1:]) {
		return parseColorMode(args)
	}

	return parsePlainMode(args)
}

func parsePlainMode(args []string) (string, string, string, string, error) {
	if len(args) > 3 {
		return "", "", "", "", errors.New(UsageMessage)
	}

	text := args[1]
	if text == "" {
		return "", "", "", "", errors.New(UsageMessage)
	}

	bannerPath := defaultBannerPath
	if len(args) == 3 {
		var err error
		bannerPath, err = normalizeBannerArg(args[2])
		if err != nil {
			return "", "", "", "", err
		}
	}

	return text, bannerPath, "", "", nil
}

func parseColorMode(args []string) (string, string, string, string, error) {
	if isColorSyntax(args[1]) {
		color, err := parseColorFlag(args[1])
		if err != nil {
			return "", "", "", "", err
		}
		switch len(args) {
		case 3:
			if args[2] == "" {
				return "", "", "", "", errors.New(ColorUsageMessage)
			}
			return args[2], defaultBannerPath, color, "", nil
		case 4:
			if args[3] == "" {
				return "", "", "", "", errors.New(ColorUsageMessage)
			}
			return args[3], defaultBannerPath, color, args[2], nil
		default:
			return "", "", "", "", errors.New(ColorUsageMessage)
		}
	}

	if len(args) == 4 || len(args) == 5 {
		if !isColorSyntax(args[3]) {
			return "", "", "", "", errors.New(ColorUsageMessage)
		}
		text := args[1]
		if text == "" {
			return "", "", "", "", errors.New(UsageMessage)
		}
		bannerPath, err := normalizeBannerArg(args[2])
		if err != nil {
			return "", "", "", "", err
		}
		color, err := parseColorFlag(args[3])
		if err != nil {
			return "", "", "", "", err
		}
		substring := ""
		if len(args) == 5 {
			substring = args[4]
		}
		return text, bannerPath, color, substring, nil
	}

	return "", "", "", "", errors.New(ColorUsageMessage)
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

func parseOutputMode(args []string) (string, []string, error) {
	if len(args) < 2 || !strings.HasPrefix(args[1], "--output=") {
		return "", nil, errors.New(OutputUsageMessage)
	}

	outputFile := strings.TrimPrefix(args[1], "--output=")
	if outputFile == "" {
		return "", nil, errors.New(OutputUsageMessage)
	}

	if countOutputSyntax(args[1:]) > 1 {
		return "", nil, errors.New(OutputUsageMessage)
	}

	argsToParse := append([]string{args[0]}, args[2:]...)
	if len(argsToParse) < 2 {
		return "", nil, errors.New(OutputUsageMessage)
	}

	return outputFile, argsToParse, nil
}

func hasOutputSyntax(args []string) bool {
	for _, arg := range args {
		if isOutputSyntax(arg) {
			return true
		}
	}
	return false
}

func countOutputSyntax(args []string) int {
	count := 0
	for _, arg := range args {
		if isOutputSyntax(arg) {
			count++
		}
	}
	return count
}

func isOutputSyntax(arg string) bool {
	return arg == "--output" ||
		strings.HasPrefix(arg, "--output=") ||
		strings.HasPrefix(arg, "-output=") ||
		strings.HasPrefix(arg, "--outputs=")
}

func hasColorSyntax(args []string) bool {
	for _, arg := range args {
		if isColorSyntax(arg) {
			return true
		}
	}
	return false
}

func countColorSyntax(args []string) int {
	count := 0
	for _, arg := range args {
		if isColorSyntax(arg) {
			count++
		}
	}
	return count
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
