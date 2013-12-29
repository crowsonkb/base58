// Package base58 implements base58 encoding as used in Bitcoin addresses.
package base58

import (
	"fmt"
	"math"
	"math/big"
)

// The 58-character encoding alphabet.
const Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// The radix of the base58 encoding system.
const Radix = len(Alphabet)

// Bits of entropy per base 58 digit.
var BitsPerDigit = math.Log2(float64(Radix))

var invAlphabet map[rune]*big.Int
var radixBig = big.NewInt(int64(Radix))

func init() {
	invAlphabet = make(map[rune]*big.Int, Radix)
	for index, value := range Alphabet {
		invAlphabet[value] = big.NewInt(int64(index))
	}
}

type CorruptInputError int

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
	return append(make([]byte, zeros), n.Bytes()...), nil
}

// EncodeInt encodes the big.Int n using base58.
func EncodeInt(n *big.Int) string {
	n = new(big.Int).Set(n)
	buf := make([]byte, 0, MaxEncodedLen(n.BitLen()))
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
	n := new(big.Int).SetBytes(src[zeros:])
	buf := append(make([]byte, zeros), EncodeInt(n)...)
	for i := 0; i < zeros; i++ {
		buf[i] = Alphabet[0]
	}
	return string(buf)
}

// MaxEncodedLen returns the maximum length of an encoding of n source bits.
func MaxEncodedLen(n int) int {
	return int(math.Ceil(float64(n) / BitsPerDigit))
}
