package main

import "testing"

func TestLanguageSearch(t *testing.T) {
	t.Run("test language search", func(t *testing.T) {
		s := "french"
		a := languageSearch(s, false, "languagesSpeech")
		final := "fr"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test language search", func(t *testing.T) {
		s := "spanish"
		a := languageSearch(s, true, "languagesTranslation")
		final := "es"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test language search", func(t *testing.T) {
		s := "fren"
		a := languageSearch(s, false, "languagesSpeech")
		final := "fr"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test language search", func(t *testing.T) {
		s := "spani"
		a := languageSearch(s, true, "languagesTranslation")
		final := "es"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
}
