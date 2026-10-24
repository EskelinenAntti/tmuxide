package ide

import "github.com/eskelinenantti/tmuxide/internal/project"

// Workspace owns the host-specific create, reuse, and focus lifecycle.
// Project discovery and editor selection do not depend on the terminal host.
type Workspace interface {
	Open(project.Project, []string) error
}

func Start(command []string, project project.Project, workspace Workspace) error {
	return workspace.Open(project, command)
}
