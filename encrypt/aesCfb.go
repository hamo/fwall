package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func test() {
	const key = "1234567890123456"
	iv := []byte(key)[:aes.BlockSize]
	var msg = "message"

	encrypted := make([]byte, len(msg))
	block, _ := aes.NewCipher([]byte(key))
	encryptStream := cipher.NewCFBEncrypter(block, iv)
	encryptStream.XORKeyStream(encrypted, []byte(msg))
	fmt.Printf("Encrypting %v %s -> %v\n", []byte(msg), msg, encrypted)
}

type aesStream struct {
	enc cipher.Stream
	dec cipher.Stream
}

/* a bit complex && redundancy due to I copy it form a random-string generator */
func genRandIv(n int) []byte {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return bytes
}

func NewAesStream (iv, key []byte) (*aesStream){
	iv = genRandIv(32)
	block, _ := aes.NewCipher([]byte(key))
	encryptStream := cipher.NewCFBEncrypter(block, iv)
	decryptStream := cipher.NewCFBDecrypter(block, iv)
	return &aesStream{
		enc: encryptStream,
		dec: decryptStream,
	}
}

func (this *aesStream) encrpt (plaintext []byte) ([]byte) {
	ciphertext := make([]byte, len(plaintext))
	this.enc.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

func (this *aesStream) decrpt (ciphertext []byte) ([]byte) {
	plaintext := make([]byte, len(ciphertext))
	this.dec.XORKeyStream(plaintext, ciphertext)
	return plaintext
}
