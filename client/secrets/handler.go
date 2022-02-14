package secrets

type Handler interface {
	GenerateKeyPair() ([]byte, []byte, error)
	EncryptWithPublicKey(msg []byte, key []byte) ([]byte, error)
	DecryptWithPrivateKey(msg []byte, key []byte) ([]byte, error)
}
