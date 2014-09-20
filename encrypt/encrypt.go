package encrypt

import ()

type Encryption interface {
	Encrypt(plaintext []byte) []byte
	Decrypt(ciphertext []byte) []byte
}
