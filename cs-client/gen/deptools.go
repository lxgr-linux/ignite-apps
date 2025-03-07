package gen

import (
	"context"

	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
)

func DepTools() []string {
	return []string{
		// buf build code generation.
		"github.com/bufbuild/buf/cmd/buf@v1.50.0",
		"github.com/DecentralCardGame/protoc-gen-cosmos-csharp@v0.1.2",
	}
}

func InstallDepTools(ctx context.Context) error {
	for _, tool := range DepTools() {
		err := gocmd.Install(ctx, ".", []string{tool})
		if err != nil {
			return err
		}
	}

	return nil
}
