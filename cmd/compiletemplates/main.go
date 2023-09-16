package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

var removeHTML = strings.NewReplacer(".html", "")

type config struct {
	TemplateName string
	Package      string
	Var          string
	FileName     string
}

func main() {
	outputDir, err := os.Getwd()
	lerr.Panic(err)

	templatesDir := outputDir
	args := os.Args[1:]
	if len(args) > 0 {
		templatesDir = filepath.Join(templatesDir, args[0])
	}
	os.Chdir(templatesDir)

	cfgFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Expect file 'config.json' with fields 'TemplateName', 'Package', 'Var', and 'FileName'")
	}
	cfg := &config{}
	lerr.Panic(json.NewDecoder(cfgFile).Decode(cfg))
	cfgFile.Close()

	files, err := filepath.Glob("*.html")
	lerr.Panic(err)

	ts := make([]string, len(files))
	for i, name := range files {
		file, err := ioutil.ReadFile(name)
		lerr.Panic(err)
		name = removeHTML.Replace(name)
		ts[i] = fmt.Sprintf("{{define \"%s\" -}}\n%s\n{{- end}}", name, file)
	}

	lerr.Panic(os.Chdir(outputDir))
	out, err := os.Create(cfg.FileName)
	lerr.Panic(err)

	fmt.Fprintf(out, "package %s\n\nimport \"html/template\"\n\nvar %s = template.Must(template.New(\"%s\").Parse(`\n%s\n`))\n", cfg.Package, cfg.Var, cfg.TemplateName, strings.Join(ts, "\n\n"))
	out.Close()

}
