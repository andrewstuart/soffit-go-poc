package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
	"github.com/dgrijalva/jwt-go"
)

const (
	jIss = "soffit-go.test.astuart.co"
	jAud = "portal.astuart.co"

	kid = "soffit-signer"
)

var (
	signingKey *rsa.PrivateKey
	tlsCert    *tls.Certificate
)

func init() {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	bs, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := getSignedCert(k)
	if err != nil {
		log.Fatal(err)
	}

	tlsCert = cert

	p := &pem.Block{
		Bytes: bs,
		Type:  "PUBLIC KEY",
	}

	log.Println("Public Key:")
	pem.Encode(os.Stdout, p)

	signingKey = k
}

func getJWT(req soffit.Payload, secret string) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256)

	t.Header["kid"] = kid

	t.Claims = map[string]interface{}{
		"iss": jIss,
		"aud": jAud,
		"sub": req.Request.User.Username,
		"org.apereo.portal.soffitRequest": req,
		"secret": secret,
	}

	return t.SignedString(signingKey)
}
