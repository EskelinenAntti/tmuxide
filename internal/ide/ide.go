package ide

import "github.com/eskelinenantti/tmuxide/internal/project"

// Multiplexer owns the host-specific create, reuse, and focus lifecycle.
// Project discovery and editor selection do not depend on the terminal host.
type Multiplexer interface {
	Open(project.Project, []string) error
}

func Start(command []string, project project.Project, mux Multiplexer) error {
	return mux.Open(project, command)
}
