package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
	"github.com/sirupsen/logrus"
)

const maxPreAllocated = 100

type parserImpl struct {
	contracts.ParsingConfig
	re     regexp.Regexp
	order  ordering
	logger *logrus.Logger
}

func New(logger *logrus.Logger, config contracts.ParsingConfig) (contracts.Parser, error) {
	re, order, err := buildRE(logger, config)
	if err != nil {
		return nil, fmt.Errorf("cannot build internal matcher: %w", err)
	}

	return &parserImpl{
		ParsingConfig: config,
		re:            *re,
		order:         *order,
		logger:        logger,
	}, nil
}

func (me *parserImpl) Parse(fileContent string) ([]contracts.ParsedComment, error) {
	lines := strings.Split(fileContent, "\n")

	results := make([]contracts.ParsedComment, 0, utils.Min(len(lines)/10, maxPreAllocated))
	for num, line := range lines {
		me.logger.Debugf("parsing line %q", line)
		matches := me.re.FindStringSubmatch(line)
		if len(matches) < expectedMatches {
			continue
		}
		me.logger.Infof("line matched %+v", matches)
		var expiry *time.Time
		if matches[me.order.matchExpiry] != "" {
			expiryValue, err := time.Parse(me.DateLayout, matches[me.order.matchExpiry])
			if err != nil {
				me.logger.Errorf("invalid date layout %q: %v", matches[me.order.matchExpiry], err)
				continue
			}
			expiry = &expiryValue
		}
		results = append(results, contracts.ParsedComment{
			CommentPrefix: matches[me.order.matchComment],
			Prefix:        matches[me.order.matchPrefix],
			Content:       matches[me.order.matchContent],
			Expiry:        expiry,
			LineNumber:    uint(num) + 1,
			OriginalLine:  strings.TrimSpace(line),
		})
	}
	return results, nil
}
