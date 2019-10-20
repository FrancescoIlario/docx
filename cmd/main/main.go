package main

//go:generate go run assets/fgenerate.go

import (
	"log"

	"github.com/FrancescoIlario/docx"
)

const (
	_filepathTestDocument     = `data/docx/TestDocument.docx`
	_filepathSimplestDocument = `data/docx/SimplestDocument.docx`
	_filepathOutputDocument   = `data/docx/HyperlinkedDocument.docx`
)

func main() {
	panicIf(run())
}

func run() error {
	doc, err := docx.ReadDocxFile(_filepathSimplestDocument)
	if err != nil {
		return err
	}
	defer doc.Close()

	err = doc.ReplaceTextWithHyperlink("This", "Allegato.txt")
	if err != nil {
		return err
	}

	err = doc.Editable().WriteToFile(_filepathOutputDocument)
	if err != nil {
		return err
	}
}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
