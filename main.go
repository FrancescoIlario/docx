package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/FrancescoIlario/docx/docx"
)

var (
	lookFor = flag.String("LookFor", "This", "The substring to look")
)

func main() {
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
