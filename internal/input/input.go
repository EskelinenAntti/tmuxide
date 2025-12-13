package input

import (
	"fmt"
	"os"
	"path/filepath"
)

type Args struct {
	Program string
	Path    string
}

func Parse(args []string) (Args, error) {
	var path string
	var program string
	var err error

	switch len(args) {
	case 1:
		path, err = os.Getwd()
	case 2:
		path, err = filepath.Abs(args[1])
	case 3:
		program = args[1]
		path, err = filepath.Abs(args[2])
	default:
		return Args{}, fmt.Errorf("Invalid number of arguments. See %s --help.", args[0])
	}

	return Args{Path: path, Program: program}, err
}
