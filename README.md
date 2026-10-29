# herdr-zipln 🪄

Press your shortcut, pick any file or folder from the fuzzy finder, you'll find yourself in a nicely named herdr workspace created for that particular location.

It doesn't matter if the workspace didn't yet exist. It'll *just work* 🪄

## Installation

Install zipln as a herdr plugin:

```bash
herdr plugin install eskelinenantti/herdr-zipln
```

You'll need herdr 0.9.1 or newer, Go 1.25.1 or newer to build the plugin, and `fd`, `fzf`, and `git` on your `$PATH`. Set `$EDITOR` to the editor you'd like to use for files.

> [!TIP]
> Add the following config to your `~/.config/herdr/config.toml` to start jumping between folders and files from anywhere in herdr. Press your prefix, then `o`.
>
> ```toml
> [[keys.command]]
> key = "prefix+o"
> type = "plugin_action"
> command = "herdr-zipln.open"
> description = "Open project or file"
> ```

## Manual

Your shortcut opens zipln in a popup where you can fuzzy find folders and files. If a workspace for given location already exists, zipln will switch to it. Otherwise, zipln creates the workspace.

Picking a folder takes you to its workspace. Picking a file opens it in your editor, in a new tab in its workspace.

### Folder targets

```txt
path/to/project/
              \
               The herdr workspace is created for the absolute path of the selected folder.
```

### File targets

```txt
                 The selected file is opened in editor configured by the $EDITOR variable
                /
path/to/file.txt
      \
       The herdr workspace is automatically created for the repository root of the selected file,
       or for the surrounding directory if file isn't inside a git repository.
```

### tmux

zipln also supports tmux sessions. With `zipln` on your `$PATH`, add this to your `tmux.conf`:

```tmux
bind-key o "run 'zipln --mux tmux'"
```
