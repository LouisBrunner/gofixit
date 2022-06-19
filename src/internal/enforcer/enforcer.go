package enforcer

import (
	"fmt"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/hako/durafmt"
	"github.com/sirupsen/logrus"
)

type enforcer struct {
	contracts.EnforcerConfig
	logger *logrus.Logger
}

func New(logger *logrus.Logger, config contracts.EnforcerConfig) (contracts.Enforcer, error) {
	return &enforcer{
		EnforcerConfig: config,
		logger:         logger,
	}, nil
}

func (me *enforcer) Check(comment contracts.ParsedComment) error {
	if comment.Expiry == nil {
		if me.Strict {
			return fmt.Errorf("%s missing expiry date", comment.Prefix)
		}
		return nil
	}

	if me.Now.After(*comment.Expiry) {
		duration := me.Now.Sub(*comment.Expiry)
		return fmt.Errorf("%s now overdue for %s", comment.Prefix, durafmt.Parse(duration).LimitFirstN(2))
	}
	return nil
}
