package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tmuxide",
	Short: "",
	Long:  ``,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	var selectedPath string
	var err error
	switch len(args) {
	case 0:
		selectedPath, err = os.Getwd()
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}
	case 1:
		selectedPath, err = dirAbsPath(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("An unexpected error occurred")
		return
	}

	fmt.Println(selectedPath)
}

// Get absolute path to enclosing directory, or to the directory itself if it's a file.
func dirAbsPath(inputPath string) (string, error) {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return "", err
	}
	absolutePath, err := filepath.Abs(inputPath)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return filepath.Dir(absolutePath), nil
	}
	return absolutePath, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
