package iso

import (
	"errors"
	"fmt"
	"unicode"
)

// porting of https://github.com/ShellRechargeSolutionsEU/mobilityid/blob/master/core/src/main/scala/com/thenewmotion/mobilityid/checkDigit.scala

var (
	negP2minus15 = matrix{0, 2, 2, 1} // -p2^(-15)
	p1s          []matrix
	p2s          []matrix

	encoding map[rune]matrix
	decoding map[matrix]rune
)

func init() {
	p1 := matrix{0, 1, 1, 1}
	p2 := matrix{0, 1, 1, 2}

	for i := 0; i < 14; i++ {
		if i == 0 {
			p1s = append(p1s, p1)
			p2s = append(p2s, p2)
		} else {
			p1s = append(p1s, p1s[i-1].multiply(p1))
			p2s = append(p2s, p2s[i-1].multiply(p2))
		}
	}

	cipher := map[rune]int{
		'0': 0, '1': 16, '2': 32,
		'3': 4, '4': 20, '5': 36,
		'6': 8, '7': 24, '8': 40,
		'9': 2, 'A': 18, 'B': 34,
		'C': 6, 'D': 22, 'E': 38,
		'F': 10, 'G': 26, 'H': 42,
		'I': 1, 'J': 17, 'K': 33,
		'L': 5, 'M': 21, 'N': 37,
		'O': 9, 'P': 25, 'Q': 41,
		'R': 3, 'S': 19, 'T': 35,
		'U': 7, 'V': 23, 'W': 39,
		'X': 11, 'Y': 27, 'Z': 43,
	}

	encoding = make(map[rune]matrix)
	decoding = make(map[matrix]rune)
	for k, v := range cipher {
		m := decode(v)

		encoding[k] = m
		decoding[m] = k
	}
}

type matrix struct {
	m11, m12, m21, m22 int
}

func (m matrix) multiply(m2 matrix) matrix {
	return matrix{
		m11: m.m11*m2.m11 + m.m12*m2.m21,
		m12: m.m11*m2.m12 + m.m12*m2.m22,
		m21: m.m21*m2.m11 + m.m22*m2.m21,
		m22: m.m21*m2.m12 + m.m22*m2.m22,
	}
}

type vec struct {
	v1, v2 int
}

func (v vec) add(v2 vec) vec {
	return vec{
		v1: v.v1 + v2.v1,
		v2: v.v2 + v2.v2,
	}
}

func (v vec) multiply(m matrix) vec {
	return vec{
		v1: v.v1*m.m11 + v.v2*m.m21,
		v2: v.v1*m.m12 + v.v2*m.m22,
	}
}

func decode(x int) matrix {
	return matrix{x & 1, (x >> 1) & 1, (x >> 2) & 3, x >> 4}
}

// ComputeCheckDigit computes and returns a check digit for `code`, if possible.
//
// It returns an error if:
//
// - `code` doesn't match the expected length, or
//
// - any of the runes in `code` is not an upper case ASCII character nor a decimal digit
func ComputeCheckDigit(code string) (rune, error) {
	if len(code) != len(p1s) {
		return -1, fmt.Errorf("code must have a length of %v", len(p1s))
	}

	for _, r := range code {
		if !(unicode.IsUpper(r) || unicode.IsDigit(r)) || r > unicode.MaxASCII {
			return -1, errors.New("code must consist of uppercase ASCII letters and digits only")
		}
	}

	sumEq := func(ps []matrix, f func(matrix) vec) (vec, error) {
		result := vec{}
		for ix, p := range ps {
			ch := code[ix]
			mx, ok := encoding[rune(ch)]
			if !ok {
				return vec{}, fmt.Errorf("invalid character: %v", ch)
			}

			qr := f(mx)
			result = result.add(qr.multiply(p))
		}

		return result, nil
	}

	t1, err := sumEq(p1s, func(m matrix) vec {
		return vec{
			m.m11,
			m.m12,
		}
	})
	if err != nil {
		return -1, fmt.Errorf("unable to compute check digit: %w", err)
	}

	t2, err := sumEq(p2s, func(m matrix) vec {
		return vec{
			m.m21,
			m.m22,
		}
	})
	if err != nil {
		return -1, fmt.Errorf("unable to compute check digit: %w", err)
	}

	t2m := t2.multiply(negP2minus15)

	m15 := matrix{t1.v1 & 1, t1.v2 & 1, t2m.v1 % 3, t2m.v2 % 3}

	if result, ok := decoding[m15]; ok {
		return result, nil
	}

	return -1, fmt.Errorf("undecodable matrix: %v", m15)
}
