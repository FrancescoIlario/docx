package docx

import (
	"log"
	"strings"

	"github.com/FrancescoIlario/docx/stringext"
	"github.com/antchfx/xmlquery"
	"github.com/cbroglie/mustache"
)

// SubstituteRunWithHyperlinkWrtTarget SubstituteRunWithHyperlinkWrtTarget
func (d *ReplaceDocx) SubstituteRunWithHyperlinkWrtTarget(chosenOne *xmlquery.Node, target, link string) error {
	splits := stringext.SplitAfterWithSeparator(chosenOne.InnerText(), target)
	var runs []*xmlquery.Node

	for _, split := range splits {
		if split == target {
			if node, err := d.GetHyperlinkOrAddForLink(link); err == nil {
				if len(runs) > 0 {
					lastRun := runs[len(runs)-1]
					node.PrevSibling = lastRun
					lastRun.LastChild = node
				}

				runs = append(runs, node)
			} else {
				log.Println(err)
			}
		} else {
			runNode := *chosenOne
			runNode.Data = split // TODO this is wrong

			if runsSize := len(runs); runsSize > 0 {
				lastRun := runs[len(runs)-1]
				runNode.PrevSibling = lastRun
				lastRun.LastChild = &runNode
			}

			runs = append(runs, &runNode)
		}
	}

	substituteNodeWithNodes(chosenOne.Parent, runs)

	return nil
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
