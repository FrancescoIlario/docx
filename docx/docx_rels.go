package docx

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/go-xmlfmt/xmlfmt"
)

// Relationship Relationship data structure
type Relationship struct {
	XMLName    xml.Name `xml:"Relationship"`
	ID         string   `xml:"Id,attr"`
	Type       string   `xml:"Type,attr"`
	Target     string   `xml:"Target,attr"`
	TargetMode string   `xml:"TargetMode,attr"`
}

// Relationships Relationships data structure
type Relationships struct {
	XMLName       xml.Name       `xml:"Relationships"`
	Relationships []Relationship `xml:"Relationship"`
}

func (r Relationship) String() string {
	return fmt.Sprintf("Id: %v, Type: %v, Target: %v, TargetMode: %v", r.ID, r.Type, r.Target, r.Target)
}

func (r Relationships) String() string {
	result := "{\n"
	for _, rel := range r.Relationships {
		result += "\t{ " + rel.String() + " },\n"
	}
	result += "}"

	return result
}

// PrintHyperlinks prints hyperlinks
func (d *ReplaceDocx) PrintHyperlinks() error {
	links := xmlfmt.FormatXML(d.links, "\t", " ")
	fmt.Println(links)

	var _rels Relationships
	if err := xml.Unmarshal([]byte(d.links), &_rels); err != nil {
		return err
	}

	fmt.Println(_rels)
	return nil
}

// GetHyperlinksRels returns the registered hyperlinks
func (d *ReplaceDocx) GetHyperlinksRels() ([]Relationship, error) {
	var _rels Relationships
	if err := xml.Unmarshal([]byte(d.links), &_rels); err != nil {
		return nil, err
	}

	var _returns []Relationship
	for _, rel := range _rels.Relationships {
		if strings.Contains(rel.Type, "hyperlink") {
			_returns = append(_returns, rel)
		}
	}
	return _returns, nil
}

func (d *ReplaceDocx) addHyperlinkRel(target string) {
	template := "<Relationship Id=\"{{rId}}\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink\" Target=\"{{target}}\" TargetMode=\"External\"/>"
	fmt.Println(template)
	strings.Replace(template, "{{target}}", target, -1)
}
