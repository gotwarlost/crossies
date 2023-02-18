package main

import (
	"fmt"
	"strings"

	"github.com/gotwarlost/crossies/internal/synonyms"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func addSynonymsCommand(root *cobra.Command) {
	var q synonyms.Query
	var alphaSort bool
	cmd := &cobra.Command{
		Use:     "synonyms",
		Aliases: []string{"syn"},
		Short:   "get synonyms for the specified word or phrase from wordhippo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no word or phrase specified")
			}
			cmd.SilenceUsage = true
			q.Word = strings.Join(args, " ")
			if alphaSort {
				q.Sort = synonyms.SortAlpha
			} else {
				q.Sort = synonyms.SortDisplay
			}
			result, err := q.Run()
			if err != nil {
				return errors.Wrap(err, "find synonyms")
			}
			for _, e := range result.Entries {
				fmt.Println(e.Synonym)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&q.StartsWith, "starts", "s", "", "letters that the synonym should start with")
	f.StringVarP(&q.EndsWith, "ends", "e", "", "letters that the synonym should end with")
	f.StringVarP(&q.Pattern, "pattern", "p", "", "RE2 pattern against which to match synonyms")
	f.IntVarP(&q.MinLetters, "min", "m", 0, "minimum letters that the synonym should have")
	f.IntVarP(&q.MaxLetters, "max", "M", 0, "maximum letters that the synonym should have (0=any number)")
	f.BoolVar(&alphaSort, "sort", false, "return words in alphabetical order")
	f.BoolVar(&q.All, "all", false, "display all synonyms including ones that are hidden behind the 'More...' link in wordhippo")
	root.AddCommand(cmd)
}
