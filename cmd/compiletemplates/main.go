package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/adamcolton/luce/lerr"
	ljson "github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/ltmpl"
	"github.com/adamcolton/luce/util/luceio"
)

type config struct {
	TemplateName string
	Package      string
	Var          string
	FileName     string
	Path         string
	Comment      string
	SkipImport   bool
}

func main() {
	cfgName := "config.json"
	if args := os.Args; len(args) > 1 {
		cfgName = args[1]
	}

	cfg := &config{}
	err := ljson.Deserializer{}.Load(cfg, cfgName)
	lerr.Panic(err)
	var outdir string
	if cfg.Path != "" {
		outdir, err = os.Getwd()
		lerr.Panic(err)
		lerr.Panic(os.Chdir(cfg.Path))
	}

	m := &lfile.Match{
		SkipDir: filter.MustRegex(lfile.FilterHidden),
	}
	m.Find.File = func(name string) bool {
		if name == cfgName {
			return false
		}
		return !strings.HasSuffix(name, ".go")
	}
	files, _, err := m.Do("./")
	lerr.Panic(err)

	loader := &ltmpl.HTMLLoader{
		Iterator: lfile.Filenames(files),
	}
	buf, sw := luceio.BufferSumWriter()
	if cfg.Package != "" {
		sw.Fprint("package %s\n\n", cfg.Package)
	}
	if !cfg.SkipImport {
		sw.WriteString("import \"html/template\"\n\n")
	}
	if cfg.Comment != "" {
		sw.Fprint("//%s\n\n", cfg.Comment)
	}
	sw.Fprint("var %s = template.Must(template.New(\"%s\").Parse(`\n", cfg.Var, cfg.TemplateName)
	loader.WriteTo(buf)
	sw.WriteString("`))\n")

	if cfg.FileName != "" {
		if cfg.Path != "" {
			os.Chdir(outdir)
		}
		out, err := os.Create(cfg.FileName)
		lerr.Panic(err)
		out.Write(buf.Bytes())
		out.Close()
	} else {
		fmt.Println(buf.String())
	}
}
