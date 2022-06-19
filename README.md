# `gofixit`

You work on a new task, encounter a difficult problem or something you don't want to implement yet, you drop a `TODO` and move on, then some time later production break because you forgot about it... has this ever happened to you?

`gofixit` is a language-agnostic linter which enforces that your TODOs and FIXMEs are dealt with in time, very similar to [this eslint plugin](https://github.com/sindresorhus/eslint-plugin-unicorn/blob/main/docs/rules/expiring-todo-comments.md).

## Example

File (`examples/example1.c`):
```c
#include <stdio.h>

void doSomething() {
  // TODO[2022-06-15]: implement later
  assert(("unimplemented", 0));
}


int main() {
  doSomething()
}
```

Result (with `gofixit --files examples/example1.c --strict` as of 2022-06-19):
```
examples/example1.c:4 TODO now overdue for 4 days 13 hours
examples/example1.c:12 FIXME missing expiry date
```

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
 * `FilesExcludePatterns`: list of patterns used to exclude files or directories
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
