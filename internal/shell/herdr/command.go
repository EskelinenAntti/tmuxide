package herdr

import "strings"

// shellCommand preserves argv boundaries when pane run submits a single shell
// command. Single quotes keep paths and editor flags from becoming shell syntax;
// embedded single quotes must be emitted outside the surrounding quoted string.
func shellCommand(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
	}
	return strings.Join(quoted, " ")
}
