package htmlplus

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/andybalholm/cascadia"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

var (
	l        sync.Mutex
	compiled = map[string]cascadia.Sel{}
)

type Node struct {
	node *html.Node
}

func Wrap(node *html.Node) *Node {
	if node == nil {
		return nil
	}
	return &Node{node: node}
}

func (n *Node) Type() html.NodeType {
	return n.node.Type
}

func (n *Node) Tag() string {
	return n.node.Data
}

func (n *Node) parseSelector(selector string) cascadia.Sel {
	l.Lock()
	defer l.Unlock()
	s := compiled[selector]
	if s != nil {
		return s
	}
	s, err := cascadia.ParseWithPseudoElement(selector)
	if err != nil {
		panic(errors.Wrapf(err, "parse: '%s'", selector))
	}
	compiled[selector] = s
	return s
}

func (n *Node) Find(selector string) *Node {
	sel := n.parseSelector(selector)
	return Wrap(cascadia.Query(n.node, sel))
}

func (n *Node) FindAll(selector string) []*Node {
	sel := n.parseSelector(selector)
	nodes := cascadia.QueryAll(n.node, sel)
	var ret []*Node
	for _, node := range nodes {
		ret = append(ret, Wrap(node))
	}
	return ret
}

func (n *Node) AttributeValue(name string) string {
	for _, a := range n.node.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

func (n *Node) Children() []*Node {
	x := n.node.FirstChild
	var out []*Node
	for x != nil {
		out = append(out, Wrap(x))
		x = x.NextSibling
	}
	return out
}

func (n *Node) ElementChildren() []*Node {
	x := n.node.FirstChild
	var out []*Node
	for x != nil {
		if x.Type == html.ElementNode {
			out = append(out, Wrap(x))
		}
		x = x.NextSibling
	}
	return out
}

func (n *Node) writeHTML(w io.Writer) {
	node := n.node
	write := func(s string) {
		_, _ = w.Write([]byte(s))
	}

	switch node.Type {
	case html.DocumentNode, html.ElementNode:
		write("<")
		write(node.Data)
		for _, a := range node.Attr {
			write(" ")
			write(a.Key)
			write("=")
			write(`"`)
			write(html.EscapeString(a.Val))
			write(`"`)
		}
		write(">")

		child := node.FirstChild
		for child != nil {
			Wrap(child).writeHTML(w)
			child = child.NextSibling
		}

		write("</")
		write(node.Data)
		write(">")
	case html.TextNode:
		write(html.EscapeString(node.Data))
	}
}

func (n *Node) HTML() string {
	var b bytes.Buffer
	n.writeHTML(&b)
	return b.String()
}

func (n *Node) writeText(w io.Writer, isInner bool) {
	node := n.node
	write := func(s string) {
		_, _ = w.Write([]byte(s))
	}

	switch node.Type {
	case html.DocumentNode, html.ElementNode:
		child := node.FirstChild
		for child != nil {
			Wrap(child).writeText(w, true)
			child = child.NextSibling
		}
	case html.TextNode:
		if isInner {
			write(" ")
		}
		write(node.Data)
	}
}

func (n *Node) InnerText() string {
	var b bytes.Buffer
	n.writeText(&b, false)
	return strings.TrimSpace(b.String())
}

type Document struct {
	Node
}

func Load(r io.Reader) (*Document, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Document{Node: *Wrap(node)}, nil
}

// LoadOptions are options for loading a URL.
type LoadOptions struct {
	Method string     // HTTP method, defaults to GET
	Params url.Values // URL or form parameters based on method
}

// LoadURL loads a document from the supplied URL.
func LoadURL(u string, opts LoadOptions) (*Document, error) {
	if opts.Method == "" {
		opts.Method = http.MethodGet
	}

	var in io.Reader
	var ct string
	switch opts.Method {
	case http.MethodGet:
		if opts.Params != nil {
			u = fmt.Sprintf("%s?%s", u, opts.Params.Encode())
		}
	case http.MethodPost:
		in = strings.NewReader(opts.Params.Encode())
		ct = "application/x-www-form-urlencoded"
	default:
		return nil, fmt.Errorf("invalid HTTP method: %s", opts.Method)
	}
	req, err := http.NewRequest(opts.Method, u, in)
	if err != nil {
		return nil, err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s return status %d", opts.Method, u, res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	doc, err := Load(bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "read and parse HTML")
	}
	return doc, nil
}
