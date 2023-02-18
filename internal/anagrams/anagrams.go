package anagrams

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gotwarlost/crossies/internal/htmlplus"
	"github.com/gotwarlost/crossies/internal/inputerror"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://www.thewordfinder.com/anagram-solver/"
)

// Query is a query to find words matching a frame.
type Query struct {
	Phrase  string `json:"phrase"`
	Partial bool   `json:"partial,omitempty"`
}

func (q *Query) initialize() error {
	if q.Phrase == "" {
		return inputerror.New("empty phrase not allowed")
	}
	q.Phrase = strings.ReplaceAll(q.Phrase, " ", "")
	return nil
}

// NewQueryFromParams returns a query object from URL parameters
func NewQueryFromParams(values url.Values) (q Query, _ error) {
	q.Phrase = values.Get("phrase")
	partialStr := values.Get("partial")
	q.Partial = partialStr == "true"
	if err := q.initialize(); err != nil {
		return q, err
	}
	return q, nil
}

// Result is the result of a query
type Result struct {
	Phrases []string `json:"phrases"` // words found in current iteration
}

func Solve(query Query) (*Result, error) {
	if err := query.initialize(); err != nil {
		return nil, err
	}
	vals := url.Values{}
	vals.Set("letters", query.Phrase)
	vals.Set("extra", "")
	vals.Set("pos", "beg")
	vals.Set("dict", "wwf")
	vals.Set("dic", "1")
	vals.Set("order", "length")

	res, err := http.PostForm(baseURL, vals)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST %s return status %d", baseURL, res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	doc, err := htmlplus.Load(bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "read and parse HTML")
	}
	var ret []string
	nodes := doc.FindAll("p.result a")

	for _, node := range nodes {
		text := node.InnerText()
		if strings.EqualFold(text, query.Phrase) {
			continue
		}
		if len(text) < len(query.Phrase) && !query.Partial {
			break
		}
		ret = append(ret, text)
	}

	if len(ret) == 0 {
		return nil, inputerror.New(fmt.Sprintf("no anagrams found for %q", query.Phrase))
	}
	sort.Slice(ret, func(i, j int) bool {
		l1, l2 := len(ret[i]), len(ret[j])
		if l1 != l2 {
			return l1 > l2
		}
		return strings.ToLower(ret[i]) < strings.ToLower(ret[j])
	})
	return &Result{
		Phrases: ret,
	}, nil
}
