// Demo only — not compiled by CI. Represents the minimal fix the CI Fixer
// applies AFTER the bump: pass a context.Context, staying within the allowed
// files (tests/**/*.go). No test assertions were changed.
package demo_test

import (
	"context"
	"testing"

	"github.com/example/widget"
)

func TestWidgetGreets(t *testing.T) {
	w := widget.New(context.Background()) // v2 signature satisfied
	if got := w.Greet("world"); got != "hello world" {
		t.Fatalf("unexpected greeting: %q", got)
	}
}
