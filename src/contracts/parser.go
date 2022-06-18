package contracts

import "time"

type ParsedComment struct {
	Prefix       string
	Content      string
	Expiry       *time.Time
	LineNumber   uint
	OriginalLine string
}

type ParsingConfig struct {
	Prefixes      []string
	ExpiryPattern string
	CaseSensitive bool
}

type Parser interface {
	Parse(fileContent string) ([]ParsedComment, error)
}
