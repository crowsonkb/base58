// Package base58 implements base58 encoding as used in Bitcoin.
package base58

import (
	"fmt"
	"math"
	"math/big"
)

// The 58-character encoding alphabet.
const Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// The radix of the base58 encoding system.
const Radix = 58

// Bits of entropy per digit, in base 58.
var BitsPerDigit = math.Log2(Radix)

var invAlphabet map[rune]*big.Int
var radixBig = big.NewInt(Radix)

func init() {
	invAlphabet = make(map[rune]*big.Int)
	for index, value := range Alphabet {
		invAlphabet[value] = big.NewInt(int64(index))
	}
}

type CorruptInputError int64

func (err CorruptInputError) Error() string {
	return fmt.Sprintf("illegal base58 data at input byte %d", err)
}

// DecodeInt returns the big.Int represented by the base58 string s.
func DecodeInt(s string) (*big.Int, error) {
	n := new(big.Int)
	for index, digit := range s {
		n.Mul(n, radixBig)
		value, ok := invAlphabet[digit]
		if !ok {
			return nil, CorruptInputError(index)
		}
		n.Add(n, value)
	}
	return n, nil
}

// Decode returns the bytes represented by the base58 string s.
func Decode(s string) ([]byte, error) {
	var zeros int
	for i := 0; i < len(s) && s[i] == Alphabet[0]; i++ {
		zeros++
	}
	n, err := DecodeInt(s)
	if err != nil {
		return nil, err
	}
	buf := n.Bytes()
	bufPadded := make([]byte, len(buf)+zeros)
	copy(bufPadded[zeros:], buf)
	return bufPadded, nil
}

// EncodeInt encodes the big.Int n using base58.
func EncodeInt(n *big.Int) string {
	buf := make([]byte, 0)
	remainder := new(big.Int)
	for n.Sign() == 1 {
		n.DivMod(n, radixBig, remainder)
		buf = append(buf, Alphabet[remainder.Int64()])
	}
	bufReverse := make([]byte, len(buf))
	for index, value := range buf {
		bufReverse[len(buf)-index-1] = value
	}
	return string(bufReverse)
}

// Encode encodes src using base58.
func Encode(src []byte) string {
	var zeros int
	for i := 0; i < len(src) && src[i] == 0; i++ {
		zeros++
	}
	n := new(big.Int)
	n.SetBytes(src[zeros:])
	buf := EncodeInt(n)
	bufPadded := make([]byte, len(buf)+zeros)
	for i := 0; i < zeros; i++ {
		bufPadded[i] = Alphabet[0]
	}
	copy(bufPadded[zeros:], buf)
	return string(bufPadded)
}

// MaxEncodedLen returns the maximum length of an encoding of n source bytes.
func MaxEncodedLen(n int) int {
	return int(math.Ceil(float64(n) * 8 / BitsPerDigit))
}
