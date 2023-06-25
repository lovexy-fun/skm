package skm

import (
	"log"
	"testing"
)

func TestName(t *testing.T) {
	encrypt, err := AES256Encrypt("1234", "")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("encrypt: %s", encrypt)
	decrypt, err := AES256Decrypt(encrypt, "")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("decrypt: %s", decrypt)
}
