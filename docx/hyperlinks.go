package docx

import (
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/cbroglie/mustache"
)

func NewHyperlinkNode(text, rID string) (*xmlquery.Node, error) {
	documentHyperlinkXML, err := getDocumentHyperlinkXMLTemplate()
	if err != nil {
		return nil, err
	}

	templateData := map[string]string{
		"Text": text,
		"rID":  rID,
	}

	hyperlinkXML, err := mustache.Render(*documentHyperlinkXML, templateData)
	if err != nil {
		return nil, err
	}
	hyperlinkXML = html.UnescapeString(hyperlinkXML)

	log.Println(hyperlinkXML)
	hyperlinkXMLReader := strings.NewReader(hyperlinkXML)
	node, err := xmlquery.Parse(hyperlinkXMLReader)
	if err != nil {
		return nil, err
	}

	newHyperlinkNode := xmlquery.Find(node, "//w:hyperlink")[0]
	return newHyperlinkNode, nil
}

//GetHyperlinkOrAddForLink GetHyperlinkOrAddForLink
func (d *ReplaceDocx) GetHyperlinkOrAddForLink(link string) (*xmlquery.Node, error) {
	node, err := d.GetHyperlinkForLink(link)
	if err != nil {
		log.Println(err)
	}
	if node != nil {
		return node, nil
	}

	return d.AddHyperlinkForLink(link)
}

// GetHyperlinkForLink GetHyperlinkForLink
func (d *ReplaceDocx) GetHyperlinkForLink(link string) (*xmlquery.Node, error) {
	hyperlinkXMLReader := strings.NewReader(d.links)
	node, err := xmlquery.Parse(hyperlinkXMLReader)
	if err != nil {
		return nil, err
	}

	xpath := fmt.Sprintf(`/Relationships/Relationship[@Target='%s']`, link)
	hyperlinks := xmlquery.Find(node, xpath)
	if len(hyperlinks) != 0 {
		return hyperlinks[0], nil
	}

	return nil, nil
}

// AddHyperlinkForLink adds a new relationship for the link and returns the id
func (d *ReplaceDocx) AddHyperlinkForLink(link string) (*xmlquery.Node, error) {
	relationshipsRoot, err := d.getRelationships()
	if err != nil {
		return nil, err
	}

	newID := extractLHigherIDFromRelationships(relationshipsRoot) + 1
	rID := fmt.Sprintf("rId%v", newID)

	node, err := NewHyperlinkNode(link, rID)
	if err != nil {
		return nil, err
	}

	lastRelationship := relationshipsRoot.LastChild
	node.PrevSibling = lastRelationship
	lastRelationship.NextSibling = node

	return node, nil
}

func (d *ReplaceDocx) getRelationships() (*xmlquery.Node, error) {
	hyperlinkXMLReader := strings.NewReader(d.links)
	node, err := xmlquery.Parse(hyperlinkXMLReader)
	if err != nil {
		return nil, err
	}

	relationshipsRoot := xmlquery.Find(node, `/Relationships`)[0]
	return relationshipsRoot, nil
}

func extractLHigherIDFromRelationships(relationshipsRoot *xmlquery.Node) int64 {
	var lastID int64
	lastID = 0
	iteratorNode := relationshipsRoot.FirstChild
	for {
		if iteratorNode == nil {
			break
		}

		for _, attr := range iteratorNode.Attr {
			if attr.Name.Local == "Id" {
				id := attr.Value[len("rId"):]
				if parsedID, err := strconv.ParseInt(id, 10, 32); err == nil && lastID < parsedID {
					lastID = parsedID
				}
			}
		}

		iteratorNode = iteratorNode.NextSibling
	}

	return lastID
}
