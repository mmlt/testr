# Minimal Go logging using logr and Go's standard library

This package implements the [logr interface](https://github.com/go-logr/logr)
in terms of Go's test package [log methods](https://pkg.go.dev/testing?tab=doc#B.Log).

Pro's and con's
- Log output is only shown for failing tests.
- Filename/line number printed in the log message is always the same (GO test package Log() can't skip stacklevels).
