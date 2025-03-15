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
	return fs.WalkDir(
		staticFiles, "static",
		func(path string, d fs.DirEntry, err error) error {
			name, err := filepath.Rel(
				"static", path,
			)
			if err != nil {
				return err
			}
			writePath := filepath.Join(g.outPath, name)

			if d.Type().IsRegular() {
				content, err := fs.ReadFile(staticFiles, path)
				if err != nil {
					return err
				}

				return os.WriteFile(
					strings.TrimRight(writePath, ".plush"),
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
