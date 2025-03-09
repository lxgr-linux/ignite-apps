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
	NameSpace string
	Services  []ServiceModel
}

type ServiceModel struct {
	Path descriptor.Desriptor
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
		NameSpace: strings.Title(g.modulePath.Package),
	}

	txClientModel := ClientModel{
		NameSpace: strings.Title(g.modulePath.Package),
	}

	for _, pkg := range pkgs {
		for _, service := range pkg.Services {
			s := ServiceModel{
				Type: service.Name,
				Path: descriptor.FromTypeUrl(pkg.Name).CutNameSpace(descriptor.FromTypeUrl(g.modulePath.Package)),
				Name: getSimpleModuleNameFromPath(pkg.Name),
			}
			switch service.Name {
			case "Query":
				queryClientModel.Services = append(queryClientModel.Services, s)
			case "Msg":
				txClientModel.Services = append(txClientModel.Services, s)
			}
		}
	}

	tmpl, err := template.New("txClient").Parse(txClientTmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(g.outPath, txClientModel.NameSpace, "TxClient.cs")
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

	path = filepath.Join(g.outPath, txClientModel.NameSpace, "QueryClient.cs")
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, queryClientModel)
}

func getSimpleModuleNameFromPath(path string) string {
	s := strings.Split(path, ".")

	if strings.Contains(s[len(s)-1], "v1") {
		return strings.Title(s[len(s)-2]) + strings.Title(s[len(s)-1])
	} else {
		return strings.Title(s[len(s)-1])
	}
}
