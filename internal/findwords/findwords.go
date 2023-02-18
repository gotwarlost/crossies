package findwords

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gotwarlost/crossies/internal/htmlplus"
	"github.com/gotwarlost/crossies/internal/inputerror"
	"github.com/gotwarlost/crossies/internal/synonyms"
	"github.com/pkg/errors"
)

const (
	placeholder      = "."
	infoPageSize     = 250
	infoPagesPerPage = 1
)

var (
	inputRE      = regexp.MustCompile(`^[a-zA-Z.]+$`)
	totalWordsRE = regexp.MustCompile(`There\s+are\s+(\d+)\s+`)
	scoreRE      = regexp.MustCompile(`[(].*`)
)

func getURL(frame string, page int) (string, url.Values, error) {
	count := len(frame)
	specifiedCount := 0
	word := ""
	for i := 0; i < count; i++ {
		ch := frame[i : i+1]
		if ch != placeholder {
			specifiedCount++
			word += ch
		} else {
			word += "_"
		}
	}
	if specifiedCount == 0 {
		return "", nil, fmt.Errorf("inputs cannot all be dots")
	}
	return fmt.Sprintf("https://www.thewordfinder.com/wordlist/at-position-%s/", word),
		url.Values{
			"dir":   []string{"ascending"},
			"field": []string{"word"},
			"pg":    []string{fmt.Sprint(page)},
			"size":  []string{fmt.Sprint(len(word))},
		}, nil
}

// Query is a query to find words matching a frame.
type Query struct {
	Frame    string   `json:"frame,omitempty"`
	Page     int      `json:"page,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
}

func (q *Query) initialize() error {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Frame == "" {
		return inputerror.New("empty frame not allowed")
	}
	if !inputRE.MatchString(q.Frame) {
		return inputerror.New("inputs can only be letters or dots")
	}
	return nil
}

// NewQueryFromParams returns a query object from URL parameters
func NewQueryFromParams(values url.Values) (q Query, _ error) {
	q.Frame = values.Get("frame")
	pageStr := values.Get("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return q, errors.Wrapf(err, "page number %q", pageStr)
		}
		q.Page = p
	}
	var syns []string
	s1, s2 := values.Get("syn1"), values.Get("syn2")
	if s1 != "" {
		syns = append(syns, s1)
	}
	if s2 != "" {
		syns = append(syns, s2)
	}
	q.Synonyms = syns
	if err := q.initialize(); err != nil {
		return q, err
	}
	return q, nil
}

// Result is the result of a query
type Result struct {
	Query          *Query   `json:"query,omitempty"`          // the query for which results are provided
	SynonymMatches []string `json:"synonymMatches,omitempty"` // words that match both frame and synonyms
	Words          []string `json:"words"`                    // words found in current iteration, includes synonym matches
	NextPage       int      `json:"nextPage"`                 // next page to read, 0 means no more pages available
	TotalWords     int      `json:"totalWords"`               // total words matching frame
}

type pageResult struct {
	words      []string
	nextPage   int
	totalWords int
}

func (q *Query) readPage() (*pageResult, error) {
	if err := q.initialize(); err != nil {
		return nil, err
	}
	u, query, err := getURL(q.Frame, q.Page)
	if err != nil {
		return nil, err
	}

	doc, err := htmlplus.LoadURL(u, htmlplus.LoadOptions{
		Params: query,
	})
	if err != nil {
		return nil, err
	}

	wordCountDiv := doc.Find("div.word-criteria-heading")
	if wordCountDiv == nil {
		return nil, fmt.Errorf("no words found that match the frame")
	}
	matches := totalWordsRE.FindStringSubmatch(wordCountDiv.InnerText())
	if matches == nil {
		return nil, fmt.Errorf("internal error: could not find word count text")
	}
	totalWords, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("internal error: %w", err)
	}
	var ret []string
	nodes := doc.FindAll("div.word-results li.word a > span:first-child")
	for _, node := range nodes {
		spanText := node.InnerText()
		spanText = strings.ReplaceAll(spanText, " ", "")
		spanText = scoreRE.ReplaceAllString(spanText, "")
		ret = append(ret, strings.ToLower(spanText))
	}

	nextPage := q.Page
	if (nextPage-1)*infoPageSize >= totalWords {
		return nil, fmt.Errorf("read past last page")
	}

	if nextPage*infoPageSize >= totalWords {
		nextPage = 0
	} else {
		nextPage++
	}
	return &pageResult{
		words:      ret,
		nextPage:   nextPage,
		totalWords: totalWords,
	}, nil
}

func (q *Query) findWords() (*Result, error) {
	var finalResult Result
	for i := 0; i < infoPagesPerPage; i++ {
		result, err := q.readPage()
		if err != nil {
			return nil, err
		}
		finalResult.Words = append(finalResult.Words, result.words...)
		finalResult.NextPage = result.nextPage
		finalResult.TotalWords = result.totalWords
		if result.nextPage == 0 {
			break
		}
		q.Page = result.nextPage
	}
	return &finalResult, nil
}

type synonymsResult struct {
	words map[string]bool
	err   error
}

func (q *Query) findSynonyms(ch chan<- synonymsResult) {
	if len(q.Synonyms) == 0 {
		ch <- synonymsResult{words: map[string]bool{}}
		return
	}
	var wg sync.WaitGroup
	var l sync.Mutex
	matches := map[string]bool{}
	var finalErr error

	setMatches := func(m []*synonyms.Entry, err error) {
		l.Lock()
		defer l.Unlock()
		for _, s := range m {
			matches[s.Synonym] = true
		}
		if err != nil {
			finalErr = err
		}
	}
	for _, s := range q.Synonyms {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			sq := synonyms.Query{Word: word}
			res, err := sq.Run()
			setMatches(res.Entries, err)
		}(s)
	}
	wg.Wait()
	ch <- synonymsResult{
		words: matches,
		err:   finalErr,
	}
}

func (q *Query) Run() (*Result, error) {
	ch := make(chan synonymsResult, 1)
	q.findSynonyms(ch)

	result, err := q.findWords()
	if err != nil {
		return nil, err
	}

	synRes := <-ch
	err = synRes.err
	if err != nil {
		return nil, err
	}

	result.Query = q
	if len(q.Synonyms) > 0 {
		for _, word := range result.Words {
			if synRes.words[word] {
				result.SynonymMatches = append(result.SynonymMatches, word)
			}
		}
	}
	return result, nil
}
