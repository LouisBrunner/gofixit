package contracts

import "time"

type ParsedComment struct {
	CommentPrefix string
	Prefix        string
	Content       string
	Expiry        *time.Time
	LineNumber    uint
	OriginalLine  string
}

type ParsingConfig struct {
	CommentPrefixes []string
	Prefixes        []string
	ExpiryPattern   string
	CaseSensitive   bool
	DateLayout      string
}

type Parser interface {
	Parse(fileContent string) ([]ParsedComment, error)
}
