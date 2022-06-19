package enforcer

import (
	"testing"
	"time"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
	"github.com/sirupsen/logrus"
)

func Test_Check(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		config  contracts.EnforcerConfig
		comment contracts.ParsedComment
		wantErr bool
	}{
		{
			name: "works, no expiry, not strict",
			config: contracts.EnforcerConfig{
				Now: now,
			},
			comment: contracts.ParsedComment{
				CommentPrefix: "//",
				Prefix:        "fixit",
				Content:       "implement",
				LineNumber:    5,
				OriginalLine:  "// fixit: implement",
			},
		},
		{
			name: "fails, no expiry, strict",
			config: contracts.EnforcerConfig{
				Now:    now,
				Strict: true,
			},
			comment: contracts.ParsedComment{
				CommentPrefix: "//",
				Prefix:        "fixit",
				Content:       "implement",
				LineNumber:    5,
				OriginalLine:  "// fixit: implement",
			},
			wantErr: true,
		},
		{
			name: "works, valid expiry, not strict",
			config: contracts.EnforcerConfig{
				Now: now,
			},
			comment: contracts.ParsedComment{
				CommentPrefix: "//",
				Prefix:        "fixit",
				Content:       "implement",
				Expiry:        utils.Pointerize(now.Add(time.Hour)),
				LineNumber:    5,
				OriginalLine:  "// fixit: implement",
			},
		},
		{
			name: "works, valid expiry, strict",
			config: contracts.EnforcerConfig{
				Now: now,
			},
			comment: contracts.ParsedComment{
				CommentPrefix: "//",
				Prefix:        "fixit",
				Content:       "implement",
				Expiry:        utils.Pointerize(now.Add(time.Hour)),
				LineNumber:    5,
				OriginalLine:  "// fixit: implement",
			},
		},
		{
			name: "fails, invalid expiry, not strict",
			config: contracts.EnforcerConfig{
				Now:    now,
				Strict: true,
			},
			comment: contracts.ParsedComment{
				CommentPrefix: "//",
				Prefix:        "fixit",
				Content:       "implement",
				Expiry:        utils.Pointerize(now.Add(-time.Hour)),
				LineNumber:    5,
				OriginalLine:  "// fixit: implement",
			},
			wantErr: true,
		},
		{
			name: "fails, invalid expiry, strict",
			config: contracts.EnforcerConfig{
				Now:    now,
				Strict: true,
			},
			comment: contracts.ParsedComment{
				CommentPrefix: "//",
				Prefix:        "fixit",
				Content:       "implement",
				Expiry:        utils.Pointerize(now.Add(-time.Hour)),
				LineNumber:    5,
				OriginalLine:  "// fixit: implement",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me, err := New(logrus.New(), tt.config)
			if err != nil {
				t.Fatalf("failed to create enforcer: %v", err)
			}

			err = me.Check(tt.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
