package docx

import (
	"regexp"
	"testing"
)

func TestSubstituteWithLinks(t *testing.T) {
	replaceDocx, err := ReadDocxFile("../data/docx/TestDocument.docx")
	if err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
	defer replaceDocx.Close()

	var subs []struct{ target, link string }
	subs = append(subs, struct{ target, link string }{
		target: "This",
		link:   "Allegato.txt",
	})

	if err = replaceDocx.SubstituteWithLinks(subs); err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
	var found *string
	hyperlinks := replaceDocx.GetHyperlinks()
	for _, hyperlink := range hyperlinks {
		finds := regexp.MustCompile(textRegexWithSubmatches).FindStringSubmatch(hyperlink)
		findsLen := len(finds)
		if findsLen != 3 {
			t.Errorf("More than one text in hyperlink %s\n", hyperlink)
		}

		foundEl := finds[2]
		if foundEl == subs[0].target {
			found = &foundEl
			break
		}
	}

	if err := replaceDocx.Editable().WriteToFile("../data/docx/produced.docx"); err != nil {
		t.Errorf("Can not save the edited file\n")
	}

	if found == nil {
		t.Errorf("No hyperlink has been generated for the target: %s\n", subs[0].target)
	}
}
