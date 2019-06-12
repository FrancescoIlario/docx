package docx

import "github.com/antchfx/xmlquery"

func inpileNodes(nodes []*xmlquery.Node, parent *xmlquery.Node) []*xmlquery.Node {
	for idx, curr := range nodes {
		if idx != 0 {
			prev := nodes[idx-1]
			prev.NextSibling = curr
			curr.PrevSibling = prev
		}
		if idx != len(nodes)-1 {
			succ := nodes[idx+1]
			succ.PrevSibling = curr
			curr.NextSibling = succ
		}
	}
	return nodes
}

func setParent(parent *xmlquery.Node, nodes ...*xmlquery.Node) {
	for _, curr := range nodes {
		curr.Parent = parent
	}
}

func substituteNodeWithPile(del *xmlquery.Node, pile []*xmlquery.Node) {
	setParent(del.Parent, pile...)

	first, last := pile[0], pile[len(pile)-1]
	if prev := del.PrevSibling; prev != nil {
		prev.NextSibling = first
		first.PrevSibling = prev
	} else if parent := del.Parent; parent != nil {
		parent.FirstChild = first
		if parent.LastChild == del {
			parent.LastChild = last
		}
	}
}
