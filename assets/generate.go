// +build ignore

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var cwd, _ = os.Getwd()
	templates := http.Dir(filepath.Join(cwd, "assets/templates/"))
	if err := vfsgen.Generate(templates, vfsgen.Options{
		Filename:     "assets/templates/templates_vfsdata.go",
		PackageName:  "assets",
		BuildTags:    "",
		VariableName: "Assets",
	}); err != nil {
		log.Fatalln(err)
	}
}
