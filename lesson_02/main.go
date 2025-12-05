package main

import (
	"fmt"
	"math"
	"strconv"
)

func FibonacciIterative(n int) int {
	if n < 2 {
		return n
	}
	x, y := 0, 1
	for range n {
		x, y = x+y, x
	}
	return x
}

func FibonacciRecursive(n int) int {
	if n < 2 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func IsPrime(n int) bool {
	if n == 2 {
		return true
	}
	if n <= 1 || n%2 == 0 {
		return false
	}
	for i := 3; i <= int(math.Sqrt(float64(n))); i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func IsBinaryPalindrome(n int) bool {
	binaryString := strconv.FormatInt(int64(n), 2)
	l, r := 0, len(binaryString)-1
	for l <= r {
		if binaryString[l] != binaryString[r] {
			return false
		}
		l++
		r--
	}
	return true
}

var openClosedParentheses = map[rune]rune{
	'(': ')',
	'[': ']',
	'{': '}',
}

func ValidParentheses(s string) bool {
	isOpenParentheses := func(ch rune) bool {
		return ch == '[' || ch == '(' || ch == '{'
	}
	isClosedParentheses := func(ch rune) bool {
		return ch == ']' || ch == ')' || ch == '}'
	}

	openParentheses := NewRuneStack()
	for _, ch := range s {
		switch {
		case isOpenParentheses(ch):
			openParentheses.Push(ch)

		case isClosedParentheses(ch):
			if openParentheses.IsEmpty() {
				// meet closed parentheses before open
				return false
			}
			open, _ := openParentheses.Pop()
			closed := openClosedParentheses[open]
			if ch != closed {
				return false
			}
		}
	}
	return openParentheses.IsEmpty() // true if no more open parentheses left
}

func Increment(num string) int {
	parsedInt, err := strconv.ParseInt(num, 2, 32)
	if err != nil {
		return 0 // error parsing binary string
	}
	return int(parsedInt) + 1
}

func main() {
	fmt.Println("FibonacciIterative(10):", FibonacciIterative(10)) // очікуємо 55
	fmt.Println("FibonacciRecursive(10):", FibonacciRecursive(10)) // очікуємо 55

	fmt.Println("IsPrime(2):", IsPrime(2))   // true
	fmt.Println("IsPrime(15):", IsPrime(15)) // false
	fmt.Println("IsPrime(29):", IsPrime(29)) // true

	fmt.Println("IsBinaryPalindrome(7):", IsBinaryPalindrome(7)) // true (111)
	fmt.Println("IsBinaryPalindrome(6):", IsBinaryPalindrome(6)) // false (110)

	fmt.Println(`ValidParentheses("[]{}()"):`, ValidParentheses("[]{}()")) // true
	fmt.Println(`ValidParentheses("[{]}"):`, ValidParentheses("[{]}"))     // false

	fmt.Println(`Increment("101") ->`, Increment("101")) // 6
}
