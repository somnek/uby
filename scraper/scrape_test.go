package main

import "testing"

func TestToNum(t *testing.T) {
	t.Run("plain number", func(t *testing.T) {
		got := ToNum("123")
		want := 123
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("number with comma", func(t *testing.T) {
		got := ToNum("1,234")
		want := 1234
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

}
