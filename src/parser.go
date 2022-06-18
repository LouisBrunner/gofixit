package gofixit

import (
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/internal/parser"
)

func NewParser(config contracts.ParsingConfig) (contracts.Parser, error) {
	return parser.New(config)
}
