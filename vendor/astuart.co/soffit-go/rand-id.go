package soffit

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

var (
	random = rand.New(rand.NewSource(time.Now().Unix()))
)

// RandHTMLID generates a random valid HTML ID to be used for namespacing
func RandHTMLID() (string, error) {
	random.Seed(time.Now().UnixNano())
	bs := make([]byte, 20)
	_, err := random.Read(bs)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(bs)

	// Replace invalid id strings
	str = strings.Replace(str, "=", "", -1)
	str = strings.Replace(str, "/", "", -1)
	str = strings.Replace(str, "\\", "", -1)

	return str, nil
}
