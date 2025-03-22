package deptools

import (
	"bytes"
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
)

//go:embed templates/cstools.go.tpl
var cstoolsTmpl string

const ToolsPath = "tools/cstools.go"

func DepTools() []string {
	return []string{
		// buf build code generation.
		"github.com/bufbuild/buf/cmd/buf",
		"github.com/DecentralCardGame/protoc-gen-cosmos-csharp",
	}
}

func isToolsFilePopulated(path string, newToolsContent []byte) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	oldContent, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return bytes.Equal(oldContent, newToolsContent)
}

func generateNewToolsContent() ([]byte, error) {
	tmpl, err := template.New("cstools").Parse(cstoolsTmpl)
	if err != nil {
		return []byte{}, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, DepTools())

	return buf.Bytes(), err
}

func ProvideTools(ctx context.Context, appPath string) error {
	toolsPath := filepath.Join(appPath, ToolsPath)

	newToolsContent, err := generateNewToolsContent()
	if err != nil {
		return err
	}

	if !isToolsFilePopulated(toolsPath, newToolsContent) {
		err = os.WriteFile(toolsPath, newToolsContent, os.ModePerm)
		if err != nil {
			return err
		}

		InstallDepTools(ctx, appPath)
	}

	return nil
}

func InstallDepTools(ctx context.Context, appPath string) error {
	if err := gocmd.ModTidy(ctx, appPath); err != nil {
		return err
	}

	if err := gocmd.Get(ctx, appPath, []string{"github.com/bufbuild/buf/cmd/buf@v1.50.0"}); err != nil {
		return err
	}

	err := gocmd.Install(ctx, appPath, DepTools())
	if err != nil {
		return err
	}

	return err
}
