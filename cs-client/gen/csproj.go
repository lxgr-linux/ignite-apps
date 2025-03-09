package gen

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/csproj.xml.tpl
var csprojTmpl string

type csprojModel struct {
	Name      string
	ShortName string
	URL       string
}

func (g generator) GenerateCsproj() error {
	m := csprojModel{
		URL:       "https://" + g.modulePath.RawPath,
		Name:      strings.Title(g.modulePath.Package),
		ShortName: g.modulePath.Package,
	}

	tmpl, err := template.New("csproj").Parse(csprojTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(g.outPath, strings.Title(g.modulePath.Package)+".csproj")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, m)
}
