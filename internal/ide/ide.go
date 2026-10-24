package ide

import "github.com/eskelinenantti/tmuxide/internal/project"

type Multiplexer interface {
	Open(project.Project, []string) error
}

func Start(command []string, project project.Project, mux Multiplexer) error {
	return mux.Open(project, command)
}
