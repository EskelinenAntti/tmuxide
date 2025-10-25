# tmuxide

Turn your terminal into an IDE by automating tmux session and window management with tmuxide.

## Enough fancy words, what does it do?

The traditional way of using tmux can be rather tedious. When creating a session, you need to
- Come up with a name for the session (or run in trouble later), and run `tmux new -s my-project` 
- Create tmux windows manually, either with key shortcuts or `tmux new-window`

When you want to reattach to the same session again, you need to
- Remember the name of the session, and remember if there was a session in the first place. Alternatively check it with `tmux ls` (this is the step where you run into trouble if you didn't name your session earlier).
- Run `tmux attach -t project` or `tmux switch -t project` depending if you are already inside tmux or not.

That's quite a lot to remember.

### Enter tmuxide

This is where tmuxide comes into play. Only command you'll need is

```bash
ide path/to/my/project
```
and you'll find yourself in beautifully named session with your favourite editor and lazygit already open in their own windows.

In order to reattach to the same session later on, all you need to do is to run the same command again, with path to the same directory or file as an argument.

It doesn't matter if you are already in tmux or not, or whether you use absolute or relative path. It'll *just work*. 🪄

## Installation

You can install it with `homebrew`

```bash
brew install eskelinenantti/cli/tmuxide
```

Alternatively, if you prefer to use `go`, you can run

```bash
go install github.com/eskelinenantti/tmuxide/cmd/ide
```

### Requirements
- [tmux](https://github.com/tmux/tmux)

### Recommended to be used with
- [fzf](https://github.com/junegunn/fzf) fuzzy finder. With fzf you can simply type `ide **<tab>` and fuzzy find your way to your project.
- [lazygit](https://github.com/jesseduffield/lazygit) is a terminal UI for Git. tmuxide recognizes whether the opened folder or file is within a git repository, and if so, opens lazygit in second window.

