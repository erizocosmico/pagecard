package content

import (
	"net/http"

	"golang.org/x/net/html"
)

// Meta represents a key-value metatag on the webpage.
type Meta struct {
	Name  string
	Value string
}

var client = &http.Client{}

// Read scans the page content at the given URL and returns a list of its metatags.
func Read(url string) ([]*Meta, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	metatags := extractMetatags(node)
	return metatagsToMetaList(metatags), nil
}

func metatagsToMetaList(metatags []*html.Node) []*Meta {
	var (
		result []*Meta
		meta   *Meta
	)

	for _, m := range metatags {
		meta = &Meta{}
		for _, attr := range m.Attr {
			switch attr.Key {
			case "property", "name":
				meta.Name = attr.Val
			case "content":
				meta.Value = attr.Val
			}
		}

		if meta.Name != "" && meta.Value != "" {
			result = append(result, meta)
		}
	}

	return result
}

func extractMetatags(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "head" {
		return findMetaNodesInHead(n)
	}

	if n.Type != html.ElementNode && n.Type != html.DocumentNode {
		return nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes := extractMetatags(c)
		if len(nodes) > 0 {
			return nodes
		}
	}

	return nil
}

func findMetaNodesInHead(n *html.Node) []*html.Node {
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "meta" {
			nodes = append(nodes, c)
		}
	}
	return nodes
}
