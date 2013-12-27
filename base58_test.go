package base58

import (
	"bytes"
	"math/big"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	in, out := big.NewInt(3*58*58+2*58+1), "432"
	result := string(EncodeInt(in))
	if result != out {
		t.Errorf("result = %v, want %v", result, out)
	}
}

func TestEncodeIntZero(t *testing.T) {
	in, out := new(big.Int), ""
	result := EncodeInt(in)
	if result != out {
		t.Errorf("result = %v, want %v", result, out)
	}
}

func TestDecodeInt(t *testing.T) {
	in, out := "432", big.NewInt(3*58*58+2*58+1)
	result, err := DecodeInt(in)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if result.Cmp(out) != 0 {
		t.Fatalf("result = %v, want %v", result, out)
	}
}

func TestDecodeIntZero(t *testing.T) {
	in, out := "", new(big.Int)
	result, err := DecodeInt(in)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if result.Cmp(out) != 0 {
		t.Fatalf("result = %v, want %v", result, out)
	}
}

func TestDecodeIntBad(t *testing.T) {
	in, out := "43=2", 2
	_, err := DecodeInt(in)
	result, ok := err.(CorruptInputError)
	if !ok {
		t.Fatalf("Error %v is not a CorruptInputError", err)
	}
	if result != 2 {
		t.Fatalf("result = %v, want %v", result, out)
	}
}

func TestEncode(t *testing.T) {
	in, out := []byte{0, 0, 0, 58}, "11121"
	result := string(Encode(in))
	if result != out {
		t.Fatalf("result = %v, want %v", result, out)
	}
}

func TestDecode(t *testing.T) {
	in, out := "11121", []byte{0, 0, 0, 58}
	result, err := Decode(in)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !bytes.Equal(result, out) {
		t.Fatalf("result = %v, want %v", result, out)
	}
}

func TestDecodeBad(t *testing.T) {
	in, out := "111=1", 3
	_, err := Decode(in)
	result, ok := err.(CorruptInputError)
	if !ok {
		t.Fatalf("Error %v is not a CorruptInputError", err)
	}
	if result != 3 {
		t.Fatalf("result = %v, want %v", result, out)
	}
}
