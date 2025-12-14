package input

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type Args struct {
	Program    string
	Path       string
	WorkingDir string
}

type flags struct {
	workingDir string
}

var ErrMissingFlagValue = errors.New("Missing flag value")

func Parse(args []string) (Args, error) {
	var err error

	flags, flagPositions, err := parseFlags(args)
	if err != nil {
		return Args{}, err
	}

	positionalArgs := parsePositionalArgs(args, flagPositions)
	var path string
	var program string

	switch len(positionalArgs) {
	case 1:
		path, err = os.Getwd()
	case 2:
		program = positionalArgs[1]
	case 3:
		program = positionalArgs[1]
		path, err = filepath.Abs(positionalArgs[2])
	default:
		return Args{}, fmt.Errorf("Invalid number of arguments. See %s --help.", args[0])
	}

	return Args{Path: path, Program: program, WorkingDir: flags.workingDir}, err
}

func parseFlags(args []string) (flags, []int, error) {
	var flags = flags{}
	var flagIndexes = []int{}
	var err error
	for i, arg := range args {
		switch arg {
		case "-c":
			if i+1 >= len(args) {
				return flags, []int{}, fmt.Errorf("%w: %s", ErrMissingFlagValue, arg)
			}
			flags.workingDir, err = filepath.Abs(args[i+1])
			flagIndexes = append(flagIndexes, i, i+1)
		}
	}
	return flags, flagIndexes, err
}

func parsePositionalArgs(args []string, flagPositions []int) []string {
	positionalArgs := []string{}
	for i := range args {
		if !slices.Contains(flagPositions, i) {
			positionalArgs = append(positionalArgs, args[i])
		}
	}
	return positionalArgs
}
