package docx

import (
	"testing"
)

func TestCheckParagraphSliceGoodExtended(t *testing.T) {
	xmlSlice := `<p></p>`
	if err := checkParagraphSlice(xmlSlice); err != nil {
		t.Errorf("Unexpected %s", err)
	}
}

func TestCheckParagraphSliceGoodCompact(t *testing.T) {
	xmlSlice := `<p/>`
	if err := checkParagraphSlice(xmlSlice); err != nil {
		t.Errorf("Unexpected %s", err)
	}
}

func TestCheckParagraphSliceGoodComplete(t *testing.T) {
	xmlSlice := `<p><r></r></p>`
	if err := checkParagraphSlice(xmlSlice); err != nil {
		t.Errorf("Unexpected %s", err)
	}
}

func TestCheckParagraphSliceWrongOnlyOpening(t *testing.T) {
	xmlSlice := `<p>`
	if err := checkParagraphSlice(xmlSlice); err == nil {
		t.Errorf("expected error, returned nil")
	}
}

func TestCheckParagraphSliceWrongOnlyClosing(t *testing.T) {
	xmlSlice := `</p>`
	if err := checkParagraphSlice(xmlSlice); err == nil {
		t.Errorf("expected error, returned nil")
	}
}

func TestCheckParagraphSliceWrongTag(t *testing.T) {
	xmlSlice := `<>`
	if err := checkParagraphSlice(xmlSlice); err == nil {
		t.Errorf("expected error, returned nil")
	}
}

func TestCheckParagraphSliceWrongExtendedNotComplete(t *testing.T) {
	xmlSlice := `<p><r></r>`
	if err := checkParagraphSlice(xmlSlice); err == nil {
		t.Errorf("expected error, returned nil")
	}
}

func TestCheckParagraphSliceWrongExtendedDifferent(t *testing.T) {
	xmlSlice := `<r></r>`
	if err := checkParagraphSlice(xmlSlice); err == nil {
		t.Errorf("expected error, returned nil")
	}
}
