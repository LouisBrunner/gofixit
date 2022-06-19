# gofixit

A go linter which sets a timer on all your TODOs and FIXMEs and ensure they are dealt in time.

## Installation

```
go install github.com/LouisBrunner/gofixit@latest
```

## Usage

```
gofixit .
```

The program will log all issues to stdout and return status code:

 * `1` if it failed because there was one or more issue
 * `2` when it failed for an internal reason (including when using `-h` or `--help`)

## Configuration

// TODO[2022-06-20]: missing still

## Issues

 * Using date layout with variable amount of digits (e.g. `6` instead of `06`) or with letters (e.g. `Jun` instead of `06`) is currently broken
 * No way to configure the utility
