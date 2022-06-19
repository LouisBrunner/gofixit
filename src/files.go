package gofixit

import (
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/internal/files"
	"github.com/sirupsen/logrus"
)

func NewFilesProcessor[T any](logger *logrus.Logger, config contracts.FilesProcessorConfig[T]) (contracts.FilesProcessor[T], error) {
	return files.New(logger, config)
}
