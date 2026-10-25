#!/bin/sh
set -eu

# The Herdr server may have a limited PATH.
export PATH="${PATH:+$PATH:}${HOME}/.local/bin:${HOME}/bin:/opt/homebrew/bin:/opt/homebrew/sbin:/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:${HOME}/.linuxbrew/bin:${HOME}/.linuxbrew/sbin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin"

exec "${HERDR_PLUGIN_ROOT:?HERDR_PLUGIN_ROOT is required}/bin/ide" --mux herdr "$@"
