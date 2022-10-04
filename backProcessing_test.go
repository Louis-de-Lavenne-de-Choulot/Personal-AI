package main

import "testing"

func TestLanguageSearch(t *testing.T) {
	t.Run("test language search", func(t *testing.T) {
		s := "french"
		a := languageSearch(s, "languagesSpeech")
		final := "fr"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test language search", func(t *testing.T) {
		s := "spanish"
		a := languageSearch(s, "languagesTranslation")
		final := "es"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test language search", func(t *testing.T) {
		s := "fren"
		a := languageSearch(s, "languagesSpeech")
		final := "fr"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test language search", func(t *testing.T) {
		s := "spani"
		a := languageSearch(s, "languagesTranslation")
		final := "es"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
}

func TestTranslate(t *testing.T) {
	t.Run("test translation", func(t *testing.T) {
		s := "Hello"
		a := Translate(s, "fr", "en")
		final := "Bonjour"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test translation", func(t *testing.T) {
		s := "Thank you"
		a := Translate(s, "en", "fr")
		final := "Thank you"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
	t.Run("test translation", func(t *testing.T) {
		s := "Bye"
		a := Translate(s, "es", "en")
		final := "Adiós"
		if a != final {
			t.Errorf("Expected  %q, got %q", final, a)
		}
	})
}
