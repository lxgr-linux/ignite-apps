package gen

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/DecentralCardGame/protoc-gen-cosmos-csharp/descriptor"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis"
)

//go:embed templates/TxClient.cs.tpl
var txClientTmpl string

//go:embed templates/QueryClient.cs.tpl
var queryClientTmpl string

type ClientModel struct {
	NameSpace descriptor.Descriptor
	Services  []ServiceModel
}

type ServiceModel struct {
	Path descriptor.Descriptor
	Name string
	Type string
}

func (g generator) GenerateClient(ctx context.Context) error {
	cache := protoanalysis.NewCache()
	pkgs, err := protoanalysis.Parse(ctx, cache, g.protoPath)
	if err != nil {
		return err
	}

	queryClientModel := ClientModel{
		NameSpace: g.csNameSpace,
	}

	txClientModel := ClientModel{
		NameSpace: g.csNameSpace,
	}

	for _, pkg := range pkgs {
		for _, service := range pkg.Services {
			name := descriptor.FromTypeUrl(pkg.Name)
			s := ServiceModel{
				Type: service.Name,
				Path: name.CutNameSpace(g.csNameSpace),
				Name: getSimpleName(name),
			}
			switch service.Name {
			case "Query":
				queryClientModel.Services = append(queryClientModel.Services, s)
			case "Msg":
				txClientModel.Services = append(txClientModel.Services, s)
			}
		}
	}

	baseOutPath := filepath.Join(g.outPath, g.csNameSpace.Name())

	tmpl, err := template.New("txClient").Parse(txClientTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(baseOutPath, "TxClient.cs")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, txClientModel)
	if err != nil {
		return err
	}

	tmpl, err = template.New("queryClient").Parse(queryClientTmpl)
	if err != nil {
		return err
	}

	path = filepath.Join(baseOutPath, "QueryClient.cs")
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, queryClientModel)
}

func getSimpleName(d descriptor.Descriptor) string {
	name := d.Name()
	if strings.Contains(name, "V1") {
		return d.Parent().Name() + name
	}
	return name
}
