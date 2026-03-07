# tmuxide 🪄

Run `ide`, pick a session, file or folder from the fuzzy finder, you'll find yourself in a nicely named tmux session created for that particular location.

It doesn't matter if you run it inside or outside tmux, or if the session didn't yet exist. It'll *just work* 🪄

> [!TIP]
> Add the following config to your `tmux.conf` to start jumping between sessions, folders, and files from anywhere.
>
> ```
> bind-key o "neww 'ide'"
> ```

## Editing files

When you pick a file, tmuxide will open it in an editor window in addition to opening the sesssion:

```bash
#                         The file given as argument is opened in editor configured by the $EDITOR variable
#                        /
path/to/dir/some/file.txt
#           \
#            The session is automatically created for the repository root of the given file,
#            or for the surrounding directory if file isn't inside a git repository.
```

## Installation

You can install it with `homebrew`

```bash
brew install eskelinenantti/cli/tmuxide
```

Alternatively, if you prefer to use `go`, you can run

```bash
go install github.com/eskelinenantti/tmuxide/cmd/ide@latest
```
