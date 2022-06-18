package main

import (
	"fmt"
	"io/ioutil"
	"os"

	gofixit "github.com/LouisBrunner/gofixit/src"
	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/LouisBrunner/gofixit/src/utils"
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
	// TODO[2022-06-26]: read from CLI/env/files
	params := args{
		commentPrefixes: []string{"//", "#", "/*"},
		prefixes:        []string{"TODO", "FIXME"},
		expiryPattern:   "{{.Prefix}}([{{.Date}}])?",
		filesPattern:    []string{"src"}, // TODO[2022-06-19]: wrong default
		// filesPattern:    []string{"."},
		recursive:  true,
		dateLayout: "2006-01-02",
	}

	parser, err := gofixit.NewParser(contracts.ParsingConfig{
		CommentPrefixes: params.commentPrefixes,
		Prefixes:        params.prefixes,
		ExpiryPattern:   params.expiryPattern,
		CaseSensitive:   params.caseSensitive,
		DateLayout:      params.dateLayout,
	})
	if err != nil {
		return false, fmt.Errorf("failed while creating parser (%w)", err)
	}

	enforcer, err := gofixit.NewEnforcer(contracts.EnforcerConfig{
		Strict: params.strict,
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

	processor, err := gofixit.NewFilesProcessor(contracts.FilesProcessorConfig[[]contracts.ParsedComment]{
		Processor:      glue,
		Recursive:      params.recursive,
		FallbackGoList: false, // TODO: expose?
	})
	if err != nil {
		return false, fmt.Errorf("failed while creating processor (%w)", err)
	}

	parsed, err := processor.ProcessFiles(params.filesPattern)
	if err != nil {
		return false, fmt.Errorf("failed while parsing files (%w)", err)
	}

	fmt.Printf("%+v\n", parsed)

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
