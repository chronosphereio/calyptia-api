package secrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const (
	DefaultEncryptionBits    = 2048
	DefaultTypeRSAPublicKey  = "RSA PUBLIC KEY"
	DefaultTypeRSAPrivateKey = "RSA PRIVATE KEY"
)

type RSA struct {
	Handler
	Bits int
}

func NewRSAHandler(bits int) *RSA {
	if bits == 0 {
		bits = DefaultEncryptionBits
	}
	return &RSA{Bits: bits}
}

func (r *RSA) EncryptWithPublicKey(msg []byte, key []byte) ([]byte, error) {
	publicKey, err := BytesToPublicKey(key)
	if err != nil {
		return nil, err
	}
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func (r *RSA) DecryptWithPrivateKey(msg []byte, key []byte) ([]byte, error) {
	privateKey, err := BytesToPrivateKey(key)
	if err != nil {
		return nil, err
	}
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, msg, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func (r *RSA) GenerateKeyPair() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, r.Bits)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := PublicKeyToBytes(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return PrivateKeyToBytes(privateKey), publicKeyBytes, nil
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  DefaultTypeRSAPrivateKey,
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
}

func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  DefaultTypeRSAPublicKey,
		Bytes: pubASN1,
	})
	return pubBytes, nil
}

func BytesToPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func BytesToPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot decode bytes to public key")
	}
	return key, nil
}
