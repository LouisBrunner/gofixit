package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
	"github.com/sirupsen/logrus"
)

type fprocessor[T any] struct {
	contracts.FilesProcessorConfig[T]
	logger          *logrus.Logger
	excludePatterns []regexp.Regexp
}

func New[T any](logger *logrus.Logger, config contracts.FilesProcessorConfig[T]) (contracts.FilesProcessor[T], error) {
	excludePatterns := make([]regexp.Regexp, 0, len(config.FilesExcludePatterns))
	for _, excludePattern := range config.FilesExcludePatterns {
		re, err := regexp.Compile(excludePattern)
		if err != nil {
			return nil, err
		}
		excludePatterns = append(excludePatterns, *re)
	}

	return &fprocessor[T]{
		FilesProcessorConfig: config,
		logger:               logger,
		excludePatterns:      excludePatterns,
	}, nil
}

func (me *fprocessor[T]) ProcessFiles(files []string) (map[string]T, error) {
	results := make(map[string]T, len(files))
	absMatches := make(map[string]struct{}, len(files))

	extras := []string{}
	for len(files) > 0 {
	fileLoop:
		for _, filename := range files {
			absFilename, err := filepath.Abs(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to generate absolute path for %s: %w", filename, err)
			}
			if _, found := absMatches[absFilename]; found {
				continue
			}

			for _, needle := range []string{filename, absFilename} {
				for _, excludePattern := range me.excludePatterns {
					if excludePattern.MatchString(needle) {
						continue fileLoop
					}
				}
			}

			info, err := os.Stat(filename)
			if err == nil {
				if !info.IsDir() {
					result, err := me.Processor(filename)
					if err != nil {
						return nil, fmt.Errorf("failed to process %s: %w", filename, err)
					}
					results[filename] = result
					absMatches[absFilename] = struct{}{}
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

			return nil, fmt.Errorf("no such file %s: %w", filename, err)
		}
		files = extras
		extras = []string{}
	}
	return results, nil
}
