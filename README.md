# tmuxide 🪄

Run `ide`, pick any file or folder from the fuzzy finder, and jump to a nicely named project workspace in tmux or Herdr.

Inside Herdr, tmuxide uses native Herdr workspaces and tabs. Otherwise it uses tmux, whether you run it inside or outside a session. Existing project workspaces are reused. 🪄

> [!TIP]
> Add the following config to your `tmux.conf` to start jumping between folders and files from anywhere.
>
> ```
> bind-key o "run ide"
> ```

## Manual

Running `ide` starts a fuzzy finder for folders and files under your home directory. If a workspace for the location exists, tmuxide focuses it; otherwise it creates one. With tmux, a workspace is a session; with Herdr, it is a native workspace.

Alternatively, you can pass folders and files as argument to the command.

```sh
ide                         # auto-detect Herdr; otherwise use tmux
ide path/to/project
ide path/to/file.txt
ide --backend tmux          # explicitly use tmux
ide --backend herdr         # explicitly use Herdr
```

Selecting a file opens `$EDITOR` in the new workspace's first terminal, or a new window/tab in an existing workspace. File paths are resolved before changing to the project directory.

The Herdr backend requires a running Herdr session. When invoked outside Herdr, it focuses the project and then attaches. Start `herdr` once first; tmuxide does not start or supervise a Herdr server. `HERDR_SESSION` or `HERDR_SOCKET_PATH` can select a different running session. Inside Herdr, the inherited session context is used, including when Herdr itself is nested inside tmux.

### Folder targets

```txt
path/to/project/
              \
                The workspace is created for the absolute path of the selected folder.
```

### File targets

```txt
                 The file given as argument is opened in editor configured by the $EDITOR variable
                /
path/to/file.txt
      \
       The workspace is automatically created for the repository root of the given file,
       or for the surrounding directory if file isn't inside a git repository.
```

## Installation

You can install it with `homebrew`


```bash
brew install eskelinenantti/cli/tmuxide
```

Required tools: `fd`, `fzf`, `git`, the selected host (`tmux` or `herdr`), and an editor configured with `$EDITOR`. Herdr does not require tmux.

### Herdr plugin

Requires Herdr 0.9.5+, Go 1.25.1+ for the install-time build, and the tools above.

```sh
herdr plugin install eskelinenantti/tmuxide
```

The plugin builds its own `ide` binary; a separate Homebrew installation is not needed. It does not add `ide` to your shell's `PATH`.

Bind a key in your Herdr configuration:

```toml
[[keys.command]]
key = "prefix+o"
type = "plugin_action"
command = "tmuxide.open"
description = "open project or file"
```

The action opens the fuzzy finder in a temporary popup. Selecting a folder creates or focuses its native workspace; selecting a file opens the editor in a native tab. Cancelling closes the popup without changing the layout.

You can also launch the picker directly:

```sh
herdr plugin action invoke tmuxide.open
```

### Development

```sh
go build -o bin/ide ./ide
herdr plugin link .         # link does not run the build step
go test ./...
```

Project discovery, selection, and editor configuration are host-independent. The `ide.Workspace` interface defines the host boundary; `internal/workspace` selects and validates a backend, and the tmux and Herdr adapters own their respective create/reuse/focus lifecycles. The plugin manifest is only packaging and a terminal entrypoint, not a second implementation of the workflow.

An optional integration test starts an isolated headless Herdr server, exercises native workspace reuse and editor tabs, verifies quoted file paths through real terminal output, and validates the plugin manifest. It does not touch your running sessions or plugin registry:

```sh
HERDR_TEST_BIN=/path/to/herdr go test -v ./internal/shell/herdr
```
