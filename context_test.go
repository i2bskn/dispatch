package pygmy

import (
	"context"
	"testing"
)

func TestContextShare(t *testing.T) {
	want := newShare("/")
	ctx := context.Background()

	obj := getShare(ctx)
	if obj != nil {
		t.Fatalf("got: nil\nwant: %#v", obj)
	}

	ctx = setShare(ctx, want)
	obj = getShare(ctx)
	if obj.path != want.path {
		t.Fatalf("got: %#v\nwant: %#v", obj.path, want.path)
	}
}

func TestContextParam(t *testing.T) {
	want := map[string]string{"id": "1234"}
	ctx := context.Background()

	id := Param(ctx, "id")
	if id != "" {
		t.Fatalf("got: %#v\nwant empty string", id)
	}

	ctx = setParam(ctx, want)
	id = Param(ctx, "id")
	if id != want["id"] {
		t.Fatalf("got: %#v\nwant: %#v", id, want["id"])
	}
}
