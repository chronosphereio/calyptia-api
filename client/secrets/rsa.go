// Package secrets provides RSA codec to work with Calyptia Cloud pipeline secrets.
package secrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash"
)

const (
	DefaultEncryptionBits  = 2048
	blockTYpeRSAPublicKey  = "RSA PUBLIC KEY"
	blockTypeRSAPrivateKey = "RSA PRIVATE KEY"
)

type RSA struct {
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
	h := sha512.New()
	ciphertext, err := Encrypt(h, publicKey, msg)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// Encrypt encrypts the given message using RSA-OAEP with the given public key.
//
// The message is split into chunks of size k - 2*hLen - 2, where k is the size of the public key in bytes and hLen is the length of the hash function used.
// Each chunk is encrypted separately and the resulting ciphertexts are concatenated.
//
// At the end it returns the ciphertext encrypted.
func Encrypt(hash hash.Hash, publicKey *rsa.PublicKey, msg []byte) ([]byte, error) {
	k := publicKey.Size()
	limit := k - 2*hash.Size() - 2
	chunks := splitByteBySizeLimit(msg, limit)
	var cipherText []byte
	for _, chunk := range chunks {
		ciphertextChunk, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, chunk, nil)
		if err != nil {
			return nil, err
		}
		cipherText = append(cipherText, ciphertextChunk...)
	}
	return cipherText, nil
}

// splitByteBySizeLimit splits the given byte array into chunks of size limit.
// If the given byte array is smaller than the limit, it returns a slice with the given byte array.
func splitByteBySizeLimit(buf []byte, sizeLimit int) [][]byte {
	if len(buf) <= sizeLimit {
		return [][]byte{buf}
	}
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/sizeLimit+1)
	for len(buf) >= sizeLimit {
		chunk, buf = buf[:sizeLimit], buf[sizeLimit:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

// Decrypt decrypts the given ciphertext using RSA-OAEP with the given private key.
//
// The ciphertext is split into chunks of size k, where k is the size of the private key in bytes.
// Each chunk is decrypted separately and the resulting plaintexts are concatenated.
//
// At the end it returns the plaintext decrypted.
func Decrypt(hash hash.Hash, privateKey *rsa.PrivateKey, msg []byte) ([]byte, error) {
	limit := privateKey.Size()
	chunks := splitByteBySizeLimit(msg, limit)
	var plaintext []byte
	for _, chunk := range chunks {
		plaintextChunk, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, chunk, nil)
		if err != nil {
			return nil, err
		}
		plaintext = append(plaintext, plaintextChunk...)
	}
	return plaintext, nil
}

func (r *RSA) DecryptWithPrivateKey(msg []byte, key []byte) ([]byte, error) {
	privateKey, err := BytesToPrivateKey(key)
	if err != nil {
		return nil, err
	}
	h := sha512.New()
	plaintext, err := Decrypt(h, privateKey, msg)
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
			Type:  blockTypeRSAPrivateKey,
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
		Type:  blockTYpeRSAPublicKey,
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
