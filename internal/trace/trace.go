package trace

import "fmt"

type Trace struct {
	enabled bool
	lines   []string
}

func New(enabled bool) *Trace {
	return &Trace{enabled: enabled}
}

func (t *Trace) Add(format string, args ...any) {
	if !t.enabled {
		return
	}
	t.lines = append(t.lines, "[trace] "+fmt.Sprintf(format, args...))
}

func (t *Trace) Lines() []string {
	return t.lines
}

