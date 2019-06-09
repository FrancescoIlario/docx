package docx

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestParseRunSimpleRun(t *testing.T) {
	xmlSlice := `<w:r><w:rPr></w:rPr><w:t>.</w:t></w:r>`

	if _, err := ParseRun(xmlSlice, false); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
}

func TestParseRunEmptyComplexRun(t *testing.T) {
	xmlSlice := `<w:r></w:r>`

	if _, err := ParseRun(xmlSlice, false); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
}

func TestParseRunEmptySimpleRun(t *testing.T) {
	xmlSlice := `<w:r/>`

	if _, err := ParseRun(xmlSlice, false); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
}

func TestExtractRunsFromParagraphComplexParagraph(t *testing.T) {
	xmlSlice := `<w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:rPr></w:rPr><w:t xml:space="preserve">This is a </w:t></w:r><w:r><w:rPr></w:rPr><w:t>.</w:t></w:r><w:r></w:r></w:p>`
	buffer := bytes.NewBufferString("")

	if err := xml.EscapeText(buffer, []byte(xmlSlice)); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
	if _, err := ExtractRunsFromParagraph(xmlSlice); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
}

func TestExtractRunsFromParagraphSimpleParagraph(t *testing.T) {
	xmlSlice := `<w:p><w:r></w:r></w:p>`
	buffer := bytes.NewBufferString("")

	if err := xml.EscapeText(buffer, []byte(xmlSlice)); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}

	if _, err := ExtractRunsFromParagraph(xmlSlice); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
}
