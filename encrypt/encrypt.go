package encrypt

import ()

type CryptoInfo struct {
	IVlen int
	GenIV func() []byte

	GenKey func(passwd string) []byte

	NewCrypto func(iv []byte, key []byte) Encryption
}

var CryptoTable map[string]*CryptoInfo = make(map[string]*CryptoInfo)

type Encryption interface {
	Encrypt(plaintext []byte) []byte
	Decrypt(ciphertext []byte) []byte
}
