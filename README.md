# tmuxide 🪄

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

This is where tmuxide comes into play. All you need to do is

```bash
ide edit path/to/my/project
```
and you'll find yourself in nicely named session with the given file or folder open in your favourite editor.

It doesn't matter if you are already in tmux or not, or whether you are already in that session or in some other session. It'll *just work* 🪄

#### Advanced usage

Another useful command you can do is

```bash
ide open path/to/my/project
```

This creates a new session or attaches to a session in the given directory, and opens it. You can also specify a command with the path, e.g.

```bash
ide open path/to/my/project lazygit
```

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

