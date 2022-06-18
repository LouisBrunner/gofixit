package gofixit

import (
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/internal/enforcer"
)

func NewEnforcer(config contracts.EnforcerConfig) (contracts.Enforcer, error) {
	return enforcer.NewEnforcer(config)
}
