package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/html"
)

var (
	box        = "myBox"
	classes    = []string{"myCheck", "myNote"}
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

	walkNodes(doc, 0, 1)

	err = html.Render(os.Stdout, doc)
	if err != nil {
		fmt.Println("failed to render html, ", err)
		return 1
	}

	return 0
}

func walkNodes(n *html.Node, boxCounter int, elemCounter int) (int, int) {
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
			fmt.Printf("Attr: key=%v, val=%v\n", a.Key, a.Val)
		}

		if contains {
			n.Attr = append(n.Attr, html.Attribute{Key: "id", Val: fmt.Sprintf("%s-%s-%d", getBoxID(boxCounter), cName, elemCounter)})
			elemCounter++
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		boxCounter, elemCounter = walkNodes(c, boxCounter, elemCounter)
	}

	return boxCounter, elemCounter
}

func getBoxID(counter int) string {
	return fmt.Sprintf("box%v", counter)
}
