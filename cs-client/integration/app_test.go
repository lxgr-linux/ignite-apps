package integration_test

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ignite/apps/cs-client/gen"
	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	envtest "github.com/ignite/cli/v29/integration"
	"github.com/stretchr/testify/require"
)

//go:embed static/*
var staticFiles embed.FS

func TestCsClient(t *testing.T) {
	var (
		require     = require.New(t)
		env         = envtest.New(t)
		app         = env.Scaffold("github.com/a/niceChain")
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "cs-client")

	homeDir := app.UseRandomHomeDir()

	require.NoError(gen.InstallStaticFiles(staticFiles, app.SourcePath(), "static"))

	env.Must(env.Exec("install cs-client app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	var outPath = "outpath"

	env.Must(env.Exec("scaffold testitem",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "scaffold", "list", "testItem", "-y"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("client gen",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "generate", "cs-client", "-y", "-o="+outPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("build",
		step.NewSteps(step.New(
			step.Exec("dotnet", "build"),
			step.Workdir(filepath.Join(app.SourcePath(), outPath)),
		)),
	))

	// Getting testing data
	o := &bytes.Buffer{}
	env.Must(env.Exec("init chain",
		step.NewSteps(step.New(
			step.Stdout(o),
			step.Exec(envtest.IgniteApp, "chain", "init", "-y", "--home", homeDir),
			step.Workdir(app.SourcePath()),
		)),
	))

	output := &bytes.Buffer{}

	env.Must(env.Exec("export keys",
		step.NewSteps(step.New(
			step.Stdout(output),
			step.Exec(app.Binary(), "keys", "export", "alice", "--unsafe", "--unarmored-hex", "--home", homeDir),
			step.Workdir(app.SourcePath()),
			step.Stdin(strings.NewReader("y\n")),
		)),
	))

	fmt.Println("yes")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		env.Must(env.Exec("should serve",
			step.NewSteps(step.New(
				step.Stdout(o),
				step.Exec(app.Binary(), "start", "--grpc.address", servers.GRPC, "--home", homeDir),
				step.Workdir(app.SourcePath()),
			)),
			envtest.ExecCtx(ctx),
		))
		wg.Done()
	}()

	time.Sleep(5)
	fmt.Println(o.String() + "\n----")

	env.Exec("run testapp",
		step.NewSteps(step.New(
			/*step.PreExec(func() error {
				return env.IsAppServed(ctx, servers.API)
			}),*/
			step.Exec("dotnet", "run", strings.TrimSuffix(output.String(), "\n"), "http://"+servers.GRPC),
			step.Workdir(filepath.Join(app.SourcePath(), "testApp")),
		)),
		//envtest.ExecRetry(),
	)

	cancel()

	//fmt.Println(o.String())

	wg.Wait()
}

func assertLocalPlugins(t *testing.T, app envtest.App, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfg, err := pluginsconfig.ParseDir(app.SourcePath())
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected local apps")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
