# tmuxide 🪄

Run `ide open path/to/dir` and you'll find yourself in a nicely named tmux session created for that particular directory. 

It doesn't matter if you run it inside or outside tmux, or if the session didn't yet exist. It'll *just work* 🪄

### Open session with command

After the dir, you can also pass in any command with its arguments, and it'll get run in the opened session:

```bash
ide open path/to/dir lazygit
```

## Edit files

tmuxide comes with handy shortcut command for editing files within tmux sessions:

```bash
#        The file given as argument is opened in editor configured by the $EDITOR variable
#       / 
ide edit path/to/dir/some/file.txt
#                  \
#                   The session is automatically created for the repository root of the given file,
#                   or for the surrounding directory if file isn't inside a git repository. 
```

This is essentially the same as running `ide open path/to/dir $EDITOR some/file.txt`.

## Installation

You can install it with `homebrew`

```bash
brew install eskelinenantti/cli/tmuxide
```

Alternatively, if you prefer to use `go`, you can run

```bash
go install github.com/eskelinenantti/tmuxide/cmd/ide@latest
```

### Requirements
- [tmux](https://github.com/tmux/tmux)

### Recommended to be used with
- [fzf](https://github.com/junegunn/fzf) fuzzy finder. With fzf you can simply type `ide open **<tab>` and fuzzy find your way to your project.
- Any terminal based editor.

### Example configurations for related tools
- [zsh-configuration](https://github.com/EskelinenAntti/zsh-configuration)
- [tmux-configuration](https://github.com/EskelinenAntti/tmux-configuration)
- [neovim-configuration](https://github.com/EskelinenAntti/neovim-configuration)
