package pkix

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"math/big"
)

const (
	rsaPrivateKeyPEMBlockType = "RSA PRIVATE KEY"
	// key bits for RSA generation
	rsaBits = 1024
)

// CreateRSAKey creates a new Key using RSA algorithm
func CreateRSAKey() (*Key, error) {
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return nil, err
	}

	return NewKey(&priv.PublicKey, priv), nil
}

type Key struct {
	Public  crypto.PublicKey
	Private crypto.PrivateKey
	// TODO(yichengq): add pemEncryptedBlock *pem.Block
}

func NewKey(pub crypto.PublicKey, priv crypto.PrivateKey) *Key {
	return &Key{Public: pub, Private: priv}
}

// NewKeyFromPrivateKeyPEM inits Key from PEM-format rsa private key bytes
func NewKeyFromPrivateKeyPEM(data []byte) (*Key, error) {
	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, errors.New("cannot find the next PEM formatted block")
	}
	if pemBlock.Type != rsaPrivateKeyPEMBlockType || len(pemBlock.Headers) != 0 {
		return nil, errors.New("unmatched type or headers")
	}

	priv, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return NewKey(&priv.PublicKey, priv), nil
}

// NewKeyFromEncryptedPrivateKeyPEM inits Key from encrypted PEM-format rsa private key bytes
func NewKeyFromEncryptedPrivateKeyPEM(data []byte, password []byte) (*Key, error) {
	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, errors.New("cannot find the next PEM formatted block")
	}
	if pemBlock.Type != rsaPrivateKeyPEMBlockType {
		return nil, errors.New("unmatched type or headers")
	}

	b, err := x509.DecryptPEMBlock(pemBlock, password)
	if err != nil {
		return nil, err
	}

	priv, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}

	return NewKey(&priv.PublicKey, priv), nil
}

// ExportPrivate exports PEM-format private key
func (k *Key) ExportPrivate() ([]byte, error) {
	var privPEMBlock *pem.Block
	switch priv := k.Private.(type) {
	case *rsa.PrivateKey:
		privBytes := x509.MarshalPKCS1PrivateKey(priv)
		privPEMBlock = &pem.Block{
			Type:  rsaPrivateKeyPEMBlockType,
			Bytes: privBytes,
		}
	default:
		return nil, errors.New("only RSA private key is supported")
	}

	buf := new(bytes.Buffer)
	if err := pem.Encode(buf, privPEMBlock); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExportEncryptedPrivate exports encrypted PEM-format private key
func (k *Key) ExportEncryptedPrivate(password []byte) ([]byte, error) {
	var privBytes []byte
	switch priv := k.Private.(type) {
	case *rsa.PrivateKey:
		privBytes = x509.MarshalPKCS1PrivateKey(priv)
	default:
		return nil, errors.New("only RSA private key is supported")
	}

	privPEMBlock, err := x509.EncryptPEMBlock(rand.Reader, rsaPrivateKeyPEMBlockType, privBytes, password, x509.PEMCipher3DES)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := pem.Encode(buf, privPEMBlock); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// GenerateSubjectKeyId generates SubjectKeyId used in Certificate
// Id is 160-bit SHA-1 hash of the value of the BIT STRING subjectPublicKey
func GenerateSubjectKeyId(pub crypto.PublicKey) ([]byte, error) {
	var pubBytes []byte
	var err error
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		pubBytes, err = asn1.Marshal(rsaPublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only RSA public key is supported")
	}

	hash := sha1.Sum(pubBytes)

	return hash[:], nil
}