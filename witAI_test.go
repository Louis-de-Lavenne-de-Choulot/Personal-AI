package main

import (
	"os"
	"testing"
)

// test message to wit.ai
func TestWitAIHandler(t *testing.T) {
	//get secretKey.txt content
	secretKey, err := os.ReadFile("secretKey.txt")
	if err != nil {
		t.Errorf("Error reading secretKey.txt: %v", err)
	}
	t.Run("test wit.ai handling", func(t *testing.T) {
		s := "Hello"
		a := witAIHandler(s, string(secretKey))
		final := "Hello"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test wit.ai handling", func(t *testing.T) {
		s := "Thanks"
		a := witAIHandler(s, string(secretKey))
		final := "You're welcome"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test wit.ai handling", func(t *testing.T) {
		s := "Bye"
		a := witAIHandler(s, string(secretKey))
		final := "Bye"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
}
