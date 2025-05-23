package runtimev2

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flarelang/flare/internal/ast"
	"github.com/flarelang/flare/internal/cache"
	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/lexer"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/lang"
	"go.uber.org/zap"
)

func (r *Runtime) importer(filename string, dg *models.Debug) (lang.Object, error) {
	zap.L().Info("importing file", zap.String("filename", filename))

	var rootDir string
	if dg != nil && dg.File != "" {
		rootDir = filepath.Dir(dg.File)
	} else {
		rootDir = "." // Use current directory as default
	}

	path := filepath.Join(rootDir, filename)
	path = filepath.Clean(path)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, errs.WithDebug(err, dg)
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(filePath) == ".fl" {
				// Calculate relative path correctly
				relPath, err := filepath.Rel(rootDir, filePath)
				if err != nil {
					return err
				}

				_, err = r.importer(relPath, dg)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, errs.WithDebug(err, dg)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, errs.WithDebug(err, dg)
	}

	builder := ast.NewBuilder()
	nodes, ok := cache.Get(path, b)
	if !ok {
		lx := lexer.New(path)
		ts, err := lx.Parse(bytes.NewReader(b))
		if err != nil {
			return nil, errs.WithDebug(err, dg)
		}
		nodes, err = builder.Build(ts)
		if err != nil {
			return nil, errs.WithDebug(err, dg)
		}
		if len(nodes) == 0 {
			return nil, nil
		}
	}

	cache.Store(path, b, nodes)
	ret, err := r.Execute(nodes)
	if err != nil {
		return nil, errs.WithDebug(err, dg)
	}

	return ret, nil
}

func (r *Runtime) evaler(code string) (lang.Object, error) {
	lx := lexer.New("<eval>")
	ts, err := lx.Parse(strings.NewReader(code))
	if err != nil {
		return nil, errs.WithDebug(err, nil)
	}

	builder := ast.NewBuilder()
	nodes, err := builder.Build(ts)
	if err != nil {
		return nil, errs.WithDebug(err, nil)
	}

	ret, err := r.Execute(nodes)
	if err != nil {
		return nil, errs.WithDebug(err, nil)
	}

	return ret, nil
}
