package enforcer

import (
	"fmt"

	"github.com/LouisBrunner/gofixit/src/contracts"
)

type enforcer struct {
	contracts.EnforcerConfig
}

func NewEnforcer(config contracts.EnforcerConfig) (contracts.Enforcer, error) {
	return &enforcer{
		EnforcerConfig: config,
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
		// TODO: use some nice library to render those durations
		return fmt.Errorf("now overdue for %s", me.Now.Sub(*comment.Expiry))
	}
	return nil
}
