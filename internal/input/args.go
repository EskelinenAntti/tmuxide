package input

import "slices"

type Args []string

func (args Args) ContainHelp() bool {
	return slices.Contains(args, "-h") || slices.Contains(args, "--help")
}

func (args Args) Command() string {
	return args[0]
}
