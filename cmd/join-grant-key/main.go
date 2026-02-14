package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate join grant key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("export FRACTURING_SPACE_JOIN_GRANT_PRIVATE_KEY=%s\n", base64.RawStdEncoding.EncodeToString(privateKey))
	fmt.Printf("export FRACTURING_SPACE_JOIN_GRANT_PUBLIC_KEY=%s\n", base64.RawStdEncoding.EncodeToString(publicKey))
}
