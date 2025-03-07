package gen

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosbuf"
)

//go:embed templates/buf.gen.yaml.tpl
var bufGenYamlTpl string

type bufGenYamlModel struct {
	PackagePrefix string
	OutDir        string
	InDir         string
}

func (g generator) createBufGenYaml() (string, error) {
	tmpl, err := template.New("buf.gen.yaml").Parse(bufGenYamlTpl)
	if err != nil {
		return "", err
	}

	model := bufGenYamlModel{
		OutDir:        ".",
		InDir:         g.protoPath,
		PackagePrefix: "yes",
	}

	wr := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(wr, model)
	if err != nil {
		return "", err
	}

	return wr.String(), nil
}

func (g generator) GenerateClients(ctx context.Context) error {
	fmt.Println("Generating clients...")

	buf, err := cosmosbuf.New(g.storage, g.appPath)
	if err != nil {
		return err
	}

	bufGenYaml, err := g.createBufGenYaml()
	if err != nil {
		return err
	}

	return buf.Generate(ctx, g.protoPath, g.outPath, bufGenYaml)
}
