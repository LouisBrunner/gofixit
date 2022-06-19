package parser

import (
	"testing"
	"time"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	tests := []struct {
		name    string
		config  contracts.ParsingConfig
		wantErr bool
	}{
		{
			name: "fails, template has no .Prefix",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"//"},
				Prefixes:        []string{"TODO"},
				ExpiryPattern:   "~{{.Date}}~",
				DateLayout:      "02/01/2006",
				CaseSensitive:   true,
			},
			wantErr: true,
		},
		{
			name: "fails, template has no .Date",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"//"},
				Prefixes:        []string{"TODO"},
				ExpiryPattern:   "{{.Prefix}}:",
				DateLayout:      "02/01/2006",
				CaseSensitive:   true,
			},
			wantErr: true,
		},
		{
			name: "template is invalid",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"//"},
				Prefixes:        []string{"TODO"},
				ExpiryPattern:   "@{{.Prefix}->{{.Date}}:",
				DateLayout:      "02/01/2006",
				CaseSensitive:   true,
			},
			wantErr: true,
		},
		{
			name: "date layout is unsupported",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"//"},
				Prefixes:        []string{"TODO"},
				ExpiryPattern:   "{{.Prefix}}(?:->{{.Date}})?:",
				DateLayout:      "02/Jan/2006",
				CaseSensitive:   true,
			},
			wantErr: true,
		},
		{
			name: "works",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"//"},
				Prefixes:        []string{"TODO"},
				ExpiryPattern:   "{{.Prefix}}(?:->{{.Date}})?:",
				DateLayout:      "02/01/2006",
				CaseSensitive:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(logrus.New(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Parse(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	tests := []struct {
		name        string
		config      contracts.ParsingConfig
		fileContent string
		want        []contracts.ParsedComment
		wantErr     bool
	}{
		{
			name: "works",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"@", "%"},
				Prefixes:        []string{"fixit", "later"},
				ExpiryPattern:   "{{.Prefix}}(?:->{{.Date}})?:",
				DateLayout:      "02/01/2006",
				CaseSensitive:   true,
			},
			fileContent: `
#include <stdio>

@fixit: made with love
do {
	sleep(1);
} while (1); @fixit->30/04/1993: better condition?

@ later: don't match me
%later : or me
%later ->19/06/2022: or me

@later: finish when I have time
%later: bit too busy right now
`,
			want: []contracts.ParsedComment{
				{
					CommentPrefix: "@",
					Prefix:        "fixit",
					Content:       "made with love",
					LineNumber:    4,
					OriginalLine:  "@fixit: made with love",
				},
				{
					CommentPrefix: "@",
					Prefix:        "fixit",
					Content:       "better condition?",
					Expiry:        utils.Pointerize(utils.Must(time.Parse("02/01/2006", "30/04/1993"))),
					LineNumber:    7,
					OriginalLine:  "@fixit->30/04/1993: better condition?",
				},
				{
					CommentPrefix: "@",
					Prefix:        "later",
					Content:       "don't match me",
					LineNumber:    9,
					OriginalLine:  "@ later: don't match me",
				},
				{
					CommentPrefix: "@",
					Prefix:        "later",
					Content:       "finish when I have time",
					LineNumber:    13,
					OriginalLine:  "@later: finish when I have time",
				},
				{
					CommentPrefix: "%",
					Prefix:        "later",
					Content:       "bit too busy right now",
					LineNumber:    14,
					OriginalLine:  "%later: bit too busy right now",
				},
			},
		},
		{
			name: "works, invalid date",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"@", "%"},
				Prefixes:        []string{"fixit", "later"},
				ExpiryPattern:   "{{.Prefix}}(?:->{{.Date}})?:",
				DateLayout:      "02/01/2006",
				CaseSensitive:   true,
			},
			fileContent: `
#include <stdio>

@fixit->99/04/1993: forgot to implement
%later: maybe after dinner
`,
			want: []contracts.ParsedComment{
				{
					CommentPrefix: "%",
					Prefix:        "later",
					Content:       "maybe after dinner",
					LineNumber:    5,
					OriginalLine:  "%later: maybe after dinner",
				},
			},
		},
		{
			name: "works, inverted .Prefix/.Date",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"@", "%"},
				Prefixes:        []string{"fixit", "later"},
				ExpiryPattern:   "{{.Date}}\\.{{.Prefix}}",
				DateLayout:      "02/01/2006",
				CaseSensitive:   true,
			},
			fileContent: `
#include <stdio>

@18/05/1991.fixit forgot to implement
%.later maybe after dinner
`,
			want: []contracts.ParsedComment{
				{
					CommentPrefix: "@",
					Prefix:        "fixit",
					Content:       "forgot to implement",
					Expiry:        utils.Pointerize(utils.Must(time.Parse("02/01/2006", "18/05/1991"))),
					LineNumber:    4,
					OriginalLine:  "@18/05/1991.fixit forgot to implement",
				},
				{
					CommentPrefix: "%",
					Prefix:        "later",
					Content:       "maybe after dinner",
					LineNumber:    5,
					OriginalLine:  "%.later maybe after dinner",
				},
			},
		},
		{
			name: "works, case insensitive",
			config: contracts.ParsingConfig{
				CommentPrefixes: []string{"@", "%"},
				Prefixes:        []string{"fixit", "later"},
				ExpiryPattern:   "{{.Prefix}}(?:->{{.Date}})?:",
				DateLayout:      "02/01/2006",
			},
			fileContent: `
#include <stdio>

@FixIt->18/05/1991: forgot to implement
%LATER: maybe after dinner

@later: MORE


%FiXiT: EVEN MORE
`,
			want: []contracts.ParsedComment{
				{
					CommentPrefix: "@",
					Prefix:        "FixIt",
					Content:       "forgot to implement",
					Expiry:        utils.Pointerize(utils.Must(time.Parse("02/01/2006", "18/05/1991"))),
					LineNumber:    4,
					OriginalLine:  "@FixIt->18/05/1991: forgot to implement",
				},
				{
					CommentPrefix: "%",
					Prefix:        "LATER",
					Content:       "maybe after dinner",
					LineNumber:    5,
					OriginalLine:  "%LATER: maybe after dinner",
				},
				{
					CommentPrefix: "@",
					Prefix:        "later",
					Content:       "MORE",
					LineNumber:    7,
					OriginalLine:  "@later: MORE",
				},
				{
					CommentPrefix: "%",
					Prefix:        "FiXiT",
					Content:       "EVEN MORE",
					LineNumber:    10,
					OriginalLine:  "%FiXiT: EVEN MORE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me, err := New(logger, tt.config)
			if err != nil {
				t.Fatalf("failed to create parser: %v", err)
			}

			res, err := me.Parse(tt.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, res)
		})
	}
}
