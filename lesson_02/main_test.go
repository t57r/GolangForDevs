package main

import "testing"

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"n=-1", -1, -1},
		{"n=0", 0, 0},
		{"n=1", 1, 1},
		{"n=2", 2, 1},
		{"n=5", 5, 5},
		{"n=10", 10, 55},
		{"n=20", 20, 6765},
	}

	for _, tt := range tests {
		t.Run("Iterative/"+tt.name, func(t *testing.T) {
			got := FibonacciIterative(tt.n)
			if got != tt.want {
				t.Fatalf("FibonacciIterative(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})

		t.Run("Recursive/"+tt.name, func(t *testing.T) {
			got := FibonacciRecursive(tt.n)
			if got != tt.want {
				t.Fatalf("FibonacciRecursive(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"n=-10", -10, false},
		{"n=0", 0, false},
		{"n=1", 1, false},
		{"n=2", 2, true},
		{"n=3", 3, true},
		{"small composite", 4, false},
		{"odd composite", 9, false},
		{"non-trivial composite", 15, false},
		{"prime 29", 29, true},
		{"large prime just under 1000", 997, true},
		{"composite with large factors", 899, false}, // 29 * 31
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPrime(tt.n)
			if got != tt.want {
				t.Fatalf("IsPrime(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestIsBinaryPalindrome(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"n=0", 0, true},              // "0"
		{"n=1", 1, true},              // "1"
		{"n=2", 2, false},             // "10"
		{"n=3", 3, true},              // "11"
		{"n=5", 5, true},              // "101"
		{"n=6", 6, false},             // "110"
		{"n=7", 7, true},              // "111"
		{"n=8", 8, false},             // "1000"
		{"n=9", 9, true},              // "1001"
		{"palindrome 585", 585, true}, // 1001001001
		{"non-pal 10", 10, false},     // 1010
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBinaryPalindrome(tt.n)
			if got != tt.want {
				t.Fatalf("IsBinaryPalindrome(%d) = %v, want %v (bin check)", tt.n, got, tt.want)
			}
		})
	}
}

func TestValidParentheses(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty", "", true},
		{"single pair round", "()", true},
		{"single pair square", "[]", true},
		{"single pair curly", "{}", true},
		{"simple combo", "[]{}()", true},
		{"nested mixed", "{[()]}", true},
		{"deep nested", "({[]})", true},
		{"wrong order", "[{]}", false},
		{"extra closing", "())", false},
		{"only opening", "(((", false},
		{"only closing", ")))", false},
		{"prefix valid then invalid", "()]", false},
		{"valid prefix but unfinished", "()[", false},
		{"interleaved wrong", "([)]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidParentheses(tt.s)
			if got != tt.want {
				t.Fatalf("ValidParentheses(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIncrement(t *testing.T) {
	tests := []struct {
		name string
		num  string
		want int
	}{
		{"simple 101", "101", 6},     // 5 + 1
		{"zero", "0", 1},             // 0 + 1
		{"one", "1", 2},              // 1 + 1
		{"carry all bits", "111", 8}, // 7 + 1
		{"no carry", "1000", 9},      // 8 + 1
		{"alternating", "1010", 11},  // 10 + 1
		{"leading zeros", "0011", 4}, // 3 + 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Increment(tt.num)
			if got != tt.want {
				t.Fatalf("Increment(%q) = %d, want %d", tt.num, got, tt.want)
			}
		})
	}
}
