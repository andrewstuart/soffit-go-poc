package main

import (
	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
	"github.com/dgrijalva/jwt-go"
)

const (
	jIss = "soffit-go.test.astuart.co"
	jAud = "portal.astuart.co"

	kid = "soffit-signer"
)

var (
	signingKey = []byte("foobarbaz")
)

func getJWT(req soffit.Request, secret string) (string, error) {
	t := jwt.New(jwt.SigningMethodHS384)

	t.Header["kid"] = kid

	t.Claims = map[string]interface{}{
		"iss": jIss,
		"aud": jAud,
		"sub": req.UserName,
		"org.apereo.portal.soffitRequest": req,
		"secret": secret,
	}

	return t.SignedString(signingKey)
}
