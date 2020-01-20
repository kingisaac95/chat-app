package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buffer bytes.Buffer
	tracer := New(&buffer)

	if tracer == nil {
		t.Error("New should not return nil")
	} else {
		tracer.Trace("Hello from trace package.")
		if buffer.String() != "Hello from trace package.\n" {
			t.Errorf("Trace should not write '%s'.", buffer.String())
		}
	}
}
