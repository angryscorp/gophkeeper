package sync

import (
	"context"
	"testing"
)

func TestUsernameCtx(t *testing.T) {
	ctx := contextWithUsername(context.Background(), "alice")
	user, ok := usernameFromCtx(ctx)
	if !ok || user != "alice" {
		t.Fatalf("got %q, ok=%v; want alice,true", user, ok)
	}
}
