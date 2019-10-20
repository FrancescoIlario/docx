package docx

import (
	"errors"
	"html"
	"io"
	"log"
	"strings"

	"github.com/antchfx/xmlquery"
)

// AddInternetLinkStyleIfMissing looks into the `word/styles.xml` file for
// InternetLink style, and creates the entry if it does not exist
func (d *ReplaceDocx) AddInternetLinkStyleIfMissing() error {
	doc, err := xmlquery.Parse(strings.NewReader(d.styles))
	if err != nil {
		return err
	}

	if !existsInternetLinkStyle(doc) {
		log.Println("Internet link style not found, adding it")
		return d.addInternetLinkStyle(doc)
	}

	log.Println("Internet link style is yet present")
	return nil
}

func existsInternetLinkStyle(doc *xmlquery.Node) bool {
	textNodes := xmlquery.Find(doc, `//w:styles/w:style[w:styleId="InternetLink"]`)
	return len(textNodes) > 0
}

func (d *ReplaceDocx) addInternetLinkStyle(doc *xmlquery.Node) error {
	nodes := xmlquery.Find(doc, `//w:styles`)
	if len(nodes) == 0 {
		return errors.New("word/styles.xml presents no w:styles element")
	}

	internetLinkStyleReader, err := internetLinkStyleXML()
	if err != nil {
		return err
	}

	internetLinkStyleDocument, err := xmlquery.Parse(internetLinkStyleReader)
	if err != nil {
		return err
	}

	internetLinkStyleNodes := xmlquery.Find(internetLinkStyleDocument, `//w:style[@w:styleId="InternetLink"]`)
	if len(internetLinkStyleNodes) == 0 {
		return errors.New("\"//w:style[@w:styleId=\"InternetLink\"]\" not found in generated document")
	}

	nodes[0].LastChild.NextSibling = internetLinkStyleNodes[0]
	d.styles = fromNodeToRootOutputXML(doc)

	return nil
}

func internetLinkStyleXML() (io.Reader, error) {
	internetLinkStyle, err := getInternetLinkStyleXMLTemplate()
	if err != nil {
		return nil, err
	}

	internetLinkStyle = html.UnescapeString(internetLinkStyle)
	return strings.NewReader(internetLinkStyle), nil
}
