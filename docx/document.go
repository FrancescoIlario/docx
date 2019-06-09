package docx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ExtractParagraphs extract paragraphs
func (d *Docx) ExtractParagraphs() ([]*Paragraph, error) {
	paragraphsContent, err := d.extractParagraphsText()
	if err != nil {
		return nil, err
	}

	var paragraphs []*Paragraph
	for _, paragraphContent := range paragraphsContent {
		paragraph, err := ParseParagraph(paragraphContent)
		if err != nil {
			return nil, err
		}
		paragraphs = append(paragraphs, paragraph)
	}

	return paragraphs, nil
}

func (d *Docx) extractParagraphsText() ([]string, error) {
	reader := strings.NewReader(d.content)
	decoder := xml.NewDecoder(reader)

	var paragraphs []string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

		switch Element := token.(type) {
		case xml.StartElement:
			if Element.Name.Local == paragraphLocalname {
				encoded, err := encodeToken(decoder, token, paragraphLocalname)
				if err != nil {
					return nil, err
				}

				paragraphs = append(paragraphs, *encoded)
			}
		}
	}

	return paragraphs, nil
}

func encodeToken(decoder *xml.Decoder, token xml.Token, wrappingTag string) (*string, error) {
	buffer := bytes.NewBufferString("")
	encoder := xml.NewEncoder(buffer)

	err := encoder.EncodeToken(token)
	if err != nil {
		return nil, err
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil, err
		}

		if err = encoder.EncodeToken(token); err != nil {
			return nil, err
		}

		element, ok := token.(xml.EndElement)
		if ok && element.Name.Local == wrappingTag {
			break
		}
	}

	encoder.Flush()

	tokenString := fmt.Sprintln(buffer)
	return &tokenString, nil
}
