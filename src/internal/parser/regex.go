package parser

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"unicode"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
	"github.com/sirupsen/logrus"
)

func layoutToRegex(layout string) (string, error) {
	matcherBuilder := &strings.Builder{}
	for _, c := range layout {
		// TODO[2022-08-10]: don't support layout with anything but digits and separators
		if unicode.IsDigit(c) {
			matcherBuilder.WriteString("[[:digit:]]")
		} else if unicode.IsLetter(c) {
			return "", fmt.Errorf("unsupported character %q in the date layout", c)
		} else {
			matcherBuilder.WriteString(regexp.QuoteMeta(string(c)))
		}
	}
	return matcherBuilder.String(), nil
}

const (
	matchEverything = iota + 1
	matchComment
	matchPrefix
	matchExpiry
	matchContent
	expectedMatches
)

type ordering struct {
	matchEverything int
	matchComment    int
	matchPrefix     int
	matchExpiry     int
	matchContent    int
}

func buildRE(logger *logrus.Logger, config contracts.ParsingConfig) (*regexp.Regexp, *ordering, error) {
	order := ordering{
		matchEverything: matchEverything,
		matchComment:    matchComment,
		matchPrefix:     matchPrefix,
		matchExpiry:     matchExpiry,
		matchContent:    matchContent,
	}

	idxPrefix := strings.Index(config.ExpiryPattern, ".Prefix")
	if idxPrefix == -1 {
		return nil, nil, fmt.Errorf("expiry template must contain {{.Prefix}}")
	}
	idxExpiry := strings.Index(config.ExpiryPattern, ".Date")
	if idxExpiry == -1 {
		return nil, nil, fmt.Errorf("expiry template must contain {{.Date}}")
	}
	if idxExpiry < idxPrefix {
		order.matchPrefix = matchExpiry
		order.matchExpiry = matchPrefix
	}

	tmpl, err := template.New("partWithExpiry").Parse(config.ExpiryPattern)
	if err != nil {
		return nil, nil, err
	}

	dateRegex, err := layoutToRegex(config.DateLayout)
	if err != nil {
		return nil, nil, err
	}

	patternBuilder := &strings.Builder{}
	err = tmpl.Execute(patternBuilder, struct {
		Prefix string
		Date   string
	}{
		Prefix: fmt.Sprintf("(%s)", strings.Join(utils.MapSlice(config.Prefixes, regexp.QuoteMeta), "|")),
		Date:   fmt.Sprintf("((?:%s)?)", dateRegex),
	})
	if err != nil {
		return nil, nil, err
	}

	flags := "(?i)"
	if config.CaseSensitive {
		flags = "(?)"
	}

	literal := fmt.Sprintf(
		"%s((%s)[[:space:]]*%s[[:space:]]*(.+)?)$",
		flags,
		strings.Join(utils.MapSlice(config.CommentPrefixes, regexp.QuoteMeta), "|"),
		patternBuilder.String(),
	)
	logger.Infof("using regex %q to parse comments", literal)

	re, err := regexp.Compile(literal)
	return re, &order, err
}
