package gen

import (
	_ "embed"
	"os"
	"path/filepath"
	"text/template"
)

const COSMCS_VERSION = "0.7.0-rc2"

//go:embed templates/csproj.xml.tpl
var csprojTmpl string

type csprojModel struct {
	Name          string
	ShortName     string
	URL           string
	CosmcsVersion string
}

func (g generator) GenerateCsproj() error {
	m := csprojModel{
		URL:           "https://" + g.modulePath.RawPath,
		Name:          g.csNameSpace.String(),
		ShortName:     g.modulePath.Package,
		CosmcsVersion: COSMCS_VERSION,
	}

	tmpl, err := template.New("csproj").Parse(csprojTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(g.outPath, g.csNameSpace.String()+".csproj")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, m)
}
