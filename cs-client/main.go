package main

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/apps/cs-client/cmd"
	"github.com/ignite/apps/cs-client/gen"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "cs-client",
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	session := cliui.New(cliui.StartSpinnerWithText("Testing spinner..."))
	defer session.End()

	g, err := gen.New(ctx, cmd, api)
	if err != nil {
		return errors.Errorf("failed to init genrator: %s", err)
	}

	_ = g

	//time.Sleep(time.Second * 5)
	session.StopSpinner()

	session.StartSpinner("Installing dependencies...")
	gen.InstallDepTools(ctx)
	session.StopSpinner()

	err = g.GenerateClients(ctx)
	if err != nil {
		return err
	}

	err = g.GenerateCsproj()
	if err != nil {
		return err
	}
	return nil
}

func (app) ExecuteHookPre(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPost(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookCleanUp(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"cs-client": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
