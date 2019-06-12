package main

import (
	"log"

	"github.com/FrancescoIlario/docx/docx"
)

const (
	filepathTestDocument     = `data/docx/TestDocument.docx`
	filepathSimplestDocument = `data/docx/SimplestDocument.docx`
	filepathOutputDocument   = `data/docx/HyperlinkedDocument.docx`
)

func main() {
	err := docx.ReplaceTextWithHyperlink(filepathTestDocument, filepathOutputDocument, "This", "Allegato.txt")
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
