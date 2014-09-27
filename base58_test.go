package base58

import (
	"bytes"
	"math/big"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	src, want := int64((3*58*58)+(2*58)+1), "432"
	n := big.NewInt(src)
	got := EncodeInt(n)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if n.Int64() != src {
		t.Fatalf("input was altered")
	}
}

func TestEncodeIntZero(t *testing.T) {
	src, want := new(big.Int), ""
	got := EncodeInt(src)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDecodeInt(t *testing.T) {
	src, want := "432", big.NewInt((3*58*58)+(2*58)+1)
	got, err := DecodeInt(src)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDecodeIntZero(t *testing.T) {
	src, want := "", new(big.Int)
	got, err := DecodeInt(src)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDecodeIntBad(t *testing.T) {
	src, want := "43=2", CorruptInputError(2)
	_, err := DecodeInt(src)
	got, ok := err.(CorruptInputError)
	if !ok {
		t.Fatalf("Error %v is not a CorruptInputError", err)
	}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestEncodeBitcoin(t *testing.T) {
	src, want := []byte{0, 0, 0, 58}, "11121"
	got := Bitcoin.Encode(src)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestEncodeFixed(t *testing.T) {
	src, want := []byte{0, 0, 0, 58}, "111121"
	got := Fixed.Encode(src)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDecodeBitcoin(t *testing.T) {
	src, want := "11121", []byte{0, 0, 0, 58}
	got, err := Bitcoin.Decode(src)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDecodeFixed(t *testing.T) {
	src, want := "111121", []byte{0, 0, 0, 58}
	got, err := Fixed.Decode(src)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDecodeBad(t *testing.T) {
	src, want := "111=1", CorruptInputError(3)
	_, err := Bitcoin.Decode(src)
	got, ok := err.(CorruptInputError)
	if !ok {
		t.Fatalf("Error %v is not a CorruptInputError", err)
	}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestMaxEncodedLen(t *testing.T) {
	src, want := 64, 11
	got := MaxEncodedLen(src)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
