package gen

import (
	"embed"
	_ "embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

func (g *generator) InstallStaticFiles() error {
	return InstallStaticFiles(staticFiles, g.outPath, "static")
}

func InstallStaticFiles(files fs.FS, outPath string, baseDir string) error {
	return fs.WalkDir(
		files, baseDir,
		func(path string, d fs.DirEntry, err error) error {
			name, err := filepath.Rel(
				baseDir, path,
			)
			if err != nil {
				return err
			}
			writePath := filepath.Join(outPath, name)

			if d.Type().IsRegular() {
				content, err := fs.ReadFile(files, path)
				if err != nil {
					return err
				}

				return os.WriteFile(
					strings.TrimSuffix(writePath, ".plush"),
					content, os.ModePerm,
				)
			} else if d.Type().IsDir() {
				_, err := os.Stat(writePath)
				if err != nil {
					return os.Mkdir(writePath, os.ModeDir|os.ModePerm)
				}
			}

			return nil
		})
}
