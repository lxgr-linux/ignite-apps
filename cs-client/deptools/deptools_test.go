package deptools_test

import (
	"context"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/ignite/apps/cs-client/deptools"
	"github.com/stretchr/testify/assert"
)

func setupDir(path string) {
	os.WriteFile(filepath.Join(path, "go.mod"), []byte(`
module github.com/test/module

go 1.23.6
`), os.ModePerm)
	os.MkdirAll(filepath.Join(path, "tools"), os.ModePerm)
}

func getGoPath() (gopath string) {
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return
}

func TestAddsFileAndDepsWithMissingFile(t *testing.T) {
	gopath := getGoPath()

	dir, _ := os.MkdirTemp("", "*")
	defer os.RemoveAll(dir)
	setupDir(dir)

	err := deptools.ProvideTools(context.Background(), dir)

	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "tools", "cstools.go"))
	assert.NoError(t, err)
	for _, dep := range deptools.DepTools() {
		_, err := os.Stat(filepath.Join(gopath, "bin", filepath.Base(dep)))
		assert.NoError(t, err)
	}
}

func TestAddsFileAndDepsWithMissPopulatedFile(t *testing.T) {
	gopath := getGoPath()

	dir, _ := os.MkdirTemp("", "*")
	defer os.RemoveAll(dir)
	setupDir(dir)
	os.WriteFile(filepath.Join(dir, "tools", "cstools.go"), []byte("THIS IS INVALID"), os.ModePerm)

	err := deptools.ProvideTools(context.Background(), dir)

	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "tools", "cstools.go"))
	assert.NoError(t, err)
	for _, dep := range deptools.DepTools() {
		_, err := os.Stat(filepath.Join(gopath, "bin", filepath.Base(dep)))
		assert.NoError(t, err)
	}
}

func TestDoNothingWhenPopulatedCorrectly(t *testing.T) {
	gopath := getGoPath()

	dir, _ := os.MkdirTemp("", "*")
	defer os.RemoveAll(dir)
	setupDir(dir)
	err := deptools.ProvideTools(context.Background(), dir)
	assert.NoError(t, err)
	toolsStat1, err := os.Stat(filepath.Join(dir, "tools", "cstools.go"))
	var stats []os.FileInfo
	for _, dep := range deptools.DepTools() {
		s, err := os.Stat(filepath.Join(gopath, "bin", filepath.Base(dep)))
		assert.NoError(t, err)
		stats = append(stats, s)
	}

	err = deptools.ProvideTools(context.Background(), dir)

	assert.NoError(t, err)
	toolsStat2, err := os.Stat(filepath.Join(dir, "tools", "cstools.go"))
	assert.True(t, toolsStat1.ModTime().Equal(toolsStat2.ModTime()))
	for idx, dep := range deptools.DepTools() {
		s, err := os.Stat(filepath.Join(gopath, "bin", filepath.Base(dep)))
		assert.NoError(t, err)
		assert.True(t, s.ModTime().Equal(stats[idx].ModTime()))
	}
}
