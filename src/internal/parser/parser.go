package parser

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
)

const maxPreAllocated = 100

const (
	// TODO[2022-06-19]: this breaks if the template reorder the Prefix and Date
	matchComment = iota + 1
	matchPrefix
	matchExpiry
	matchContent
	expectedMatches
)

type parserImpl struct {
	contracts.ParsingConfig
	re regexp.Regexp
}

func New(config contracts.ParsingConfig) (contracts.Parser, error) {
	re, err := buildRE(config)
	if err != nil {
		return nil, fmt.Errorf("cannot build internal matcher: %w", err)
	}

	return &parserImpl{
		ParsingConfig: config,
		re:            *re,
	}, nil
}

func layoutToRegex(layout string) string {
	matcherBuilder := &strings.Builder{}
	for _, c := range layout {
		// TODO[2022-10-04]: don't support layout with anything but digits and separators
		if unicode.IsDigit(c) {
			matcherBuilder.WriteString("[[:digit:]]")
		} else {
			matcherBuilder.WriteString(regexp.QuoteMeta(string(c)))
		}
	}
	return matcherBuilder.String()
}

func buildRE(config contracts.ParsingConfig) (*regexp.Regexp, error) {
	tmpl, err := template.New("partWithExpiry").Parse(config.ExpiryPattern)
	if err != nil {
		return nil, err
	}

	patternBuilder := &strings.Builder{}
	tmpl.Execute(patternBuilder, struct {
		Prefix string
		Date   string
	}{
		Prefix: fmt.Sprintf("(%s)", strings.Join(utils.MapSlice(config.Prefixes, regexp.QuoteMeta), "|")),
		Date:   layoutToRegex(config.DateLayout),
	})

	literal := fmt.Sprintf(
		"(%s)[[:space:]]*%s[[:space:]]*(.+)?$",
		strings.Join(utils.MapSlice(config.CommentPrefixes, regexp.QuoteMeta), "|"),
		patternBuilder.String(),
	)
	// TODO: add the logging for the regexp

	return regexp.Compile(literal)
}

func (me *parserImpl) Parse(fileContent string) ([]contracts.ParsedComment, error) {
	lines := strings.Split(fileContent, "\n")

	results := make([]contracts.ParsedComment, 0, utils.Min(len(lines)/10, maxPreAllocated))
	for num, line := range lines {
		matches := me.re.FindStringSubmatch(line)
		if len(matches) < expectedMatches {
			continue
		}
		var expiry *time.Time
		if matches[matchExpiry] != "" {
			expiryValue, err := time.Parse(me.DateLayout, matches[matchExpiry])
			if err != nil {
				// TODO: add logging?
				continue
			}
			expiry = &expiryValue
		}
		results = append(results, contracts.ParsedComment{
			CommentPrefix: matches[matchComment],
			Prefix:        matches[matchPrefix],
			Content:       matches[matchContent],
			Expiry:        expiry,
			LineNumber:    uint(num) + 1,
			OriginalLine:  strings.TrimSpace(line),
		})
	}
	return results, nil
}
