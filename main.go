package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/FrancescoIlario/docx/docx"
	"github.com/antchfx/xmlquery"
)

type void struct{}

func main() {
	replaceDocx, err := docx.ReadDocxFile("data/docx/TestDocument.docx")
	panicIf(err)
	defer replaceDocx.Close()

	content := replaceDocx.Content
	doc, err := xmlquery.Parse(strings.NewReader(content))
	if err != nil {
		panic(err)
	}

	channels := xmlquery.Find(doc, `//w:p/w:r`)
	runs := make(map[string]void)
	var member void
	var chosenOne *xmlquery.Node
	if channels != nil && len(channels) > 0 {
		fmt.Println("runs:")
		for _, channel := range channels {
			fmt.Println(channel.Parent.OutputXML(true))

			if n := channel.SelectElement("title"); n != nil {
				runs[n.OutputXML(true)] = member
				fmt.Printf("runs: %s\n", n.OutputXML(true))
			}
		}

		chosenOne = channels[0].Parent
		fmt.Printf("Using %s\n", chosenOne.OutputXML(true))
	}

	if chosenOne != nil {
		// prevSibling := chosenOne.PrevSibling

		// xmlquery.Parse()
	}

}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
