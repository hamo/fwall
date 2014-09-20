package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type aesStream struct {
	enc cipher.Stream
	dec cipher.Stream
}

/* a bit complex && redundancy due to I copy it form a random-string generator */
func GenRandIv(n int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return bytes
}

func NewAesStream(iv, key []byte) *aesStream {
	// iv = GenRandIv(32)
	block, _ := aes.NewCipher(key)
	encryptStream := cipher.NewCFBEncrypter(block, iv)
	decryptStream := cipher.NewCFBDecrypter(block, iv)
	return &aesStream{
		enc: encryptStream,
		dec: decryptStream,
	}
}

func (this *aesStream) Encrypt(plaintext []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	this.enc.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

func (this *aesStream) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	this.dec.XORKeyStream(plaintext, ciphertext)
	return plaintext
}
