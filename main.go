package main

import (
	"flag"
	"fmt"
	"log"

	"strings"

	"github.com/FrancescoIlario/docx/docx"
	"github.com/antchfx/xmlquery"
)

func main() {
	replaceDocx, err := docx.ReadDocxFile("data/docx/TestDocument.docx")
	panicIf(err)
	defer replaceDocx.Close()

	content := replaceDocx.Content()
	doc, err := xmlquery.Parse(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	channel := xmlquery.FindOne(doc, "//w:p//w:r")
	fmt.Printf("title: %s\n", channel.OutputXML(true))

	if n := channel.SelectElement("title"); n != nil {
		fmt.Printf("title: %s\n", n.OutputXML(true))
	}
}

var (
	lookFor = flag.String("LookFor", "This", "The substring to look")
)

func mainO() {
	var replaceDocx *docx.ReplaceDocx
	var err error

	replaceDocx, err = docx.ReadDocxFile("data/docx/TestDocument.docx")
	panicIf(err)
	defer replaceDocx.Close()

	text := replaceDocx.GetText()
	fmt.Println(text)

	document := replaceDocx.Editable()
	paragraphs, err := document.ExtractParagraphs()
	panicIf(err)

	for _, paragraph := range paragraphs {
		fmt.Println(paragraph.Text(true))
	}

	formattedContent := replaceDocx.FormattedContent("")
	fmt.Println(formattedContent)
}

func panicIf(err error) {
	if err != nil {
		if errData, ok := err.(*docx.WrongXMLSlice); ok {
			log.Printf("%s\n", errData.XMLSlice)
		}
		log.Panicln(err)
	}
}
