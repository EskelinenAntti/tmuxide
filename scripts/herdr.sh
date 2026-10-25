#!/bin/sh

# GUI launches may inherit a minimal PATH. Keep existing entries first.
PATH="${PATH:+$PATH:}/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin"
script_dir=$(CDPATH= cd -P "${0%/*}" && pwd) || exit 1
export PATH

exec "$script_dir/../bin/ide" "$@"
