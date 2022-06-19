package gofixit

import (
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/internal/enforcer"
	"github.com/sirupsen/logrus"
)

func NewEnforcer(logger *logrus.Logger, config contracts.EnforcerConfig) (contracts.Enforcer, error) {
	return enforcer.NewEnforcer(logger, config)
}
