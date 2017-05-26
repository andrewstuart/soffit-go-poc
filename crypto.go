package main

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"fmt"
)

func PBKDF1(pass, salt []byte, count, l int) ([]byte, error) {
	if l > 20 {
		return nil, fmt.Errorf("Derived key too long")
	}

	derived := make([]byte, len(pass)+len(salt))
	copy(derived, pass)
	copy(derived[len(pass):], salt)

	for i := 0; i < count; i++ {
		dr := md5.Sum(derived)
		derived = dr[:]
	}

	return derived[:l], nil
}

func decrypt(encrypted []byte, password string) ([]byte, error) {
	if len(encrypted) < saltLen {
		return nil, fmt.Errorf("Invalid encrypted text")
	}

	salt := encrypted[:des.BlockSize]
	encrypted = encrypted[des.BlockSize:]

	key, err := PBKDF1([]byte(password), salt, 1000, des.BlockSize*2)
	if err != nil {
		return nil, err
	}

	iv := key[des.BlockSize:]
	key = key[:des.BlockSize]

	b, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bm := cipher.NewCBCDecrypter(b, iv)
	bm.CryptBlocks(encrypted, encrypted)

	last := len(encrypted) - 1
	pad := int(encrypted[last])

	return encrypted[:len(encrypted)-pad], nil
}
