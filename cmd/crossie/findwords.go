package main

import (
	"fmt"

	"github.com/gotwarlost/crossies/internal/findwords"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func addFindWordsCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "find-words frame",
		Aliases: []string{"find"},
		Short:   "find words with missing values indicated by the . character",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("exactly one word must be specified")
			}
			cmd.SilenceUsage = true
			next := 1
			for {
				q := findwords.Query{Frame: args[0], Page: next}
				result, err := q.Run()
				if err != nil {
					return errors.Wrap(err, "find words")
				}
				for _, w := range result.Words {
					fmt.Println(w)
				}
				next = result.NextPage
				if next == 0 {
					break
				}
			}
			return nil
		},
	}
	root.AddCommand(cmd)
}
