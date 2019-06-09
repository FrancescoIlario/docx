package docx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	runLocalname       = "r"
	hyperlinkLocalname = "hyperlink"
)

// Run w:r
type Run struct {
	XML         string `json:"xml"`
	Text        string `xml:"t" json:"text"`
	InHyperlink bool   `json:"hyperlink"`

	InnerXML string   `xml:",innerxml"`
	XMLName  xml.Name `xml:"r"`
}

func (r *Run) String() string {
	jsonData, err := json.Marshal(r)
	if err != nil {
		custom := fmt.Sprintf("[XML: %s]\n", r.XML)
		errMessage := fmt.Sprintf("cannot generate json for paragraph %s\n", custom)
		log.Println(errMessage)
		return custom
	}
	return fmt.Sprintln(string(jsonData))
}

// ParseRun initializes a new run
func ParseRun(xmlSlice string, inHyperlink bool) (*Run, error) {
	var run Run
	if err := xml.Unmarshal([]byte(xmlSlice), &run); err != nil {
		return nil, err
	}

	run.XML = xmlSlice
	run.InHyperlink = inHyperlink

	return &run, nil
}

// ExtractRunsFromParagraph extracts the list of runs contained in a paragraph
func ExtractRunsFromParagraph(xmlSlice string) ([]*Run, error) {
	strReader := strings.NewReader(xmlSlice)
	decoder := xml.NewDecoder(strReader)

	var runs []*Run
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

		// check if is in hyperlink
		isHyperlink := false
		if element, ok := token.(xml.StartElement); ok && element.Name.Local == hyperlinkLocalname {
			isHyperlink = true
			token, err = decoder.Token()

			if err == io.EOF {
				break
			}
		}

		// parse run
		element, ok := token.(xml.StartElement)
		if ok && element.Name.Local == runLocalname {
			runString, err := encodeToken(decoder, token, runLocalname)
			if err != nil {
				return nil, err
			}
			log.Printf("extracted run %s\n", *runString)

			run, err := ParseRun(*runString, isHyperlink)
			if err != nil {
				return nil, err
			}

			runs = append(runs, run)
		}
	}

	return runs, nil
}
