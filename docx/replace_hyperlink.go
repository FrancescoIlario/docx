package docx

import (
	"html"
	"log"
	"regexp"
	"strings"

	"github.com/go-xmlfmt/xmlfmt"

	"github.com/FrancescoIlario/docx/stringext"
	"github.com/cbroglie/mustache"
)

const (
	runRegex                = `<w:r(?: .*?)?>.*?<w:t(?: .*?)?>.*?{{Text}}.*?</w:t></w:r>`
	runRegexWithSubmatches  = `<w:r( .*?)?>(.*?)<w:t( .*?)?>.*?{{Text}}.*?</w:t></w:r>`
	textRegex               = `<w:t(?: .*?)?>.*?</w:t>`
	textRegexWithSubmatches = `<w:t( .*?)?>(.*?)</w:t>`
	runTemplate             = `<w:r{{rAttrs}}>{{rPr}}<w:t>{{Text}}</w:t></w:r>`
	hyperlinkTemplate       = `<w:hyperlink r:id="{{ID}}"><w:r><w:rPr><w:rStyle w:val="InternetLink"/></w:rPr><w:t>{{Text}}</w:t></w:r></w:hyperlink>`
)

// hyperlinkReplaceData data for substitution
type hyperlinkReplaceData struct {
	ID   string
	Text string
	Link string
}

// SubstituteWithLinks SubstituteWithLinks
func (d *ReplaceDocx) SubstituteWithLinks(subs []struct{ target, link string }) error {
	for _, sub := range subs {
		rID, err := d.AddHyperlinkRel(sub.link)
		if err != nil {
			return err
		}

		substitutionData := hyperlinkReplaceData{
			ID:   *rID,
			Text: sub.target,
		}

		if err := d.applySubstitution(substitutionData); err != nil {
			return err
		}
	}

	return nil
}

// SubstituteWithLink SubstituteWithLink
func (d *ReplaceDocx) SubstituteWithLink(target, link string) error {
	rID, err := d.AddHyperlinkRel(target)
	if err != nil {
		return err
	}

	substitutionData := hyperlinkReplaceData{
		ID:   *rID,
		Text: target,
	}

	if err := d.applySubstitution(substitutionData); err != nil {
		return err
	}

	return nil
}

func (d *ReplaceDocx) applySubstitutions(substitutionData []hyperlinkReplaceData) error {
	for _, substitution := range substitutionData {
		if err := d.applySubstitution(substitution); err != nil {
			return err
		}
	}

	return nil
}

func (d *ReplaceDocx) applySubstitution(substitution hyperlinkReplaceData) error {
	subRunRegex, err := mustache.Render(runRegex, substitution)
	if err != nil {
		return err
	}

	regex := regexp.MustCompile(subRunRegex)
	foundAll := regex.FindAllString(d.content, -1)
	log.Println(foundAll)

	d.content = regex.ReplaceAllStringFunc(d.content, func(match string) string {
		log.Printf("matched: %s\n", match)
		if hrun, err := replaceRunWithHyperlink(match, &substitution); err != nil {
			log.Printf("Error obtaining replace for run %s: %s\n", match, err)
		} else {
			log.Printf("Obtained replace for run %s: %s\n", match, *hrun)
			return *hrun
		}
		return match
	})

	log.Printf("New content:\n%s\n", xmlfmt.FormatXML(d.content, "", "  "))
	return nil
}

func replaceRunWithHyperlink(match string, substitutionData *hyperlinkReplaceData) (*string, error) {
	var runs []string

	subRunRegex, err := mustache.Render(runRegexWithSubmatches, substitutionData)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(subRunRegex)

	submatchesRun := regex.FindStringSubmatch(match)
	rAttrs, rPr := submatchesRun[1], submatchesRun[2]

	submatchesText := regexp.MustCompile(textRegexWithSubmatches).FindStringSubmatch(match)
	text := submatchesText[1]

	splits := stringext.SplitAfterWithSeparator(text, substitutionData.Text)
	for _, split := range splits {
		if split == substitutionData.Text {
			substitutionMapping := map[string]string{
				"ID":   substitutionData.ID,
				"Text": split,
			}

			hyperlinkEntry, err := mustache.Render(hyperlinkTemplate, substitutionMapping)
			if err != nil {
				return nil, err
			}

			runs = append(runs, html.UnescapeString(hyperlinkEntry))
		} else {
			substitutionMapping := map[string]string{
				"rAttrs": rAttrs,
				"rPr":    rPr,
				"Text":   split,
			}

			runEntry, err := mustache.Render(runTemplate, substitutionMapping)
			if err != nil {
				return nil, err
			}

			runs = append(runs, html.UnescapeString(runEntry))
		}
	}

	runsString := strings.Join(runs, "")
	return &runsString, nil
}
