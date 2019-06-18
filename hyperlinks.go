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

// NewHyperlinkRelNode  todo
func NewHyperlinkRelNode(target, rID string) (*xmlquery.Node, error) {
	documentHyperlinkXML, err := getDocumentHyperlinkRelXMLTemplate()
	if err != nil {
		return nil, err
	}

	templateData := map[string]string{
		"Target": target,
		"rID":    rID,
	}

	hyperlinkXML, err := mustache.Render(*documentHyperlinkXML, templateData)
	if err != nil {
		return nil, err
	}
	hyperlinkXML = html.UnescapeString(hyperlinkXML)

	log.Printf("hyperlink xml\n%s\n", hyperlinkXML)
	hyperlinkXMLReader := strings.NewReader(hyperlinkXML)
	node, err := xmlquery.Parse(hyperlinkXMLReader)
	if err != nil {
		return nil, err
	}
	xpath := fmt.Sprintf(`//Relationship[@Target='%s']`, target)
	newHyperlinkNode := xmlquery.Find(node, xpath)[0]

	return newHyperlinkNode, nil
}

// NewHyperlinkNode todo
func NewHyperlinkNode(text, rID string) (*xmlquery.Node, error) {
	documentHyperlinkXML, err := getDocumentHyperlinkXMLTemplate()
	if err != nil {
		return nil, err
	}

	var attrs string
	if needsAttrsStringXMLSpacePreserve(text, []string{}) {
		attrs = " " + spacePreserveAttr
	}

	templateData := map[string]string{
		"Attrs": attrs,
		"Text":  text,
		"rID":   rID,
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

// GetOrAddHyperlinkForLink looks for an existing relationship for the provided link and if the research
// is unsuccessfull, it creates and returns a new relationship describing the provided link
func (d *ReplaceDocx) GetOrAddHyperlinkForLink(link string) (*xmlquery.Node, error) {
	node, err := d.GetHyperlinkForLink(link)
	if err != nil {
		log.Println(err)
	}
	if node != nil {
		return node, nil
	}

	return d.AddHyperlinkForLink(link)
}

// GetHyperlinkForLink looks in document.rel.xml for a relationship that represents an Hyperlink to the
// given link. If a matching relationship is found it is returned, otherwise is returned the nil value
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

// AddHyperlinkForLink adds a new relationship for the link and returns the it
func (d *ReplaceDocx) AddHyperlinkForLink(link string) (*xmlquery.Node, error) {
	relationshipsRoot, err := d.GetRelationships()
	if err != nil {
		return nil, err
	}

	newID := extractLHigherIDFromRelationships(relationshipsRoot) + 1
	rID := fmt.Sprintf("rId%v", newID)

	node, err := NewHyperlinkRelNode(link, rID)
	if err != nil {
		return nil, err
	}

	lastRelationship := relationshipsRoot.LastChild
	node.PrevSibling = lastRelationship
	lastRelationship.NextSibling = node

	d.links = fromNodeToRootOutputXML(relationshipsRoot)
	return node, nil
}

func fromNodeToRootOutputXML(node *xmlquery.Node) string {
	var doc string

	root := node
	visited := []*xmlquery.Node{}
	for {
		for _, vis := range visited {
			if vis == root {
				log.Println("Loop found")
				for idx, visLog := range visited {
					log.Printf("%v - %v\n", idx, visLog)
				}
				panic("loop found")
			}
		}
		if root.Parent != nil {
			visited = append(visited, root)
			root = root.Parent
		} else {
			break
		}
	}

	visited = []*xmlquery.Node{}
	for {
		for _, vis := range visited {
			if vis == root {
				log.Println("Loop found")
				for idx, visLog := range visited {
					log.Printf("%v - %v\n", idx, visLog)
				}
				panic("loop found")
			}
		}
		if root.PrevSibling != nil {
			visited = append(visited, root)
			root = root.PrevSibling
		} else {
			break
		}
	}

	visited = []*xmlquery.Node{}
	root = root.FirstChild
	for {
		if root == nil {
			break
		}

		for _, vis := range visited {
			if vis == root {
				log.Println("Loop found")
				for idx, visLog := range visited {
					log.Printf("%v - %v\n", idx, visLog)
				}
				panic("loop found")
			}
		}

		doc += html.UnescapeString(root.OutputXML(true))
		visited = append(visited, root)
		root = root.NextSibling
	}

	return html.UnescapeString(doc)
}

// GetRelationships returns the relationships of the docx file
func (d *ReplaceDocx) GetRelationships() (*xmlquery.Node, error) {
	log.Printf("links\n%s\n", d.links)
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
