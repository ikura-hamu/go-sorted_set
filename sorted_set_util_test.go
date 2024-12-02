package gosortedset_test

import "testing"

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func requireTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatalf("expected true, got false")
	}
}
