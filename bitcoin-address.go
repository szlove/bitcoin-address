package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	// 0 - Having a private ECDSA key
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// 1 - Take the corresponding public key generated with it (33 bytes, 1 byte 0x02 (y-coord is even), and 32 bytes corresponding to X coordinate)
	publicKey := &privateKey.PublicKey
	comb1 := make([]byte, 65)
	comb1[0] = 0x02
	copy(comb1[1:33], publicKey.X.Bytes())
	copy(comb1[33:65], publicKey.Y.Bytes())

	// 2 - Perform SHA-256 hashing on the public key
	hash2 := sha256.New()
	hash2.Write(comb1)
	digest2 := hash2.Sum(nil)

	// 3 - Perform RIPEMD-160 hashing on the result of SHA-256
	hash3 := ripemd160.New()
	hash3.Write(digest2)
	digest3 := hash3.Sum(nil) // returning [20]byte

	// 4 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	v4 := make([]byte, 21)
	v4[0] = 0x00
	copy(v4[1:], digest3)

	// 5 - Perform SHA-256 hash on the extended RIPEMD-160 result
	hash5 := sha256.New()
	hash5.Write(v4)
	digest5 := hash5.Sum(nil)

	// 6 - Perform SHA-256 hash on the result of the previous SHA-256 hash
	hash6 := sha256.New()
	hash6.Write(digest5)
	digest6 := hash6.Sum(nil)

	// 7 - Take the first 4 bytes of the second SHA-256 hash. This is the address checksum
	checksum := digest6[:4]

	// 8 - Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
	comb8 := make([]byte, 25)
	copy(comb8[:21], v4)
	copy(comb8[21:], checksum)

	// 9 - Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
	address := base58.Encode(comb8)

	// print. ex) 1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs
	fmt.Println("Address:", address)
}
