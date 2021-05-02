package hasher

import (
	"crypto/sha512"
	"encoding/base64"
	"log"
)

/*
Encode function returns an encoded string
*/
func Encode(str string) (string, error) {
	h := sha512.Sum512([]byte(str))
	value := base64.StdEncoding.EncodeToString(h[:])
	log.Println(value)
	return value, nil
}
