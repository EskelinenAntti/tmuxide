# tmuxide 🪄

Run `ide`, pick any file or folder from the fuzzy finder, you'll find yourself in a nicely named tmux session created for that particular location.

It doesn't matter if you run it inside or outside tmux, or if the session didn't yet exist. It'll *just work* 🪄

> [!TIP]
> Add the following config to your `tmux.conf` to start jumping between folders and files from anywhere.
>
> ```
> bind-key o "run ide"
> ```

## Manual

Running `ide` will start a fuzzy finder where you can fuzzy find folders and files. If a session for given location already exists, tmuxide will attach to it. Otherwise, tmuxide creates the session.

Alternatively, you can pass folders and files as argument to the command.

### Search ranking

tmuxide remembers how often you select each file or folder, including targets
passed as arguments. Previously selected paths appear first, ordered by selection
count (ties are alphabetical). Matches keep this order while you type. Selected
folders are also searched in frequency order before the full home-directory scan,
and duplicate results and missing history paths are skipped.

History is stored locally in `$XDG_STATE_HOME/tmuxide/history.json`, or
`~/.local/state/tmuxide/history.json` when `XDG_STATE_HOME` is unset. Delete that
file to reset the ranking. Search and navigation still work if history cannot be
read or saved.

### Folder targets

```txt
path/to/project/
              \
               The session is created for the absolute path of the selected folder.
```

### File targets

```txt
                 The file given as argument is opened in editor configured by the $EDITOR variable
                /
path/to/file.txt
      \
       The session is automatically created for the repository root of the given file,
       or for the surrounding directory if file isn't inside a git repository.
```

## Installation

You can install it with `homebrew`


```bash
brew install eskelinenantti/cli/tmuxide
```
