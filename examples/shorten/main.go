package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/multiformats/go-multihash"
)

func main() {
	// ignores errors for simplicity.
	// don't do that at home.
	// Decode a SHA1 hash to a binary buffer
	buf, _ := hex.DecodeString("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33")

	// Create a new multihash with it.
	mHashBuf, _ := multihash.EncodeName(buf, "sha1")
	// Print the multihash as hex string
	fmt.Printf("hex: %s\n", hex.EncodeToString(mHashBuf))

	// Parse the binary multihash to a DecodedMultihash
	mHash, _ := multihash.Decode(mHashBuf)
	// Convert the sha1 value to hex string
	sha1hex := hex.EncodeToString(mHash.Digest)
	// Print all the information in the multihash
	fmt.Printf("obj: %v 0x%x %d %s\n", mHash.Name, mHash.Code, mHash.Length, sha1hex)

	minLen := 5
	url := []byte("https://www.baidu.com")
	sum, _ := multihash.Sum(url, multihash.SHA2_256, -1)
	urlHash := sum.B58String()
	proposed := strings.ToLower(urlHash[len(urlHash)-minLen:])
	fmt.Printf("%s - %s\n", urlHash, proposed)
}
