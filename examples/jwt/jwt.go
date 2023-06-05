package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type claims struct {
	Subject string `json:"sub"`
}

func token() (string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generate: %v", err)
	}
	signingKey := jose.SigningKey{Algorithm: jose.RS256, Key: priv}
	signer, err := jose.NewSigner(signingKey, nil)
	if err != nil {
		log.Fatalf("generate: %v", err)
	}
	token, err := jwt.Signed(signer).Claims(claims{"foo"}).CompactSerialize()
	if err != nil {
		log.Fatalf("creating jwt: %v", err)
	}

	return token, err
}

func generateJSONWebKey() (*jose.JSONWebKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}
	return &jose.JSONWebKey{
		Key:       key,
		KeyID:     "my-key-id",
		Algorithm: string(jose.RS256),
	}, nil
}

func main() {
	key, err := generateJSONWebKey()
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(key)
	if err != nil {
		panic(err)
	}

	log.Printf("%s", bytes)
}
