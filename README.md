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

Running `ide` will start a fuzzy finder where you can fuzzy find sessions, folders and files. If a session for given location already exists, tmuxide will attach to it. Otherwise, tmuxide creates the session.

Alternatively, you can pass sessions, folders and files as argument to the command.

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

### Session targets

```txt
project-1a5f
           \
            Selecting a session just opens it. Nothing fancy here.
```

## Installation

You can install it with `homebrew`


```bash
brew install eskelinenantti/cli/tmuxide
```
