package docx

import (
	"log"

	"github.com/antchfx/xmlquery"
)

func substituteNode(toBeRemoved, toBeInserted *xmlquery.Node) {
	prevSibling := toBeRemoved.PrevSibling

	toBeInserted.NextSibling = toBeRemoved.NextSibling
	prevSibling.NextSibling = toBeInserted
}

func substituteNodeWithNodes(toBeRemoved *xmlquery.Node, toBeInserted []*xmlquery.Node) *xmlquery.Node {
	if len(toBeInserted) == 0 {
		return toBeRemoved
	}
	first, last := toBeInserted[0], toBeInserted[len(toBeInserted)-1]
	prevSibling := toBeRemoved.PrevSibling

	for idx, el := range toBeInserted {
		log.Printf("%v - %v\n", idx, el.OutputXML(true))
	}

	if prevSibling != nil {
		log.Printf("Previous sibling\n%s\n", prevSibling.OutputXML(true))
		prevSibling.NextSibling = first
		first.PrevSibling = prevSibling
	} else {
		if toBeRemoved.Parent != nil {
			// this should be wrong -> rev. why?
			log.Printf("Parent\n%s\n", toBeRemoved.Parent.OutputXML(true))
			toBeRemoved.Parent.FirstChild = first
		} else {
			log.Printf("Node has neither a sibling nor a parent.... It seems unbelievable: %v\n", *toBeRemoved)
		}
	}

	last.NextSibling = toBeRemoved.NextSibling

	log.Printf("after substitution\n%s\n", fromNodeToRootOutputXML(first))
	return first
}
