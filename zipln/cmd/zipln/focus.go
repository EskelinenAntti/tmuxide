package cmd

import (
	"fmt"
	"os"

	"github.com/eskelinenantti/herdr-zipln/internal/history"
	"github.com/eskelinenantti/herdr-zipln/internal/picker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newFocusCmd())
}

func newFocusCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "_focus query matches",
		Hidden: true,
		Args:   cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do not override navigation made while the query was being matched.
			preferred := history.Preferred(args[0])
			if os.Getenv("FZF_ACTION") != "rebind" || preferred == "" {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "unbind(result)")
				return err
			}
			matches, err := os.Open(args[1])
			if err != nil {
				return err
			}
			defer matches.Close()
			position, err := picker.PreferredPosition(preferred, matches)
			if err != nil {
				return err
			}
			if position == 0 {
				// Keep waiting for a preferred path that is still being scanned.
				return nil
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "pos(%d)+track-current+unbind(result)\n", position)
			return err
		},
	}
}
