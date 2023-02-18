package synonyms

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gotwarlost/crossies/internal/htmlplus"
	"github.com/pkg/errors"
)

const baseURL = "https://wordhippo.com"

// Sort defines the sorting order in which synonyms are returned
type Sort string

// available sort orders
const (
	SortAlpha   Sort = "alpha"
	SortDisplay Sort = "display"
)

// Query is a query for synonyms
type Query struct {
	Word       string `json:"word,omitempty"`       // word or phrase for which to find a synonym
	Pattern    string `json:"pattern,omitempty"`    // RE2 pattern to match
	StartsWith string `json:"startsWith,omitempty"` // synonym starts with this string
	EndsWith   string `json:"endsWith,omitempty"`   // synonym ends with this string
	Sort       Sort   `json:"sort,omitempty"`       // sorting of output
	MinLetters int    `json:"minLetters,omitempty"` // min letters in synonym
	MaxLetters int    `json:"maxLetters,omitempty"` // max letters in synonym
	All        bool   `json:"all,omitempty"`        // whether to show all synonyms or just the closest ones
	pat        *regexp.Regexp
}

// NewQueryFromParams returns a query object from URL parameters
func NewQueryFromParams(values url.Values) (q Query, _ error) {
	q.Word = values.Get("word")
	if q.Word == "" {
		return q, fmt.Errorf("no word specified")
	}

	q.StartsWith = values.Get("startsWith")
	q.EndsWith = values.Get("endsWith")
	q.Pattern = values.Get("pattern")
	q.All = values.Get("all") == "true"
	q.Sort = SortDisplay
	if values.Get("sort") == string(SortAlpha) {
		q.Sort = SortAlpha
	}
	var err error
	minStr, maxStr := values.Get("minLetters"), values.Get("maxLetters")
	if minStr != "" {
		q.MinLetters, err = strconv.Atoi(minStr)
		if err != nil {
			return q, errors.Wrapf(err, "invalid min letters %q", minStr)
		}
	}
	if maxStr != "" {
		q.MaxLetters, err = strconv.Atoi(maxStr)
		if err != nil {
			return q, errors.Wrapf(err, "invalid max letters %q", maxStr)
		}
	}
	return q, nil
}

func (q *Query) normalize(text string) string {
	return strings.ToLower(text)
}

func (q *Query) initialize() error {
	q.Word = q.normalize(q.Word)
	q.StartsWith = q.normalize(q.StartsWith)
	q.EndsWith = q.normalize(q.EndsWith)
	q.Pattern = q.normalize(q.Pattern) // can this cause a problem?
	if q.Pattern != "" {
		var err error
		q.pat, err = regexp.Compile(q.Pattern)
		if err != nil {
			return errors.Wrapf(err, "bad regex %q", q.Pattern)
		}
	}
	if q.MinLetters > 0 && q.MaxLetters > 0 && q.MinLetters > q.MaxLetters {
		q.MinLetters, q.MaxLetters = q.MaxLetters, q.MinLetters
	}
	return nil
}

func (q *Query) shouldInclude(answer string, extended bool) bool {
	answer = q.normalize(answer)
	if !q.All && extended {
		return false
	}
	if q.MinLetters != 0 && len(answer) < q.MinLetters {
		return false
	}
	if q.MaxLetters != 0 && len(answer) > q.MaxLetters {
		return false
	}
	if q.StartsWith != "" && !strings.HasPrefix(answer, q.StartsWith) {
		return false
	}
	if q.EndsWith != "" && !strings.HasSuffix(answer, q.EndsWith) {
		return false
	}
	if q.pat != nil && !q.pat.MatchString(answer) {
		return false
	}
	return true
}

func (q *Query) sortEntries(entries []*Entry) {
	sort.SliceStable(entries, func(i, j int) bool {
		left := entries[i]
		right := entries[j]
		if q.Sort == SortAlpha {
			return strings.ToLower(left.Synonym) < strings.ToLower(right.Synonym)
		} else {
			return left.Priority < right.Priority
		}
	})
}

// Entry is a result of finding a synonym
type Entry struct {
	Synonym  string `json:"synonym,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

// Result is the result of synonym query
type Result struct {
	Query   *Query   `json:"query,omitempty"`   // query for which results are provided
	Entries []*Entry `json:"entries,omitempty"` // matching entries
}

// Run searches wordhippo for synonyms and returns results based on specified filters and sort order.
func (q *Query) Run() (*Result, error) {
	err := q.initialize()
	if err != nil {
		return nil, err
	}
	if q.Word == "" {
		return nil, fmt.Errorf("no word specified")
	}

	u := fmt.Sprintf("%s/what-is/another-word-for/%s.html", baseURL, url.PathEscape(q.Word))
	doc, err := htmlplus.LoadURL(u, htmlplus.LoadOptions{})
	if err != nil {
		return nil, err
	}

	sel := doc.FindAll("div.relatedwords > div.wb")
	uniq := map[string]*Entry{}
	counter := 0
	for _, node := range sel {
		counter++
		text := node.InnerText()
		priority := counter
		extended := node.AttributeValue("id") != ""
		if !q.shouldInclude(text, extended) {
			continue
		}
		if extended {
			priority += 10000
		}
		if _, ok := uniq[text]; !ok {
			uniq[text] = &Entry{
				Synonym:  text,
				Priority: priority,
			}
		}
	}
	if len(uniq) == 0 {
		return nil, fmt.Errorf("no synonyms for word %q that match the supplied filters", q.Word)
	}
	entries := make([]*Entry, 0, len(uniq))
	for _, e := range uniq {
		entries = append(entries, e)
	}
	q.sortEntries(entries)
	return &Result{Query: q, Entries: entries}, nil
}
