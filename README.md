# tmuxide 🪄

Run `ide`, pick any file or folder from the fuzzy finder, you'll find yourself in a nicely named tmux session created for that particular location.

It doesn't matter if you run it inside or outside tmux, or if the session didn't yet exist. It'll *just work* 🪄

> [!TIP]
> Add the following config to your `tmux.conf` to start jumping between sessions, folders, and files from anywhere.
>
> ```
> bind-key o "run ide"
> ```

## Manual

Running `ide` will start a fuzzy finder where you can fuzzy find sessions, folders and files.

Alternatively, you can pass sessions, folders and files as argument to the command.

### Folder targets

```bash
ide project/
#          \
#           The session is created for the absolute path of the selected folder.
```

### File targets

```bash
#                         The file given as argument is opened in editor configured by the $EDITOR variable
#                        /
ide project/dir/file.txt
#          \
#           The session is automatically created for the repository root of the given file,
#           or for the surrounding directory if file isn't inside a git repository.
```

### Session targets

```bash
ide project-1a5f
#               \
#                Opens the session. A shortcut for `tmux attach` and `tmux switch` which works inside and outside tmux sessions.
```

## Installation

You can install it with `homebrew`


```bash
brew install eskelinenantti/cli/tmuxide
```
