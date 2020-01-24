package trace

import (
	"fmt"
	"io"
)

// Tracer interface describes an object capable of
// tracing events throughout the code
type Tracer interface {
	Trace(...interface{})
}

// tracer writes to an io.Writer
type tracer struct {
	out io.Writer
}

// Trace writes the arguments to this Tracer's io.Writer
func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// New creates a Tracer that writes the output to the
// specified io.Writer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// nilTracer
type nilTracer struct{}

// Trace for a nilTracer does nothing
func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to Trace
func Off() Tracer {
	return &nilTracer{}
}
