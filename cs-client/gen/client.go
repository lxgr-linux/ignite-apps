package gen

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/DecentralCardGame/protoc-gen-cosmos-csharp/descriptor"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis"
)

//go:embed templates/TxClient.cs.tpl
var txClientTmpl string

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
			fmt.Printf("%s, %s, %s\n", pkg.ModuleName(), service.Name, strings.Title(pkg.Name))

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

	path := filepath.Join(g.outPath, txClientModel.NameSpace, "Client.cs")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, txClientModel)
}

func getSimpleModuleNameFromPath(path string) string {
	s := strings.Split(path, ".")

	if strings.Contains(s[len(s)-1], "v1") {
		return strings.Title(s[len(s)-2]) + strings.Title(s[len(s)-1])
	} else {
		return strings.Title(s[len(s)-1])
	}
}
