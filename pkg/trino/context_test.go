package trino

import (
	"context"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestSetQueryArgs_Empty(t *testing.T) {
	ctx := context.Background()
	args := setQueryArgs(ctx)
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestSetQueryArgs_WithUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKeyTrinoUser, &backend.User{Login: "testuser"})
	args := setQueryArgs(ctx)
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
}

func TestSetQueryArgs_WithAccessToken(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKeyAccessToken, "my-token")
	args := setQueryArgs(ctx)
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
}

func TestSetQueryArgs_WithClientTags(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKeyClientTags, "tag1,tag2")
	args := setQueryArgs(ctx)
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
}

func TestSetQueryArgs_AllValues(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyTrinoUser, &backend.User{Login: "user1"})
	ctx = context.WithValue(ctx, ctxKeyAccessToken, "token123")
	ctx = context.WithValue(ctx, ctxKeyClientTags, "tag1,tag2,tag3")

	args := setQueryArgs(ctx)
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
}
