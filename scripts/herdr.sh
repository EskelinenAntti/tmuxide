#!/bin/sh

# GUI launches may inherit a minimal PATH. Keep existing entries first.
PATH="${PATH:+$PATH:}/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin"
export PATH

exec "$(dirname "$0")/../bin/ide" "$@"
