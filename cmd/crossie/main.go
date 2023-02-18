package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

const exe = "crossie"

func setup() *cobra.Command {
	root := &cobra.Command{
		Use:   exe,
		Short: "crossword tools",
	}
	addSynonymsCommand(root)
	addFindWordsCommand(root)
	addAnagramsCommand(root)
	return root
}

func main() {
	log.SetFlags(0)
	cmd := setup()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
