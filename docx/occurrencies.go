package docx

import (
	"encoding/xml"
	"io"
	"strings"
)

// CountOccurrences retrieve all the occurrences indices of a string in the document
func (d *ReplaceDocx) CountOccurrences(lookFor string, matchCase bool, wrapWord bool) int {
	var occurrences int
	if wrapWord {
		occurrences = d.countOccurrences(lookFor, matchCase)
	} else {
		occurrences = d.countSeparatedOccurrences(lookFor, matchCase)
	}
	return occurrences
}

func (d *ReplaceDocx) countOccurrences(lookFor string, matchCase bool) int {
	text := d.GetText()
	_lookFor := lookFor

	if !matchCase {
		text = strings.ToLower(text)
		_lookFor = strings.ToLower(lookFor)
	}

	var result []int

	bi, li, idx := 0, 0, 0
	vtext := text
	for {
		vtext = vtext[li:]
		if idx = strings.Index(vtext, _lookFor); idx == -1 {
			break
		}

		result = append(result, idx+bi)
		bi = li
		li = idx + len(lookFor)
	}

	return len(result)
}

func (d *ReplaceDocx) countSeparatedOccurrences(lookFor string, matchCase bool) int {
	text := d.GetText()
	if !matchCase {
		text = strings.ToLower(text)
		lookFor = strings.ToLower(lookFor)
	}

	fields := strings.Fields(text)
	counter := 0
	for _, field := range fields {
		if strings.Trim(field, `,<.>/?\|"';:?~!@#$%^&*()[]{}-=_+`) == lookFor {
			counter++
		}
	}

	return counter
}

// GetText get the text in the document
func (d *ReplaceDocx) GetText() string {
	var texts []string

	strReader := strings.NewReader(d.content)
	decoder := xml.NewDecoder(strReader)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

		switch Element := token.(type) {
		case xml.StartElement:
			if Element.Name.Local == "t" {

				if t, err := decoder.Token(); t != nil {
					if err != nil {
						break
					}

					switch REl := t.(type) {
					case xml.CharData:
						cd := string(REl)
						texts = append(texts, cd)
					}
				}
			} else if Element.Name.Local == "p" {
				texts = append(texts, "\n")
			}
			break
		}
	}

	return strings.Join(texts[1:], "")
}
