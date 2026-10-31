# zipln 🪄

Open `zipln`, pick any file or folder from the fuzzy finder, and you'll find yourself in the right herdr workspace.

It doesn't matter if the workspace doesn't exist yet, or the file is buried in a subfolder of a subfolder. It'll *just work* 🪄

## Installation

Install zipln as a herdr plugin:

```bash
herdr plugin install eskelinenantti/zipln
```

You'll need herdr 0.9.1 or newer, Go 1.25.1 or newer to build the plugin, and `fd`, `fzf`, and `git`. If you are already using homebrew and are missing any of these dependencies, run `brew install fd fzf git`.

> [!TIP]
> Add the following config to your `~/.config/herdr/config.toml` to start jumping between folders and files from anywhere in herdr with the `prefix+f` shortcut.
>
> ```toml
> [[keys.command]]
> key = "prefix+f"
> type = "plugin_action"
> command = "zipln.open"
> description = "Open project or file"
> ```

## Manual

Your shortcut opens zipln in a popup where you can fuzzy find folders and files. If a workspace for the selected location already exists, zipln will switch to it. Otherwise, zipln creates the workspace.

If a required tool such as `fd` or `fzf` is missing, zipln displays installation instructions. The Herdr plugin popup keeps errors visible until you press Enter. After installing dependencies, restart Herdr if it still cannot find them on `$PATH`.

Picking a folder takes you to its workspace. Picking a file opens it in your editor, in a new tab in its workspace.

### Folder targets

```txt
path/to/project/
              \
               The herdr workspace is created for the absolute path of the selected folder.
```

### File targets

```txt
                 The selected file is opened in the editor configured by the $EDITOR variable
                /
path/to/file.txt
      \
       The herdr workspace is automatically created for the repository root of the selected file,
       or for the surrounding directory if the file isn't inside a git repository.
```

### tmux

zipln also supports tmux sessions. Install `zipln` with `go install github.com/eskelinenantti/zipln/zipln@latest` and add this to your `tmux.conf`:

```tmux
bind-key o "run 'zipln --mux tmux'"
```
