package input

import (
	"fmt"
	"os"
	"path/filepath"
)

func Path(args Args) (string, error) {
	var target string
	var err error

	switch len(args) {
	case 1:
		target, err = os.Getwd()
	case 2:
		target, err = filepath.Abs(args[1])
	default:
		return "", fmt.Errorf("Invalid number of arguments. See %s --help.", args.Command())
	}

	return target, err
}
