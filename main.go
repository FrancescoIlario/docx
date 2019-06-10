package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/cbroglie/mustache"

	"github.com/FrancescoIlario/docx/docx"
	"github.com/antchfx/xmlquery"
)

const (
	documentTemplatePath = `templates/docx_document.xml`
	runTemplatePath      = `templates/run.xml`
	hyperlinkTemplate    = `<w:hyperlink r:id="{{ID}}"><w:r><w:rPr><w:rStyle w:val="InternetLink"/></w:rPr><w:t>{{Text}}</w:t></w:r></w:hyperlink>`
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
			fmt.Println(channel.OutputXML(true))

			if n := channel.SelectElement("title"); n != nil {
				runs[n.OutputXML(true)] = member
				fmt.Printf("runs: %s\n", n.OutputXML(true))
			}
		}

		chosenOne = channels[0]
		fmt.Printf("Using %s\n", chosenOne.OutputXML(true))
	}

	if chosenOne != nil {
		node, err := getRunNode("AAAAAAAAAAAAAAAAAAAAAAAAAAA")
		panicIf(err)

		substituteNode(chosenOne, node)

		fmt.Println("\n\nfinal result: ")
		printDocument(doc)
	}
}

func printDocument(document *xmlquery.Node) {
	iteratorNode := document.FirstChild
	for iteratorNode != nil {
		fmt.Println(iteratorNode.OutputXML(true))
		iteratorNode = iteratorNode.NextSibling
	}
}

func substituteNode(toBeRemoved, toBeInserted *xmlquery.Node) {
	prevSibling := toBeRemoved.PrevSibling
	toBeInserted.NextSibling = toBeRemoved.NextSibling
	prevSibling.NextSibling = toBeInserted
}

func getRunNode(text string) (*xmlquery.Node, error) {
	documentRunXML, err := getDocumentRunXMLTemplate()
	if err != nil {
		return nil, err
	}

	templateData := map[string]string{
		"Text": text,
	}

	runXML, err := mustache.Render(*documentRunXML, templateData)
	if err != nil {
		return nil, err
	}
	runXMLReader := strings.NewReader(runXML)
	node, err := xmlquery.Parse(runXMLReader)
	if err != nil {
		return nil, err
	}

	newRunNode := xmlquery.Find(node, "//w:r")[0]
	return newRunNode, nil
}

func parseFile(filepath string) (*string, error) {
	data, err := ioutil.ReadFile(runTemplatePath)
	if err != nil {
		return nil, err
	}

	stringData := string(data)
	return &stringData, nil
}

func getRunXMLTemplate() (*string, error) {
	data, err := parseFile(runTemplatePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getDocumentRunXMLTemplate() (*string, error) {
	documentTemplate, err := getDocumentXMLTemplate()
	if err != nil {
		return nil, err
	}

	runTemplate, err := getRunXMLTemplate()
	if err != nil {
		return nil, err
	}

	documentTemplateData := map[string]string{
		"Template": *runTemplate,
	}

	documentRunXML, err := mustache.Render(*documentTemplate, documentTemplateData)
	if err != nil {
		return nil, err
	}

	return &documentRunXML, nil
}

func getDocumentXMLTemplate() (*string, error) {
	data, err := parseFile(documentTemplatePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
