package docx

import "github.com/antchfx/xmlquery"

// XMLNode XMLNode
type XMLNode xmlquery.Node

func getDocumentXML(document *xmlquery.Node) string {
	var doc string

	iteratorNode := document.FirstChild
	for iteratorNode != nil {
		doc += iteratorNode.OutputXML(true)
		iteratorNode = iteratorNode.NextSibling
	}

	return doc
}
