package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/FrancescoIlario/docx/docx"
	"github.com/antchfx/xmlquery"
)

const (
	documentTemplatePath = `templates/docx_document.xml`
	runTemplatePath      = `templates/run.xml`
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
	if channels != nil && len(channels) > 0 {
		fmt.Println("runs:")
		for _, channel := range channels {
			fmt.Println(channel.OutputXML(true))
			textNodes := xmlquery.Find(channel, `/w:t`)
			if len(textNodes) == 1 {
				// does not substitute
				err = replaceDocx.SubstituteRunWithHyperlinkWrtTarget(textNodes[0], "This", "Allegato.txt")
				log.Printf("content after substitution\n%v\n", replaceDocx.Content)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	replaceDocx.Content = getDocumentXML(doc)
	err = replaceDocx.Editable().WriteToFile("data/docx/HyperlinkedDocument.docx")
	panicIf(err)
}

func getDocumentXML(document *xmlquery.Node) string {
	var doc string
	iteratorNode := document.FirstChild
	for iteratorNode != nil {
		doc += iteratorNode.OutputXML(true)
		iteratorNode = iteratorNode.NextSibling
	}
	return doc
}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
