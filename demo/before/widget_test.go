// Demo only — not compiled by CI. Represents the tests/ code BEFORE the bump,
// which no longer compiles against github.com/example/widget v2 because
// widget.New now requires a context.Context.
package demo_test

import (
	"testing"

	"github.com/example/widget"
)

func TestWidgetGreets(t *testing.T) {
	w := widget.New() // v2: not enough arguments — wants (context.Context)
	if got := w.Greet("world"); got != "hello world" {
		t.Fatalf("unexpected greeting: %q", got)
	}
}
