// Package base58 implements base58 encoding as used in Bitcoin addresses.
package base58

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

// The 58-character encoding alphabet.
const Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// The radix of the base58 encoding system.
const Radix = len(Alphabet)

// Bits of entropy per base 58 digit.
var BitsPerDigit = math.Log2(float64(Radix))

var invAlphabet map[byte]*big.Int
var radixBig = big.NewInt(int64(Radix))

// An Encoding defines a base58 encoding/decoding scheme. The 'Bitcoin' and
// 'Fixed' encodings differ only in how they pad the encoded string with initial
// '1's. The Bitcoin encoding is compatible with the Bitcoin address format,
// which produces variable length output strings given the the same number of
// input bytes. The Fixed encoding always pads with the maximum number of
// initial '1's such that the output string is always the same length given the
// same number of input bytes.
type Encoding struct {
	Decode func(string) ([]byte, error)
	Encode func([]byte) string
}

// The Bitcoin encoding is compatible with the Bitcoin address format, which
// produces variable length output strings given the the same number of input
// bytes.
var Bitcoin = Encoding{decodeBitcoin, encodeBitcoin}

// The Fixed encoding always pads with the maximum number of initial '1's such
//that the output string is always the same length given the same number of
// input bytes.
var Fixed = Encoding{decodeFixed, encodeFixed}

func init() {
	invAlphabet = make(map[byte]*big.Int, Radix)
	for index, value := range []byte(Alphabet) {
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
	for index, digit := range []byte(s) {
		n.Mul(n, radixBig)
		value, ok := invAlphabet[digit]
		if !ok {
			return nil, CorruptInputError(index)
		}
		n.Add(n, value)
	}
	return n, nil
}

func decodeBitcoin(s string) ([]byte, error) {
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

func decodeFixed(s string) ([]byte, error) {
	n, err := DecodeInt(s)
	if err != nil {
		return nil, err
	}
	buf := n.Bytes()
	zeros := DecodedLen(len(s)) - len(buf)
	if zeros <= 0 {
		return buf, nil
	}
	return append(make([]byte, zeros), buf...), nil
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

func encodeBitcoin(src []byte) string {
	var zeros int
	for i := 0; i < len(src) && src[i] == 0; i++ {
		zeros++
	}
	n := new(big.Int).SetBytes(src[zeros:])
	return strings.Repeat(Alphabet[:1], zeros) + EncodeInt(n)
}

func encodeFixed(src []byte) string {
	n := new(big.Int).SetBytes(src)
	buf := []byte(EncodeInt(n))
	zeros := MaxEncodedLen(len(src)*8) - len(buf)
	return strings.Repeat(Alphabet[:1], zeros) + string(buf)
}

// DecodedLen returns the decoded length in bytes of an n-digit base58 string
// with no initial padding.
func DecodedLen(n int) int {
	return int(math.Floor(float64(n) * BitsPerDigit / 8))
}

// MaxEncodedLen returns the maximum length in bytes of an encoding of n source
// bits.
func MaxEncodedLen(n int) int {
	return int(math.Ceil(float64(n) / BitsPerDigit))
}
