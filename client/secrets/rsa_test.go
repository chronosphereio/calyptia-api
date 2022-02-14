package secrets

import (
	"testing"
)

func TestRSA_EncryptWithPublicKey(t *testing.T) {
	t.Run("test encrypt/decrypt with valid public key", func(t *testing.T) {
		handler := NewRSAHandler(DefaultEncryptionBits)
		private, public, err := handler.GenerateKeyPair()
		if err != nil {
			t.Errorf("err != nil, err: %v", err)
			return
		}
		toEncrypt := []byte("testing")
		encrypted, err := handler.EncryptWithPublicKey(toEncrypt, public)
		if err != nil {
			t.Errorf("err != nil, err: %v", err)
			return
		}

		decrypted, err := handler.DecryptWithPrivateKey(encrypted, private)
		if err != nil {
			t.Errorf("err != nil, err: %v", err)
			return
		}

		if want, got := string(toEncrypt), string(decrypted); want != got {
			t.Errorf("want: %v != got: %v", want, got)
			return
		}
	})
	t.Run("test encrypt/decrypt with invalid public key", func(t *testing.T) {
		handler := NewRSAHandler(DefaultEncryptionBits)
		private, _, err := handler.GenerateKeyPair()
		if err != nil {
			t.Errorf("err != nil, err: %v", err)
			return
		}

		_, newPublicKey, err := handler.GenerateKeyPair()
		if err != nil {
			t.Errorf("err != nil, err: %v", err)
			return
		}
		toEncrypt := []byte("testing")
		encrypted, err := handler.EncryptWithPublicKey(toEncrypt, newPublicKey)
		if err != nil {
			t.Errorf("err != nil, err: %v", err)
			return
		}

		_, err = handler.DecryptWithPrivateKey(encrypted, private)
		if err == nil {
			t.Errorf("err == nil, should raise a decryption error")
			return
		}
	})
}
