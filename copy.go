package docx

import (
	"html"
	"strings"

	"github.com/FrancescoIlario/xmlquery"
	"github.com/cbroglie/mustache"
)

func docxRunDeepCopy(src *xmlquery.Node) (*xmlquery.Node, error) {
	docTemplate, err := getDocumentXMLTemplate()
	if err != nil {
		return nil, err
	}

	xmlCode := src.OutputXML(true)
	docXMLCode, err := mustache.Render(*docTemplate, map[string]string{"Template": xmlCode})
	if err != nil {
		return nil, err
	}
	docXMLCode = html.UnescapeString(docXMLCode)

	doc, err := xmlquery.Parse(strings.NewReader(docXMLCode))
	if err != nil {
		return nil, err
	}

	out := xmlquery.Find(doc, "//w:r")[0]
	return out, nil
}
