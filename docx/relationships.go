package docx

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/go-xmlfmt/xmlfmt"
)

const (
	closedRelationshipsTag        = `</Relationships>`
	hyperlinkType                 = `http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink`
	relationshipHyperlinkTemplate = `<Relationship Id="{{rID}}" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="{{target}}" TargetMode="External"/>`
)

// Relationship Relationship data structure
type Relationship struct {
	XMLName    xml.Name `xml:"Relationship"`
	InnerXML   string   `xml:",innerxml"`
	ID         string   `xml:"Id,attr"`
	Type       string   `xml:"Type,attr"`
	Target     string   `xml:"Target,attr"`
	TargetMode string   `xml:"TargetMode,attr"`
}

// Relationships Relationships data structure
type Relationships struct {
	XMLName       xml.Name       `xml:"Relationships"`
	InnerXML      string         `xml:",innerxml"`
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

// GetRels returns the registered relationships
func (d *ReplaceDocx) GetRels() ([]Relationship, error) {
	var _rels Relationships
	if err := xml.Unmarshal([]byte(d.links), &_rels); err != nil {
		return nil, err
	}

	return _rels.Relationships, nil
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

// CreateHyperlinkRel CreateHyperlinkRel
func createHyperlinkRel(rID, target string) (string, error) {
	hyperlinkData := map[string]string{
		"rID":    rID,
		"target": target,
	}

	hyperlink, err := mustache.Render(relationshipHyperlinkTemplate, hyperlinkData)
	return hyperlink, err
}

//AddHyperlinkRel returns the id of the inserted hyperlink rel
func (d *ReplaceDocx) AddHyperlinkRel(link string) (*string, error) {
	hRelID, err := d.GetHyperlinkRel(link)
	if err != nil {
		return nil, err
	}
	if hRelID != nil {
		return hRelID, nil
	}

	rID, err := d.newRelsID()
	if err != nil {
		return nil, err
	}

	hyperlink, err := createHyperlinkRel(*rID, link)
	if err != nil {
		return nil, err
	}

	// inserting new link
	d.links = strings.Replace(d.links, closedRelationshipsTag, hyperlink+closedRelationshipsTag, 1)

	return rID, nil
}

func (d *ReplaceDocx) newRelsID() (*string, error) {
	hlinks, err := d.GetRels()
	if err != nil {
		return nil, err
	}

	max := int64(0)
	for _, hlink := range hlinks {
		re := regexp.MustCompile("[0-9]+")
		finds := re.FindAllString(hlink.ID, -1)
		if len(finds) != 1 {
			return nil, err
		}

		if id, err := strconv.ParseInt(finds[0], 10, 64); err != nil {
			return nil, err
		} else if id > max {
			max = id
		}
	}

	newID := fmt.Sprintf("rId%v", max+1)
	return &newID, nil
}

// GetHyperlinkRel if exists returns the id of the hyperlink, otherwise nil
func (d *ReplaceDocx) GetHyperlinkRel(target string) (*string, error) {
	hRels, err := d.GetHyperlinksRels()
	if err != nil {
		return nil, err
	}

	// check if exists
	for _, hRel := range hRels {
		if hRel.Type == hyperlinkType && hRel.Target == target {
			return &hRel.ID, nil
		}
	}

	return nil, nil
}
