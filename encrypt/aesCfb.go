package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
)

type aesStream struct {
	enc cipher.Stream
	dec cipher.Stream
}

var aesCfb256 CryptoInfo = CryptoInfo{
	IVlen:  16,
	GenIV:  aesGenRandIv(16),
	GenKey: aesGenKey(16),

	NewCrypto: NewAesStream,
}

func init() {
	CryptoTable["aes-256-cfb"] = &aesCfb256
}

func aesGenRandIv(n int) func() []byte {
	return func() []byte {
		bytes := make([]byte, n)
		rand.Read(bytes)
		return bytes
	}
}

func aesGenKey(n int) func(string) []byte {
	return func(passwd string) []byte {
		s := sha512.Sum512([]byte(passwd))
		return s[:n]
	}
}

func NewAesStream(iv, key []byte) Encryption {
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
