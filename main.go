package main

import (
	"flag"
	"fmt"
	"log"

	fdocx "github.com/FrancescoIlario/docx/docx"
)

var (
	lookFor = flag.String("LookFor", "This", "The substring to look")
)

func main() {
	var editable *fdocx.ReplaceDocx
	var err error

	editable, err = fdocx.ReadDocxFile("data/docx/TestDocument.docx")
	panicIf(err)

	text := editable.GetText()
	fmt.Println(text)

	occs := editable.GetOccurrences(*lookFor, true)
	for _, occ := range occs {
		fmt.Println(occ)
	}
}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
