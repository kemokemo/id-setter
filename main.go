package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/html"
)

func main() {
	os.Exit(run())
}

func run() int {
	b, err := ioutil.ReadFile("source.html")
	if err != nil {
		fmt.Println("failed to load source.html, ", err)
		return 1
	}

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		fmt.Println("failed to parse html, ", err)
		return 1
	}

	walkNodes(doc)

	return 0
}

func walkNodes(n *html.Node) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			fmt.Printf("Attr: namespace=%v, key=%v, val=%v\n", a.Namespace, a.Key, a.Val)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNodes(c)
	}
}
