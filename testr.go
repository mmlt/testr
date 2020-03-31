// Package testr implements github.com/go-logr/logr.Logger in terms of
// Go's test package log method.
package testr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-logr/logr"
)

// The global verbosity level.  See SetVerbosity().
var globalVerbosity int = 0

// SetVerbosity sets the global level against which all info logs will be
// compared.  If this is greater than or equal to the "V" of the logger, the
// message will be logged.  A higher value here means more logs will be written.
// The previous verbosity value is returned.  This is not concurrent-safe -
// callers must be sure to call it from only one goroutine.
func SetVerbosity(v int) int {
	old := globalVerbosity
	globalVerbosity = v
	return old
}

// New returns a logr.Logger which is implemented by Go's test package Log().
//
// Example: testr.New(t)
func New(std TestLogger) logr.Logger {
	return logger{
		std:    std,
		level:  0,
		prefix: "",
		values: nil,
	}
}

// TestLogger is the subset of the Go test package Log API that is needed for
// this adapter.
type TestLogger interface {
	// Log matches https://pkg.go.dev/testing?tab=doc#B.Log
	Log(args ...interface{})
}

type logger struct {
	std    TestLogger
	level  int
	prefix string
	values []interface{}
	//depth  int
}

func (l logger) clone() logger {
	out := l
	l.values = copySlice(l.values)
	return out
}

func copySlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	copy(out, in)
	return out
}

func flatten(kvList ...interface{}) string {
	keys := make([]string, 0, len(kvList))
	vals := make(map[string]interface{}, len(kvList))
	for i := 0; i < len(kvList); i += 2 {
		k, ok := kvList[i].(string)
		if !ok {
			panic(fmt.Sprintf("key is not a string: %s", pretty(kvList[i])))
		}
		var v interface{}
		if i+1 < len(kvList) {
			v = kvList[i+1]
		}
		keys = append(keys, k)
		vals[k] = v
	}
	sort.Strings(keys)
	buf := bytes.Buffer{}
	for i, k := range keys {
		v := vals[k]
		if i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(pretty(k))
		buf.WriteString("=")
		buf.WriteString(pretty(v))
	}
	return buf.String()
}

func pretty(value interface{}) string {
	jb, _ := json.Marshal(value)
	return string(jb)
}

func (l logger) Info(msg string, kvList ...interface{}) {
	if l.Enabled() {
		lvlStr := flatten("level", l.level)
		msgStr := flatten("msg", msg)
		fixedStr := flatten(l.values...)
		userStr := flatten(kvList...)
		l.std.Log(l.prefix, lvlStr, msgStr, fixedStr, userStr)
	}
}

func (l logger) Enabled() bool {
	return globalVerbosity >= l.level
}

func (l logger) Error(err error, msg string, kvList ...interface{}) {
	msgStr := flatten("msg", msg)
	var loggableErr interface{}
	if err != nil {
		loggableErr = err.Error()
	}
	errStr := flatten("error", loggableErr)
	fixedStr := flatten(l.values...)
	userStr := flatten(kvList...)
	l.std.Log(l.prefix, errStr, msgStr, fixedStr, userStr)
}

func (l logger) V(level int) logr.InfoLogger {
	r := l.clone()
	r.level = level
	return r
}

// WithName returns a new logr.Logger with the specified name appended.
// testr uses '/' characters to separate name elements.  Callers should not pass '/'
// in the provided name string, but this library does not actually enforce that.
func (l logger) WithName(name string) logr.Logger {
	r := l.clone()
	if len(l.prefix) > 0 {
		r.prefix = l.prefix + "/"
	}
	r.prefix += name
	return r
}

func (l logger) WithValues(kvList ...interface{}) logr.Logger {
	r := l.clone()
	r.values = append(r.values, kvList...)
	return r
}
