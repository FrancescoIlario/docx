package main

//go:generate go run assets/fgenerate.go

import (
	"log"

	"github.com/FrancescoIlario/docx"
)

const (
	filepathTestDocument     = `data/docx/TestDocument.docx`
	filepathSimplestDocument = `data/docx/SimplestDocument.docx`
	filepathOutputDocument   = `data/docx/HyperlinkedDocument.docx`
)

func main() {
	doc, err := docx.ReadDocxFile(filepathSimplestDocument)
	panicIf(err)
	defer doc.Close()

	err = doc.ReplaceTextWithHyperlink("This", "Allegato.txt")
	panicIf(err)

	err = doc.Editable().WriteToFile(filepathOutputDocument)
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
