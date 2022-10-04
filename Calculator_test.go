package main

import (
	"testing"
)

func TestCalculation(t *testing.T) {
	t.Run("test simple calculation", func(t *testing.T) {
		e := "1+1"
		result, _ := calculation(e)
		if result != 2 {
			t.Errorf("expected 2, got %v", result)
		}
	})
	t.Run("test simple calculation", func(t *testing.T) {
		e := "1+2-3*4+3"
		result, _ := calculation(e)
		if result != -6 {
			t.Errorf("expected 6, got %v", result)
		}
	})
	t.Run("test simple calculation", func(t *testing.T) {
		e := "1+2+3*4+3"
		result, _ := calculation(e)
		if result != 18 {
			t.Errorf("expected 18, got %v", result)
		}
	})
	t.Run("test medium calculation", func(t *testing.T) {
		e := "4*5 + (5-4%3) / (2-1*(3/1))"
		result, _ := calculation(e)
		println(result)
		if result != 16 {
			t.Errorf("expected 16, got %v", result)
		}
	})
	t.Run("test complex calculation", func(t *testing.T) {
		e := "4*5 * (5-4%3) / (2-1*(3/1)) + 1 - 2 * 3 / 4"
		result, _ := calculation(e)
		if result != -80.5 {
			t.Errorf("expected -80.5, got %v", result)
		}
	})
	t.Run("super duper complex calculation with lot of parenthesis", func(t *testing.T) {
		e := "4*5 * (5-4%3) / ((2-1*(3/1)) + 1 - 2) * (3 / 4 + (1 + 2)) * (3 + 4) / (5 + 6)"
		result, _ := calculation(e)
		if result != -95.45 {
			t.Errorf("expected  -95.45, got %v", result)
		}
	})
	t.Run("test with division of minus , multiplication of minus", func(t *testing.T) {
		e := "4*-5 * (5-4%3) / (2-1*(3/1)) + 1 - 2 * 3 / 4"
		result, _ := calculation(e)
		if result != 79.5 {
			t.Errorf("expected -80.5, got %v", result)
		}
	})
	t.Run("Calculate Big numbers", func(t *testing.T) {
		e := "4 * 5 /5 (4 *4*5/9) + 1 - 2 * 3 / 4"
		result, _ := calculation(e)
		if result != 35.05 {
			t.Errorf("expected 35.05, got %v", result)
		}
	})
	t.Run("parenthesis number error", func(t *testing.T) {
		e := "4 * (3+4"
		_, err := calculation(e)
		//get error message
		if err.Error() != "no closing bracket" {
			t.Errorf("expected error 'no closing bracket', got %s", err.Error())
		}
	})
	t.Run("math with powers", func(t *testing.T) {
		e := "4 * 5 ^ 2*3 + 3^3^2"
		result, _ := calculation(e)
		if result != 19983 {
			t.Errorf("expected 19983, got %v", result)
		}
	})
}
