package window

import "github.com/eskelinenantti/tmuxide/internal/input"

type Path interface {
	Contains(path string) bool
}

type Window input.Args
