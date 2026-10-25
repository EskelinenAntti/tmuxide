#!/bin/sh
set -eu

# The Herdr server may have a limited PATH.
export PATH="${PATH:+$PATH:}/opt/homebrew/bin:/usr/local/bin:/home/linuxbrew/.linuxbrew/bin"

exec "${HERDR_PLUGIN_ROOT:?HERDR_PLUGIN_ROOT is required}/bin/ide" --mux herdr "$@"
