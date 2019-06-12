package docx

import (
	"log"
	"strings"

	"github.com/FrancescoIlario/docx/stringext"
	"github.com/FrancescoIlario/xmlquery"
)

//ReplaceTextWithHyperlink ReplaceTextWithHyperlink
func ReplaceTextWithHyperlink(filePath, outputFilePath, lookFor, link string) error {
	replaceDocx, err := ReadDocxFile(filePath)
	if err != nil {
		return err
	}
	defer replaceDocx.Close()

	doc, err := xmlquery.Parse(strings.NewReader(replaceDocx.content))
	if err != nil {
		return err
	}

	textNodes := xmlquery.Find(doc, `//w:p/w:r/w:t`)
	for _, textNode := range textNodes {
		replaceDocx.SubstituteRunWithHyperlinkWrtTarget(textNode, lookFor, link)
	}

	docXML := getDocumentXML(doc)
	replaceDocx.content = docXML

	return replaceDocx.Editable().WriteToFile(outputFilePath)
}

// SubstituteRunWithHyperlinkWrtTarget SubstituteRunWithHyperlinkWrtTarget
func (d *ReplaceDocx) SubstituteRunWithHyperlinkWrtTarget(chosenOne *xmlquery.Node, target, link string) {
	splits := stringext.SplitAfterWithSeparator(chosenOne.InnerText(), target)
	var nodes []*xmlquery.Node

	for _, split := range splits {
		node, err := d.getConfiguredNodeForSplit(chosenOne, split, target, link, nodes)
		if err == nil {
			nodes = append(nodes, node)
		}
	}

	for _, run := range nodes {
		log.Println(run.OutputXML(true))
	}
	parent := chosenOne.Parent
	pile := inpileNodes(nodes, parent.Parent)

	pileRoot := substituteNodeWithPile(parent, pile)
	d.content = fromNodeToRootOutputXML(pileRoot)
}

func (d *ReplaceDocx) getConfiguredNodeForSplit(chosenOne *xmlquery.Node, split, target, link string, runs []*xmlquery.Node) (*xmlquery.Node, error) {
	if split == target {
		return d.getConfiguredHyperlinkNode(chosenOne, target, link, runs)
	}
	return getTextRunFromRun(chosenOne.Parent, split)
}

func (d *ReplaceDocx) getConfiguredHyperlinkNode(chosenOne *xmlquery.Node, target, link string, runs []*xmlquery.Node) (*xmlquery.Node, error) {
	hyperlinkRelNode, err := d.GetHyperlinkOrAddForLink(link)
	if err != nil {
		return nil, err
	}

	rID := hyperlinkRelNode.SelectAttr("Id")
	node, err := NewHyperlinkNode(target, rID)
	if err != nil {
		return nil, err
	}
	return node, err
}

func getTextRunFromRun(run *xmlquery.Node, text string) (*xmlquery.Node, error) {
	newRunNode, err := docxRunDeepCopy(run)
	if err != nil {
		return nil, err
	}

	textNode := xmlquery.Find(newRunNode, `//w:t`)[0]
	textNode.FirstChild.Data = text

	return newRunNode, nil
}
