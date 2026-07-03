// Package cli parses command-line and interactive options.
package cli

// PickMode returns the legacy argument tuple for argument or interactive mode.
func PickMode(args []string) (string, string, string, string, bool, error) {
	if len(args) == 1 {
		return RunInteractive()
	}

	text, bannerPath, color, substring, err := RunArgumentMode(args)
	return text, bannerPath, color, substring, false, err
}

// PickModeOptions returns the full option set for argument or interactive mode.
func PickModeOptions(args []string) (Options, bool, error) {
	if len(args) == 1 {
		return RunInteractiveOptions()
	}

	opts, err := RunArgumentModeOptions(args)
	return opts, false, err
}
