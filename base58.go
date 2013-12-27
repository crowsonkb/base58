package base58

import (
	"fmt"
	"math/big"
)

const Base = 58

var baseBig = big.NewInt(Base)
var Dict = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
var invDict map[byte]*big.Int

func init() {
	invDict = make(map[byte]*big.Int)
	for index, value := range Dict {
		invDict[value] = big.NewInt(int64(index))
	}
}

type CorruptInputError int64

func (err CorruptInputError) Error() string {
	return fmt.Sprintf("illegal base58 data at input byte %d", err)
}

func DecodeInt(src []byte) (*big.Int, error) {
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

func EncodeInt(src *big.Int) []byte {
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
	return bufReverse
}

func Decode(src []byte) ([]byte, error) {
	var leadingOnes int
	for i := 0; i < len(src) && src[i] == '1'; i++ {
		leadingOnes++
	}
	n, err := DecodeInt(src[leadingOnes:])
	if err != nil {
		return nil, err
	}
	buf := n.Bytes()
	paddedBuf := make([]byte, len(buf)+leadingOnes)
	copy(paddedBuf[leadingOnes:], buf)
	return paddedBuf, nil
}

func Encode(src []byte) []byte {
	var leadingZeros int
	for i := 0; i < len(src) && src[i] == 0; i++ {
		leadingZeros++
	}
	n := new(big.Int)
	n.SetBytes(src[leadingZeros:])
	buf := EncodeInt(n)
	paddedBuf := make([]byte, len(buf)+leadingZeros)
	for i := 0; i < leadingZeros; i++ {
		paddedBuf[i] = '1'
	}
	copy(paddedBuf[leadingZeros:], buf)
	return paddedBuf
}
