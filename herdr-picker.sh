#!/bin/sh
# Actions have no terminal. Open the interactive picker in a managed popup.
exec "${HERDR_BIN_PATH:-herdr}" plugin pane open --plugin tmuxide --entrypoint picker
