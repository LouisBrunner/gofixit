package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	gofixit "github.com/LouisBrunner/gofixit/src"
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
	"github.com/sirupsen/logrus"
)

type args struct {
	commentPrefixes []string
	prefixes        []string
	caseSensitive   bool
	expiryPattern   string
	dateLayout      string
	strict          bool
	recursive       bool
	filesPattern    []string
	loggingLevel    logrus.Level
}

func main() {
	success, err := work()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gofixit: %s\n", err.Error())
		os.Exit(1)
	}
	if !success {
		os.Exit(2)
	}
}

func work() (bool, error) {
	// TODO[2022-06-20]: read from CLI/env/files
	params := args{
		commentPrefixes: []string{"//", "#", "/*"},
		prefixes:        []string{"TODO", "FIXME"},
		expiryPattern:   "{{.Prefix}}(?:\\[{{.Date}}\\])?",
		filesPattern:    []string{"."},
		recursive:       true,
		dateLayout:      "2006-01-02",
		loggingLevel:    logrus.FatalLevel,
	}

	log := logrus.New()
	log.SetLevel(params.loggingLevel)
	log.SetOutput(os.Stderr)

	parser, err := gofixit.NewParser(log, contracts.ParsingConfig{
		CommentPrefixes: params.commentPrefixes,
		Prefixes:        params.prefixes,
		ExpiryPattern:   params.expiryPattern,
		CaseSensitive:   params.caseSensitive,
		DateLayout:      params.dateLayout,
	})
	if err != nil {
		return false, fmt.Errorf("failed while creating parser (%w)", err)
	}

	enforcer, err := gofixit.NewEnforcer(log, contracts.EnforcerConfig{
		Strict: params.strict,
		Now:    time.Now(),
	})
	if err != nil {
		return false, fmt.Errorf("failed while creating enforcer (%w)", err)
	}

	glue := func(filepath string) ([]contracts.ParsedComment, error) {
		// FIXME: should really be streaming files better than this
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		return parser.Parse(string(content))
	}

	processor, err := gofixit.NewFilesProcessor(log, contracts.FilesProcessorConfig[[]contracts.ParsedComment]{
		Processor: glue,
		Recursive: params.recursive,
	})
	if err != nil {
		return false, fmt.Errorf("failed while creating processor (%w)", err)
	}

	parsed, err := processor.ProcessFiles(params.filesPattern)
	if err != nil {
		return false, fmt.Errorf("failed while parsing files (%w)", err)
	}

	hadError := false
	for _, entry := range utils.SortedMap(parsed) {
		for _, comment := range entry.Value {
			err := enforcer.Check(comment)
			if err == nil {
				continue
			}
			fmt.Printf("%s:%d %s\n", entry.Key, comment.LineNumber, err.Error())
			hadError = true
		}
	}
	return !hadError, nil
}
