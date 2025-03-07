package gen

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v28/ignite/config"
	chainconfig "github.com/ignite/cli/v28/ignite/config/chain"
	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/ignite/cli/v28/ignite/version"
	/*chainConfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/plugin"*/)

const cacheFileName = "ignite_cache.db"

type generator struct {
	//modulePath gomodulepath.Path
	config    *chainconfig.Config
	storage   cache.Storage
	appPath   string
	protoPath string
	/*outPath      string
	csModulePath string*/
}

func New(ctx context.Context, api plugin.ClientAPI) (*generator, error) {
	/*flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	c, err := getChain(cmd)
	if err != nil {
		return nil, err
	}
	config, err := c.Config()
	if err != nil {
		return nil, err
	}

	p, appPath, err := getPath(cmd)
	if err != nil {
		return nil, err
	}

	outFlag, _ := flags.GetString("out")
	if outFlag == "" {
		outFlag = "./cs"
		}*/

	chainInfo, err := api.GetChainInfo(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get chain info: %s", err)
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

	storage, err := newCache()
	if err != nil {
		return nil, err
	}

	fmt.Println(ch.Name())

	err = InstallDepTools(ctx)
	if err != nil {
		return nil, err
	}

	gen := generator{
		config:    config,
		storage:   storage,
		appPath:   ch.AppPath(),
		protoPath: filepath.Join(ch.AppPath(), config.Build.Proto.Path),
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

func getPath(cmd *plugin.ExecutedCommand) (gomodulepath.Path, string, error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return gomodulepath.Path{}, "", err
	}
	path, _ := flags.GetString("path")
	absPath, err := filepath.Abs(path)
	if err != nil {
		return gomodulepath.Path{}, "", err
	}

	return gomodulepath.Find(absPath)
}

func getModulePath(rawPath string) string {
	return strings.Join(strings.Split(rawPath, "/")[1:], ".")
}

func newCache() (cache.Storage, error) {
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

	return storage, nil
}
