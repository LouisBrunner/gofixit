package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
)

type fprocessor[T any] struct {
	contracts.FilesProcessorConfig[T]
}

func New[T any](config contracts.FilesProcessorConfig[T]) (contracts.FilesProcessor[T], error) {
	return &fprocessor[T]{
		FilesProcessorConfig: config,
	}, nil
}

func (me *fprocessor[T]) ProcessFiles(files []string) (map[string]T, error) {
	results := make(map[string]T, len(files))

	extras := []string{}
	for len(files) > 0 {
		for _, filename := range files {
			filename, err := filepath.Abs(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to generate absolute path for %s: %w", filename, err)
			}
			if _, found := results[filename]; found {
				continue
			}

			info, err := os.Stat(filename)
			if err == nil {
				if !info.IsDir() {
					result, err := me.Processor(filename)
					if err != nil {
						return nil, fmt.Errorf("failed to process %s: %w", filename, err)
					}
					results[filename] = result
				} else if me.Recursive {
					files, err := os.ReadDir(filename)
					if err != nil {
						return nil, fmt.Errorf("failed to list %s: %w", filename, err)
					}
					filenames := utils.MapSlice(files, func(entry fs.DirEntry) string {
						return filepath.Join(filename, entry.Name())
					})
					extras = append(extras, filenames...)
				} else {
					return nil, fmt.Errorf("cannot process directory %s: %w", filename, err)
				}

				continue
			}

			if !me.FallbackGoList {
				return nil, fmt.Errorf("no such file %s: %w", filename, err)
			}

			// TODO: find a way to be able to support ./... and such
			panic("unimplemented")
		}
		files = extras
		extras = []string{}
	}
	return results, nil
}
