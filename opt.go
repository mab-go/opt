// Package opt contains a set of optional types. Optional types provide programmers
// with a better means of distinguishing a value's "is set" state compared to using
// a possibly-nil pointer.
package opt // import "gopkg.in/mab-go/opt.v0"

//go:generate opt-gen gen-data/builtins.json
