package main

import (
	"fmt"
	"strings"

	"github.com/gotwarlost/crossies/internal/anagrams"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func addAnagramsCommand(root *cobra.Command) {
	var q anagrams.Query
	cmd := &cobra.Command{
		Use:     "anagrams phrase",
		Aliases: []string{"anag"},
		Short:   "get anagrams for the supplied phrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no word or phrase specified")
			}
			cmd.SilenceUsage = true
			q.Phrase = strings.Join(args, " ")
			cmd.SilenceUsage = true

			result, err := anagrams.Solve(q)
			if err != nil {
				return errors.Wrap(err, "find anagrams")
			}
			for _, p := range result.Phrases {
				fmt.Println(p)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&q.Partial, "partial", "p", false, "return partial anagrams")
	root.AddCommand(cmd)
}
