package docx

import (
	"errors"
	"log"

	"github.com/antchfx/xmlquery"
)

func substituteNode(toBeRemoved, toBeInserted *xmlquery.Node) {
	prevSibling := toBeRemoved.PrevSibling

	toBeInserted.NextSibling = toBeRemoved.NextSibling
	prevSibling.NextSibling = toBeInserted
}

func substituteNodeWithNodes(toBeRemoved *xmlquery.Node, toBeInserted []*xmlquery.Node) error {
	if len(toBeInserted) == 0 {
		return errors.New("empty list of xml nodes passed")
	}

	first, last := toBeInserted[0], toBeInserted[len(toBeInserted)-1]

	last.NextSibling = toBeRemoved.NextSibling

	prevSibling := toBeRemoved.PrevSibling
	
	if prevSibling != nil {
		log.Printf("Previous sibling\n%s\n", prevSibling.OutputXML(true))
		prevSibling.NextSibling = first
	} else {
		if toBeRemoved.Parent != nil {
			log.Printf("Parent\n%s\n", toBeRemoved.Parent.OutputXML(true))
			toBeRemoved.Parent.FirstChild = first
		} else {
			log.Printf("Node has neither a sibling nor a parent.... It seems unbelievable: %v\n", *toBeRemoved)
		}
	}

	return nil
}
