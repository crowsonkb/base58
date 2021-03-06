// The base58 command-line tool encodes and decodes base58 data (as used in
// Bitcoin addresses). base58 can also compute base58-encoded message digests
// using md5, sha1, sha256, or sha512.
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"github.com/crowsonkb/base58"
	"io/ioutil"
	"os"
	"strings"
)

var usageMessage = `
base58 encodes and decodes base58 data (as used in Bitcoin addresses). With no
options, base58 reads raw data from stdin and writes encoded data to stdout.
`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage[1:])
	fmt.Fprintln(os.Stderr, "Usage of base58:")
	flag.PrintDefaults()
}

var (
	flagDecode       bool
	flagFixedPadding bool
	flagHash         string
)

func init() {
	flag.Usage = usage

	flag.BoolVar(&flagDecode, "decode", false,
		"Read base58 data and output binary data.")
	flag.BoolVar(&flagFixedPadding, "fixed", false,
		"Use a fixed-length padding scheme instead of the Bitcoin scheme.")
	flag.StringVar(&flagHash, "hash", "",
		"Hash the input. Valid algorithms: md5, sha1, sha256, sha512")
}

func exit(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flagDecode && flagHash != "" {
		exit("invalid combination of options")
	}

	encoding := base58.Bitcoin
	if flagFixedPadding {
		encoding = base58.Fixed
	}

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		exit(err)
	}

	if flagDecode {
		result, err := encoding.Decode(strings.TrimSpace(string(buf)))
		if err != nil {
			exit(err)
		}
		os.Stdout.Write(result)
		return
	}

	switch flagHash {
	case "":
	case "md5":
		arr := md5.Sum(buf)
		buf = arr[:]
	case "sha1":
		arr := sha1.Sum(buf)
		buf = arr[:]
	case "sha256":
		arr := sha256.Sum256(buf)
		buf = arr[:]
	case "sha512":
		arr := sha512.Sum512(buf)
		buf = arr[:]
	default:
		exit("invalid hash algorithm")
	}
	result := encoding.Encode(buf)
	fmt.Println(result)
}
