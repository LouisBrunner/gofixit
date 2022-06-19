package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/LouisBrunner/gofixit/src/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type args struct {
	commentPrefixes     []string
	prefixes            []string
	caseInsensitive     bool
	expiryPattern       string
	dateLayout          string
	strict              bool
	recursive           bool
	files               []string
	filesExcludePattern []string
	loggingLevel        logrus.Level
}

func addAllParentFoldersUntilRoot() error {
	start, err := os.Getwd()
	if err != nil {
		return err
	}

	// FIXME: is this dangerous?
	for {
		viper.AddConfigPath(start)
		if start == "/" {
			break
		}
		start = filepath.Dir(start)
	}
	return nil
}

func addDefault[V any](name []string, value V, typ func(name string, value V, help string) *V, help string) {
	configName := strings.Join(utils.MapSlice(name, strings.Title), "")
	viper.SetDefault(configName, value)

	flagName := strings.Join(utils.MapSlice(name, strings.ToLower), "-")
	typ(flagName, value, help)
}

func getArgs() (*args, error) {
	// Default values & flags
	addDefault([]string{"Comment", "Prefixes"}, []string{"//", "#", "/*"}, pflag.StringSlice, "strings which define what a comment definition looks like")
	addDefault([]string{"Prefixes"}, []string{"TODO", "FIXME"}, pflag.StringSlice, "strings which define what a TODO looks like")
	addDefault([]string{"Expiry", "Pattern"}, "{{.Prefix}}(?:\\[{{.Date}}\\])?", pflag.String, "Go template used to generate a regex to match the prefix and expiry date together, careful of escaping any regex character in here")
	addDefault([]string{"Files"}, []string{"."}, pflag.StringSlice, "list of files to parse")
	addDefault([]string{"Files", "Exclude", "Patterns"}, []string{""}, pflag.StringSlice, "list of patterns used to exclude files or directories")
	addDefault([]string{"Recursive"}, true, pflag.Bool, "will process directories recursively")
	addDefault([]string{"Strict"}, false, pflag.Bool, "will force all matched comments to have an expiry date")
	addDefault([]string{"Date", "Layout"}, "2006-01-02", pflag.String, "date layout format, as specified by Golang's date parsing")
	addDefault([]string{"Logging", "Level"}, "fatal", pflag.String, "logrus log level for internal debugging")
	addDefault([]string{"Case", "Sensitive"}, true, pflag.Bool, "should prefixes be matched as case sensitive or not")

	// Env
	viper.SetEnvPrefix("GOFIXIT")
	viper.AutomaticEnv()

	// Config files
	viper.SetConfigName(".gofixit.config")
	viper.SetConfigType("toml")
	err := addAllParentFoldersUntilRoot()
	if err != nil {
		return nil, err
	}
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return nil, err
	}

	// Flags
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Parsing
	logLevel, err := logrus.ParseLevel(viper.GetString("LoggingLevel"))
	if err != nil {
		return nil, err
	}

	return &args{
		commentPrefixes:     viper.GetStringSlice("CommentPrefixes"),
		prefixes:            viper.GetStringSlice("Prefixes"),
		caseInsensitive:     viper.GetBool("CaseSensitive"),
		expiryPattern:       viper.GetString("ExpiryPattern"),
		dateLayout:          viper.GetString("DateLayout"),
		strict:              viper.GetBool("Strict"),
		recursive:           viper.GetBool("Recursive"),
		files:               viper.GetStringSlice("Files"),
		filesExcludePattern: viper.GetStringSlice("FilesExcludePatterns"),
		loggingLevel:        logLevel,
	}, nil
}
