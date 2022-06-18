package gofixit

import (
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/internal/files"
)

func NewFilesProcessor[T any](config contracts.FilesProcessorConfig[T]) (contracts.FilesProcessor[T], error) {
	return files.New(config)
}
