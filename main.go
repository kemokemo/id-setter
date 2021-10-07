package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/html"
)

var (
	box        = "myBox"
	classes    = []string{"myCheck", "myNote", "mySelect"}
	classesMap = make(map[string]struct{})
)

func init() {
	for _, v := range classes {
		classesMap[v] = struct{}{}
	}
}

func main() {
	os.Exit(run())
}

func run() int {
	source := "source.html"

	b, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Println("failed to load source.html, ", err)
		return 1
	}

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		fmt.Println("failed to parse html, ", err)
		return 1
	}

	hash := sha256.New()
	docURI := fmt.Sprintf("%x", hash.Sum([]byte(source)))

	walkNodes(doc, 0, 1, docURI)

	err = html.Render(os.Stdout, doc)
	if err != nil {
		fmt.Println("failed to render html, ", err)
		return 1
	}

	return 0
}

func walkNodes(n *html.Node, boxCounter int, elemCounter int, docURI string) (int, int) {
	if n.Data == "body" {
		n.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "div",
			Attr: []html.Attribute{
				{Key: "id", Val: "myDocURI"},
				{Key: "value", Val: docURI},
			}},
		)
	}

	if n.Type == html.ElementNode {
		var contains bool
		var cName string

		for _, a := range n.Attr {
			if a.Key != "class" {
				continue
			}

			if a.Val == box {
				boxCounter++
				n.Attr = append(n.Attr, html.Attribute{Key: "id", Val: getBoxID(boxCounter)})
				continue
			}

			_, ok := classesMap[a.Val]
			if ok {
				contains = true
				cName = a.Val
			}
		}

		if contains {
			id := fmt.Sprintf("%s-%s-%d", getBoxID(boxCounter), cName, elemCounter)
			n.Attr = append(n.Attr, html.Attribute{Key: "id", Val: id})
			elemCounter++
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		boxCounter, elemCounter = walkNodes(c, boxCounter, elemCounter, docURI)
	}

	return boxCounter, elemCounter
}

func getBoxID(counter int) string {
	return fmt.Sprintf("box%v", counter)
}
