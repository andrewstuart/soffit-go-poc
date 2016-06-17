package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"net/http"
	"os"

	vault "github.com/hashicorp/vault/api"
)

const stuartCa = `-----BEGIN CERTIFICATE-----
MIIDLzCCAhegAwIBAgIJALsD4dhjW9hAMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNV
BAMMCVN0dWFydCBDQTAeFw0xNjAzMTMxOTM5MjBaFw0yNjAzMTExOTM5MjBaMBQx
EjAQBgNVBAMMCVN0dWFydCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBANbhNmArgIfJKSbh92o9lBxV2SoPshvdOvrTkmrHTuoDkPqWoTXUzaKaUTAH
6t7wHTxoxWShHbQaohYwkdJa74DVEw3pkChr5jOCN4XkgKNv4JUkFSKoArBdFSMN
QxtuYTgedumlwyFG7kAumE2wwiNA5t2tLNArZFJapks9iPyMbO5oCzXiWpn8/6KV
OZXOnCYJsDcMfJ7Jq0lzRVU9y/mQYF8YndK23CTGegMGFsg0i8/2nmxogJFYC+hi
Cd7+PnGp4usno+pVlHflBy25lkSx7Udq+5EMi6s7ebfezGFn2Ia0yvud8uMHDmuE
UWehJfZrJj1dsaG+8Wo+N79/gm8CAwEAAaOBgzCBgDAdBgNVHQ4EFgQUtS0ilhtv
JOcwwukbR54BryyTmfEwRAYDVR0jBD0wO4AUtS0ilhtvJOcwwukbR54BryyTmfGh
GKQWMBQxEjAQBgNVBAMMCVN0dWFydCBDQYIJALsD4dhjW9hAMAwGA1UdEwQFMAMB
Af8wCwYDVR0PBAQDAgEGMA0GCSqGSIb3DQEBCwUAA4IBAQCaJYZ/4Y21Zl5c2POh
mJgtWgay6SMxrjSo9CmcU3NMCWPlLOCSjJoXjlV99pebfiSPM2q762rLghOSEluN
7v0H6MOnfQFrkwfcbctPZtgBDVq0uLnNdJWKoFVf/puLsalCATZyFaoayxZgorjh
+csM51pgD0SHjSdmHdmweBCu5xOLJBq5o7ek5OxANf2msgLm97wFkrnC4/9lQxiZ
LvOT5Q8e91FUqoVFd6d/ZQqzAcg49Neug4uvVQ0m9uFul1jsYXtO0zSHeshyHDCw
YTf6BmFDa52efAb3F88mNwCpVTU1GGnWDv4THfqUbnN4nm3vvKMUVAnhk1SwgV4W
AvQT
-----END CERTIFICATE-----
`

var (
	token = os.Getenv("VAULT_TOKEN")
	cname = os.Getenv("SERVICE_NAME")
)

func getSignedCert(k *rsa.PrivateKey) (*tls.Certificate, error) {
	if cname == "" {
		cname = "test.svc.cluster.local"
	}

	cp := x509.NewCertPool()
	cp.AppendCertsFromPEM([]byte(stuartCa))
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: cp},
		},
	}

	cfg := vault.Config{
		Address:    "https://vault.astuart.co:8200",
		HttpClient: cli,
	}

	vCli, err := vault.NewClient(&cfg)
	if err != nil {
		return nil, err
	}

	vCli.SetToken(token)

	tpl := &x509.CertificateRequest{
		Subject:        pkix.Name{CommonName: cname},
		EmailAddresses: []string{"andrew.stuart2@gmail.com"},
	}

	csrBs, err := x509.CreateCertificateRequest(rand.Reader, tpl, k)
	if err != nil {
		return nil, err
	}

	pemB := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBs,
	}

	data := map[string]interface{}{
		"csr":         string(pem.EncodeToMemory(pemB)),
		"common_name": tpl.Subject.CommonName,
		"format":      "pem_bundle",
	}

	secret, err := vCli.Logical().Write("pki/sign/kube", data)
	if err != nil {
		return nil, err
	}

	pubBs := []byte(secret.Data["certificate"].(string))

	pb := &pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(k),
		Type:  "RSA PRIVATE KEY",
	}

	crt, err := tls.X509KeyPair(pubBs, pem.EncodeToMemory(pb))
	return &crt, err
}
