package schema

import "testing"

func TestTransform(t *testing.T) {
	var a = UpdatePayload{
		{"+", "a.b.c/a/b/c"},
		{"+", "a.b.c/a/b/c"},
		{"+", "/a"},
		{"+", ""},
	}
	t.Log(a.Transform())
}
