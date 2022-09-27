package din

import (
	"math"
	"strings"
)

// porting of https://github.com/ShellRechargeSolutionsEU/mobilityid/blob/master/core/src/main/scala/com/thenewmotion/mobilityid/checkDigit.scala

var numericValues map[rune]int

func init() {
	numericValues = make(map[rune]int)
	for i, x := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		numericValues[x] = i
	}
}

func mult(value, coeff int) int {
	return value * int(math.Pow(2, float64(coeff)))
}

func ComputeCheckDigit(code string) rune {
	var sum, mod, coeff int
	var lookupResults []int

	for _, x := range strings.ToUpper(code) {
		lookupResults = append(lookupResults, numericValues[x])
	}

	for _, x := range lookupResults {
		if x < 10 {
			sum += mult(x, coeff)
			coeff++
		} else {
			sum += mult(x/10, coeff) + mult(x%10, coeff+1)
			coeff += 2
		}
	}

	mod = sum % 11

	if mod >= 10 {
		return 'X'
	}

	return '0' + rune(mod)
}
