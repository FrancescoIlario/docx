package docx

import (
	"log"
	"regexp"
)

const (
	hyperlinkRegex = `<w:hyperlink r:id=".*?".*?>.*?</w:hyperlink>`
)

// GetHyperlinks returns the xml source code for all the hyperlinks
func (d *ReplaceDocx) GetHyperlinks() []string {
	regex := regexp.MustCompile(hyperlinkRegex)
	matches := regex.FindAllString(d.content, -1)

	log.Println(matches)
	return matches
}
