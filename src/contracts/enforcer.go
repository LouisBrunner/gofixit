package contracts

import "time"

type EnforcerConfig struct {
	Now    time.Time
	Strict bool
}

type Enforcer interface {
	Check(comment ParsedComment) error
}
