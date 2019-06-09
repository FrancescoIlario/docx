package docx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

const (
	paragraphLocalname = "p"
)

// Paragraph w:p
type Paragraph struct {
	XML           string
	PPtr          string `xml:"pPr" json:"pPr"`
	Runs          []*Run `xml:"r" json:"runs"`
	HyperlinkRuns []*Run `xml:"hyperlink>r" json:"hyperlinkRuns"`

	InnerXML string   `xml:",innerxml"`
	XMLName  xml.Name `xml:"p"`
}

func (p *Paragraph) String() string {
	jsonData, err := json.Marshal(p)
	if err != nil {
		custom := fmt.Sprintf("[XML: %s]\n", p.XML)
		errMessage := fmt.Sprintf("cannot generate json for paragraph %s\n", custom)
		log.Println(errMessage)
		return custom
	}
	return fmt.Sprintln(string(jsonData))
}

// Text The text of all the runs
func (p *Paragraph) Text(includeHyperlinks bool) string {
	var text string

	for _, r := range p.Runs {
		if includeHyperlinks || !r.InHyperlink {
			text += r.Text
		}
	}

	return text
}

// ParseParagraph initializes a new paragraph
func ParseParagraph(xmlSlice string) (*Paragraph, error) {
	runs, err := ExtractRunsFromParagraph(xmlSlice)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	paragraph := &Paragraph{
		XML:  xmlSlice,
		Runs: runs,
	}

	return paragraph, nil
}
