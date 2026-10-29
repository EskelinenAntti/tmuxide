package cmd

import (
	"os"

	"github.com/eskelinenantti/herdr-zipln/internal/ranked"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:    "_rank candidates query",
		Hidden: true,
		Args:   cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer input.Close()
			return ranked.Filter(input, cmd.OutOrStdout(), args[1])
		},
	})
}
