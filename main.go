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
	s := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
<channel>
  <title>W3Schools Home Page</title>
  <link>https://www.w3schools.com</link>
  <description>Free web building tutorials</description>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
  </item>
  <item>
    <title>XML Tutorial</title>
    <link>https://www.w3schools.com/xml</link>
    <description>New XML tutorial on W3Schools</description>
  </item>
</channel>
</rss>`

	doc, err := xmlquery.Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	channel := xmlquery.FindOne(doc, "//channel")
	if n := channel.SelectElement("title"); n != nil {
		fmt.Printf("title: %s\n", n.OutputXML(true))
	}
	if n := channel.SelectElement("link"); n != nil {
		fmt.Printf("link: %s\n", n.Parent.OutputXML(true))
	}
	for i, n := range xmlquery.Find(doc, "//item/title") {
		fmt.Printf("#%d %s\n", i, n.InnerText())
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
