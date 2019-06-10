package docx

import (
	"io/ioutil"

	"github.com/cbroglie/mustache"
)

const (
	documentTemplatePath  = `templates/docx_document.xml`
	runTemplatePath       = `templates/run.xml`
	hyperlinkTemplatePath = "templates/hyperlink.xml"
)

func parseFile(filepath string) (*string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	stringData := string(data)
	return &stringData, nil
}

func getHyperlinkXMLTemplate() (*string, error) {
	return parseFile(hyperlinkTemplatePath)
}

func getRunXMLTemplate() (*string, error) {
	return parseFile(runTemplatePath)
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
