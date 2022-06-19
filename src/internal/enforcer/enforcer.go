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
			return fmt.Errorf("missing expiry date")
		}
		return nil
	}

	if me.Now.After(*comment.Expiry) {
		duration := me.Now.Sub(*comment.Expiry)
		return fmt.Errorf("now overdue for %s", durafmt.Parse(duration).LimitFirstN(2))
	}
	return nil
}
