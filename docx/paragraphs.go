package docx

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

const (
	paragraphLocalname = "p"
	runLocalname       = "r"
)

// Paragraph w:p
type Paragraph struct {
	XML  string
	PPtr string
	Text string
	Runs []*Run
}

// Run w:r
type Run struct {
	XML       string
	Text      string
	PPtr      string
	Paragraph *Paragraph
}

// WrongXMLSlice returned if not valid XmlSlice is passed
type WrongXMLSlice struct {
	XMLSlice string
	Reason   string
}

func (e *WrongXMLSlice) Error() string {
	return ""
}

// ParseParagraph initializes a new paragraph
func ParseParagraph(xmlSlice string) (*Paragraph, error) {
	// sanity checks
	if err := checkParagraphSlice(xmlSlice); err != nil {
		return nil, err
	}

	var paragraph *Paragraph
	paragraph.XML = xmlSlice

	strReader := strings.NewReader(xmlSlice)
	decoder := xml.NewDecoder(strReader)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		var runs []*Run

		switch Element := token.(type) {
		case xml.StartElement:
			localname := Element.Name.Local
			if localname == paragraphLocalname {
				log.Println("Paragraph StartElement Found")
			} else if localname == runLocalname {
				run := extractRun(&token)
				runs = append(runs, run)
			}
		case xml.EndElement:
			localname := Element.Name.Local
			if localname == paragraphLocalname {
				log.Println("Paragraph EndElement Found")
			} else if localname == runLocalname {
				reason := "Run: close-tag found without any open-tag"
				log.Println(reason)
				return nil, &WrongXMLSlice{
					XMLSlice: xmlSlice,
					Reason:   reason,
				}
			}
		}
	}

	return paragraph, nil
}

func extractRun(token *xml.Token) *Run {
	log.Println("Run StartElement Found")
	run := &Run{}

	// TODO: Complete

	log.Println("Run EndElement Found")
	return run
}

func checkParagraphSlice(xmlSlice string) error {
	regexSimple := `^<p/>$`
	regexComplex := `^<p>.*</p>$`
	matchedSimple, err := regexp.Match(regexSimple, []byte(xmlSlice))
	if err != nil {
		return err
	}

	matchedComplex, err := regexp.Match(regexComplex, []byte(xmlSlice))
	if err != nil {
		return err
	}

	if !matchedComplex && !matchedSimple {
		reason := fmt.Sprintf("the slice is not respecting the regex: %s", regexComplex)

		return &WrongXMLSlice{
			XMLSlice: xmlSlice,
			Reason:   reason,
		}
	}
	return nil
}

// ParseRun initializes a new run
func ParseRun(xmlSlice string) Run {
	return Run{}
}

// Paragraphs returns all the document's paragraph
func (d *ReplaceDocx) Paragraphs() []Paragraph {
	return make([]Paragraph, 0)
}
