# `gofixit`

A go linter which sets a timer on all your TODOs and FIXMEs and ensure they are dealt in time.

## Installation

```
go install github.com/LouisBrunner/gofixit@latest
```

## Usage

```
gofixit
```

The program will log all issues to stdout and return status code:

 * `1` if it failed because there was one or more issue
 * `2` when it failed for an internal reason (including when using `-h` or `--help`)

### Configuration

`gofixit` supports getting settings through:

 * Environment variables: all env vars must be prefixed with `GOFIXIT_` and be in all caps
 * Configuration file: any file named `.gofixit.config.toml` in the current directory or any parent will be used, it must be a TOML formatted file
 * Command-line arguments: flag names use `kebab-case` instead of `PascalCase` as for the configuration file

All settings can be set through any of those sources, they are ordered by ascending priority (environment variables < configuration file < command line arguments).

Settings:

 * `CommentPrefixes`: strings which define what a comment definition looks like (default `[//,#,/*]`)
 * `Prefixes`: strings which define what a TODO looks like (default `[TODO,FIXME]`)
 * `CaseSensitive`: should prefixes be matched as case sensitive or not (default `true`)
 * `ExpiryPattern`: Go template used to generate a regex to match the prefix and expiry date together, careful of escaping any regex character in here (default `"{{.Prefix}}(?:\\[{{.Date}}\\])?"`), see [here](https://pkg.go.dev/text/template) for details about Go templating and [here](https://github.com/google/re2/wiki/Syntax) for details about Go regex
 * `DateLayout`: date layout format, as specified by Golang's date parsing (default `"2006-01-02"`), see [here](https://pkg.go.dev/time#Parse) for more details about format
 * `Strict`: will force all matched comments to have an expiry date
 * `Recursive`: will process directories recursively (default `true`)
 * `Files`: list of files to parse (default `[.]`)
 * `LoggingLevel`: logrus log level for internal debugging (default `"fatal"`)

Example:

```bash
# CommentPrefixes can be set through all these mechanisms

# Environment variable
GOFIXIT_COMMENTPREFIXES='//,/*' gofixit

# Configuration file
cat > .gofixit.config.toml <<HEREDOC
CommentPrefixes = ['//', '/*']
HEREDOC
gofixit

# Command-line argument
gofixit --comment-prefixes='//,/*'
```


## Issues

 * Using date layout with variable amount of digits (e.g. `6` instead of `06`) or with letters (e.g. `Jun` instead of `06`) is currently broken
 * No way to configure the utility
