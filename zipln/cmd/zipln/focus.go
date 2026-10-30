package cmd

import (
	"fmt"
	"os"

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
			if os.Getenv("FZF_ACTION") != "rebind" {
				return nil
			}
			matches, err := os.Open(args[1])
			if err != nil {
				return err
			}
			defer matches.Close()
			position, err := picker.PreferredPosition(args[0], matches)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "pos(%d)\n", position)
			return err
		},
	}
}
