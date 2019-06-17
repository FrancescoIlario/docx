package docx

import (
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"strings"

	assets "github.com/FrancescoIlario/docx/assets/templates"
	"github.com/antchfx/xmlquery"
	"github.com/cbroglie/mustache"
)

const (
	documentTemplatePath     = `docx_document.xml`
	runTemplatePath          = `run.xml`
	textTemplatePath         = `text.xml`
	hyperlinkTemplatePath    = "hyperlink.xml"
	hyperlinkRelTemplatePath = "hyperlink_rel.xml"
	spacePreserveAttr        = `xml:space="preserve"`

	// documentTemplatePath     = `templates/docx_document.xml`
	// runTemplatePath          = `templates/run.xml`
	// textTemplatePath         = `templates/text.xml`
	// hyperlinkTemplatePath    = "templates/hyperlink.xml"
	// hyperlinkRelTemplatePath = "templates/hyperlink_rel.xml"

)

func parseFile(filepath string) (*string, error) {
	file, err := assets.Assets.Open(filepath)
	if err != nil {
		panic(err)
	}

	// data, err := ioutil.ReadFile(filepath)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	stringData := string(data)
	return &stringData, nil
}

func getHyperlinkXMLTemplate() (*string, error) {
	return parseFile(hyperlinkTemplatePath)
}

func getHyperlinkRelXMLTemplate() (*string, error) {
	return parseFile(hyperlinkRelTemplatePath)
}

func getRunXMLTemplate() (*string, error) {
	return parseFile(runTemplatePath)
}

func getTextXMLTemplate() (*string, error) {
	return parseFile(textTemplatePath)
}

func getDocumentTextXMLTemplate() (*string, error) {
	documentTemplate, err := getDocumentXMLTemplate()
	if err != nil {
		return nil, err
	}

	runTemplate, err := getTextXMLTemplate()
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

func getDocumentHyperlinkRelXMLTemplate() (*string, error) {
	documentTemplate, err := getDocumentXMLTemplate()
	if err != nil {
		return nil, err
	}

	hyperlinkTemplate, err := getHyperlinkRelXMLTemplate()
	if err != nil {
		return nil, err
	}

	documentTemplateData := map[string]string{
		"Template": *hyperlinkTemplate,
	}

	documentHyperlinkXML, err := mustache.Render(*documentTemplate, documentTemplateData)
	if err != nil {
		return nil, err
	}

	return &documentHyperlinkXML, nil
}

func getDocumentHyperlinkXMLTemplate() (*string, error) {
	documentTemplate, err := getDocumentXMLTemplate()
	if err != nil {
		return nil, err
	}

	hyperlinkTemplate, err := getHyperlinkXMLTemplate()
	if err != nil {
		return nil, err
	}

	documentTemplateData := map[string]string{
		"Template": *hyperlinkTemplate,
	}

	documentHyperlinkXML, err := mustache.Render(*documentTemplate, documentTemplateData)
	if err != nil {
		return nil, err
	}

	return &documentHyperlinkXML, nil
}

func getDocumentXMLTemplate() (*string, error) {
	data, err := parseFile(documentTemplatePath)
	if err != nil {
		return nil, err
	}
	return data, nil
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

func getTextNode(text string, attrs ...string) (*xmlquery.Node, error) {
	documentTemplate, err := getDocumentXMLTemplate()
	if err != nil {
		return nil, err
	}

	pTextXMLTemplate, err := getTextXMLTemplate()
	if err != nil {
		return nil, err
	}
	textXMLTemplate := html.UnescapeString(*pTextXMLTemplate)

	documentTemplateData := map[string]string{
		"Template": textXMLTemplate,
	}

	documentRunXML, err := mustache.Render(*documentTemplate, documentTemplateData)
	if err != nil {
		return nil, err
	}
	documentRunXML = html.UnescapeString(documentRunXML)

	var attrsStr string
	if needsAttrsStringXMLSpacePreserve(text, attrs) {
		attrsStr = strings.Join([]string{attrsStr, spacePreserveAttr}, " ")
	}

	documentTemplateData = map[string]string{
		"Text":  text,
		"Attrs": attrsStr,
	}

	textDocumentXML, err := mustache.Render(documentRunXML, documentTemplateData)
	if err != nil {
		return nil, err
	}
	textDocumentXML = html.UnescapeString(textDocumentXML)

	txtTemplateReader := strings.NewReader(textDocumentXML)
	documentNewTxtNode, err := xmlquery.Parse(txtTemplateReader)
	if err != nil {
		return nil, err
	}

	newTxtNode := xmlquery.Find(documentNewTxtNode, `//w:t`)[0]
	return newTxtNode, nil
}

func needsAttrsXMLSpacePreserve(text string, attrs []xml.Attr) bool {
	var attrsStrs []string
	for _, attr := range attrs {
		attrStr := fmt.Sprintf(`%s:%s="%s"`, attr.Name.Space, attr.Name.Local, attr.Value)
		attrsStrs = append(attrsStrs, attrStr)
	}

	return needsAttrsStringXMLSpacePreserve(text, attrsStrs)
}

func needsAttrsStringXMLSpacePreserve(text string, attrs []string) bool {
	attrsStr := strings.Join(attrs, " ")
	return (strings.HasPrefix(text, " ") || strings.HasSuffix(text, " ")) &&
		!strings.Contains(attrsStr, spacePreserveAttr)
}
