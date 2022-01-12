package crawler

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
)

type Crawler struct {
	Root *html.Node
}

func NewCrawler(r io.Reader) (*Crawler, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return &Crawler{Root: doc}, nil
}

func (c *Crawler) findElement(node *html.Node, f func(n *html.Node) *html.Node) (*Crawler, error) {
	cr := f(node)
	if cr != nil {
		return &Crawler{cr}, nil
	}

	var result *Crawler
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result, _ = c.findElement(child, f)
		if result != nil {
			break
		}
	}

	if result != nil {
		return result, nil
	}

	return nil, errors.New("can't find an element")
}

func (c *Crawler) FindElementByID(id string) (*Crawler, error) {
	result, err := c.findElement(c.Root, func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && len(n.Attr) > 0 {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					return n
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't find an element with ID: %s", id))
	}

	return result, nil
}

func (c *Crawler) FindElementByTag(tag string) (*Crawler, error) {
	result, err := c.findElement(c.Root, func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == tag {
			return n
		}

		return nil
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't find an element with tag: %s", tag))
	}

	return result, nil
}

func (c *Crawler) Children() ([]*Crawler, error) {
	var children []*Crawler
	for child := c.Root.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			children = append(children, &Crawler{child})
		}
	}

	if len(children) == 0 {
		return nil, errors.New("element doesn't have any children")
	}

	return children, nil
}

func (c *Crawler) InnerText() (string, error) {
	for child := c.Root.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			return child.Data, nil
		}
	}

	return "", errors.New("inner text is empty")
}
