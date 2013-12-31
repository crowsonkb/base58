// base58 encodes and decodes base58 data (as used in Bitcoin addresses).
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
options, base58 reads raw data from stdin and writes encoded data to stdout.`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage[1:])
	fmt.Fprintln(os.Stderr, "Usage of base58:")
	flag.PrintDefaults()
}

var flagDecode bool
var flagHash string

func init() {
	flag.Usage = usage

	flag.BoolVar(&flagDecode, "decode", false,
		"Read base58 data and output binary data.")
	flag.StringVar(&flagHash, "hash", "",
		"Hash the input. Valid algorithms: md5, sha1, sha256, sha512")
}

func main() {
	flag.Parse()

	if flagDecode && flagHash != "" {
		fmt.Fprintln(os.Stderr, "invalid combination of options")
		os.Exit(1)
	}

	if flagDecode {
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result, err := base58.Decode(strings.TrimSpace(string(buf)))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Stdout.Write(result)
		return
	}

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, "invalid hash algorithm")
		os.Exit(1)
	}
	result := base58.Encode(buf)
	fmt.Println(result)
}