package gofixit

import (
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/internal/parser"
	"github.com/sirupsen/logrus"
)

func NewParser(logger *logrus.Logger, config contracts.ParsingConfig) (contracts.Parser, error) {
	return parser.New(logger, config)
}
