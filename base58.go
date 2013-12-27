// Package base58 implements base58 encoding as used in Bitcoin.
package base58

import (
	"fmt"
	"math/big"
)

const Base = 58
const Dict = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var baseBig = big.NewInt(Base)
var invDict map[rune]*big.Int

func init() {
	invDict = make(map[rune]*big.Int)
	for index, value := range Dict {
		invDict[value] = big.NewInt(int64(index))
	}
}

type CorruptInputError int64

func (err CorruptInputError) Error() string {
	return fmt.Sprintf("illegal base58 data at input byte %d", err)
}

func DecodeInt(src string) (*big.Int, error) {
	n := new(big.Int)
	for index, digit := range src {
		n.Mul(n, baseBig)
		value, ok := invDict[digit]
		if !ok {
			return nil, CorruptInputError(index)
		}
		n.Add(n, value)
	}
	return n, nil
}

func EncodeInt(src *big.Int) string {
	buf := make([]byte, 0)
	remainder := new(big.Int)
	for src.Sign() == 1 {
		src.DivMod(src, baseBig, remainder)
		buf = append(buf, Dict[remainder.Int64()])
	}
	bufReverse := make([]byte, len(buf))
	for index, value := range buf {
		bufReverse[len(buf)-index-1] = value
	}
	return string(bufReverse)
}

func Decode(src string) ([]byte, error) {
	var zeros int
	for i := 0; i < len(src) && src[i] == Dict[0]; i++ {
		zeros++
	}
	n, err := DecodeInt(src)
	if err != nil {
		return nil, err
	}
	buf := n.Bytes()
	bufPadded := make([]byte, len(buf)+zeros)
	copy(bufPadded[zeros:], buf)
	return bufPadded, nil
}

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
		bufPadded[i] = Dict[0]
	}
	copy(bufPadded[zeros:], buf)
	return string(bufPadded)
}
