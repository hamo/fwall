package encrypt

import ()

type encryption interface {
	encrypt(plaintext []byte) []byte
	decrypt(ciphertext []byte) []byte
}
