package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

const (
	jIss = "soffit-go.test.astuart.co"
	jAud = "portal.astuart.co"

	kid = "soffit-signer"
)

const (
	KeySize = 1 << 12
)

var (
	signingKey *rsa.PrivateKey
	tlsCert    *tls.Certificate
)

func init() {
	k, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		log.Fatal(err)
	}
	signingKey = k

	bs, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	p := &pem.Block{
		Bytes: bs,
		Type:  "PUBLIC KEY",
	}

	log.Println("Public Key:")
	pem.Encode(os.Stdout, p)
}

func getJWT(req Payload, secret string) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256)

	t.Header["kid"] = kid

	t.Claims = map[string]interface{}{
		"iss": jIss,
		"aud": jAud,
		"sub": req.User.Username,
		"org.apereo.portal.soffitRequest": req,
		"secret": secret,
	}

	return t.SignedString(signingKey)
}
