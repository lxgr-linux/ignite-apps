package gen

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/DecentralCardGame/protoc-gen-cosmos-csharp/descriptor"
	"github.com/ignite/apps/cs-client/customconfig"
	"github.com/ignite/cli/v29/ignite/config"
	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/dircache"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	"github.com/ignite/cli/v29/ignite/version"
)

const (
	cacheFileName  = "ignite_cache.db"
	flagClearCache = "clear-cache"
)

type generator struct {
	modulePath  gomodulepath.Path
	config      *chainconfig.Config
	storage     cache.Storage
	appPath     string
	protoPath   string
	outPath     string
	csNameSpace descriptor.Descriptor
	/*csModulePath string*/
}

func New(ctx context.Context, cmd *plugin.ExecutedCommand, chainInfo *plugin.ChainInfo) (*generator, error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	ch, err := chain.New(chainInfo.AppPath)
	if err != nil {
		return nil, errors.Errorf("failed to create a new chain object from app path: %s", err)
	}

	config, err := ch.Config()
	if err != nil {
		return nil, errors.Errorf("failed to get config: %s", err)
	}

	//csModulePath := getModulePath(p.RawPath)

	storage, err := newCache(cmd)
	if err != nil {
		return nil, err
	}

	yamlConfig, err := customconfig.Read(ch.ConfigPath())
	if err != nil {
		return nil, err
	}

	out, _ := flags.GetString("out")
	if out == "" {
		if yamlConfig.Client != nil && yamlConfig.Client.CsClient != nil {
			out = yamlConfig.Client.CsClient.Path
		} else {
			out = "./cs"
		}
	}

	modPath, _, err := gomodulepath.Find(ch.AppPath())
	if err != nil {
		return nil, err
	}

	gen := generator{
		modulePath:  modPath,
		outPath:     out,
		config:      config,
		storage:     storage,
		appPath:     ch.AppPath(),
		protoPath:   filepath.Join(ch.AppPath(), config.Build.Proto.Path),
		csNameSpace: descriptor.FromTypeUrl(modPath.Package),
	}

	return &gen, nil
}

func getChain(cmd *plugin.ExecutedCommand, chainOption ...chain.Option) (*chain.Chain, error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	var (
		home, _ = flags.GetString("home")
		path, _ = flags.GetString("path")
	)
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return chain.New(absPath, chainOption...)
}

func getModulePath(rawPath string) string {
	return strings.Join(strings.Split(rawPath, "/")[1:], ".")
}

func flagGetClearCache(cmd *plugin.ExecutedCommand) bool {
	flags, err := cmd.NewFlags()
	if err != nil {
		return false
	}

	clearCache, _ := flags.GetBool(flagClearCache)
	return clearCache
}

func newCache(cmd *plugin.ExecutedCommand) (cache.Storage, error) {
	cacheRootDir, err := config.DirPath()
	if err != nil {
		return cache.Storage{}, err
	}

	storage, err := cache.NewStorage(
		filepath.Join(cacheRootDir, cacheFileName),
		cache.WithVersion(version.Version),
	)
	if err != nil {
		return cache.Storage{}, err
	}

	if flagGetClearCache(cmd) {
		if err := storage.Clear(); err != nil {
			return cache.Storage{}, err
		}
		if err := dircache.ClearCache(); err != nil {
			return cache.Storage{}, err
		}
	}

	return storage, nil
}
