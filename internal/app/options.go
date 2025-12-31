package app

import "portarbiter/internal/trace"

type Options struct {
	Port   int
	DryRun bool
	Kill   bool
	Force  bool
	Yes    bool

	Trace *trace.Trace
}

