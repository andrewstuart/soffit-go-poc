package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
	"github.com/dgrijalva/jwt-go"
)

const (
	jIss = "soffit-go.test.astuart.co"
	jAud = "portal.astuart.co"

	kid = "soffit-signer"
)

var signingKey *rsa.PrivateKey

func init() {
	//signingKey = []byte("foobarbaz")

	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	bs, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Public Key:")
	fmt.Println("-----BEGIN PUBLIC KEY-----  ")
	fmt.Println(base64.StdEncoding.EncodeToString(bs))
	fmt.Println("-----END PUBLIC KEY-----  ")

	signingKey = k
}

func getJWT(req soffit.Request, secret string) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256)

	t.Header["kid"] = kid

	t.Claims = map[string]interface{}{
		"iss": jIss,
		"aud": jAud,
		"sub": req.UserDetails.UserName,
		"org.apereo.portal.soffitRequest": req,
		"secret": secret,
	}

	return t.SignedString(signingKey)
}
