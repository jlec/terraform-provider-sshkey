// Package keygen handles the creation of new SSH key pairs.
package keygen

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/user"

	"golang.org/x/crypto/ssh"
)

// KeyType represents a type of SSH key.
type KeyType string

// Supported key types.
const (
	RSA     KeyType = "rsa"
	ED25519 KeyType = "ed25519"
	ECDSA   KeyType = "ecdsa"
)

const RsaDefaultBits = 4096

//nolint:gochecknoglobals
var (
	SSHKeyTypes        = []KeyType{RSA, ED25519, ECDSA}
	SSHKeyTypesStrings = []string{"rsa", "ed25519", "ecdsa"}
	SSHRsaBits         = []int64{1024, 2048, 4096}
)

// ErrMissingSSHKeys indicates we're missing some keys that we expected to
// have after generating. This should be an extreme edge case.
var ErrMissingSSHKeys = errors.New(
	"missing one or more keys; did something happen to them after they were generated?",
)

// UnsupportedKeyTypeError indicates an unsupported key type.
type UnsupportedKeyTypeError struct {
	Type string
}

// Error implements the error interface for ErrUnsupportedKeyType.
func (e UnsupportedKeyTypeError) Error() string {
	err := "unsupported key type"
	if e.Type != "" {
		err += ": " + e.Type
	}

	return err
}

// FilesystemError is used to signal there was a problem creating keys at the
// filesystem-level. For example, when we're unable to create a directory to
// store new SSH keys in.
type FilesystemError struct {
	Err error
}

// Error returns a human-readable string for the error. It implements the error
// interface.
func (e FilesystemError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e FilesystemError) Unwrap() error {
	return e.Err
}

// SSHKeyPairConfig holds the SSH key pair configuration.
type SSHKeyPairConfig struct {
	// Type is the type of the SSH key pair.
	Type KeyType
	// Bits - RSA bit size
	Bits uint16
	// Comment for the ssh key pair
	Comment string
	// Passphrase
	Passphrase []byte
}

// SSHKeyPair holds a pair of SSH keys and associated methods.
type SSHKeyPair struct {
	Passphrase    []byte
	Type          KeyType
	Bits          uint16
	PrivateKeyRaw crypto.PrivateKey
	Comment       string
}

func (s *SSHKeyPair) pemBlock() (*pem.Block, error) {
	key := s.PrivateKey()
	if key == nil {
		return nil, ErrMissingSSHKeys
	}

	switch s.Type {
	case RSA, ED25519, ECDSA:
		if len(s.Passphrase) > 0 {
			//nolint:wrapcheck
			return ssh.MarshalPrivateKeyWithPassphrase(key, s.Comment, nil)
		}

		//nolint:wrapcheck
		return ssh.MarshalPrivateKey(key, s.Comment)
	default:
		return nil, UnsupportedKeyTypeError{string(s.Type)}
	}
}

// generateED25519Keys creates a pair of EdD25519 keys for SSH auth.
func (s *SSHKeyPair) generateED25519Keys() error {
	// Generate keys
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	s.PrivateKeyRaw = &privateKey

	return nil
}

// generateED25519Keys creates a pair of EdD25519 keys for SSH auth.
func (s *SSHKeyPair) generateECDSAKeys(curve elliptic.Curve) error {
	// Generate keys
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	s.PrivateKeyRaw = privateKey

	return nil
}

// generateRSAKeys creates a pair for RSA keys for SSH auth.
func (s *SSHKeyPair) generateRSAKeys() error {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, int(s.Bits))
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}
	// Validate private key
	err = privateKey.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate key: %w", err)
	}

	s.PrivateKeyRaw = privateKey

	return nil
}

// PrivateKey returns the unencrypted private key.
func (s *SSHKeyPair) PrivateKey() crypto.PrivateKey {
	switch s.Type {
	case RSA, ED25519, ECDSA:
		return s.PrivateKeyRaw
	default:
		return nil
	}
}

// PrivateKeyPEM returns the unencrypted private key in OPENSSH PEM format.
func (s *SSHKeyPair) PrivateKeyPEM() []byte {
	block, err := s.pemBlock()
	if err != nil {
		return nil
	}

	return pem.EncodeToMemory(block)
}

func (s *SSHKeyPair) publicKeyRaw() crypto.PublicKey {
	var pkey crypto.PublicKey
	// Prepare public key
	switch s.Type {
	case RSA:
		key, ok := s.PrivateKeyRaw.(*rsa.PrivateKey)
		if !ok {
			return nil
		}

		pkey = key.Public()

	case ED25519:
		key, ok := s.PrivateKeyRaw.(*ed25519.PrivateKey)
		if !ok {
			return nil
		}

		pkey = key.Public()

	case ECDSA:
		key, ok := s.PrivateKeyRaw.(*ecdsa.PrivateKey)
		if !ok {
			return nil
		}

		pkey = key.Public()

	default:
		return nil
	}

	return pkey
}

// PublicKey returns the SSH public key (RFC 4253). Ready to be used in an
// OpenSSH authorized_keys file.
func (s *SSHKeyPair) PublicKey() []byte {
	pkey, err := ssh.NewPublicKey(s.publicKeyRaw())
	if err != nil {
		return nil
	}

	// serialize public key
	ak := ssh.MarshalAuthorizedKey(pkey)

	return bytes.TrimSpace(fmt.Appendf(bytes.TrimSpace(ak), " %s", s.Comment))
}

func (s *SSHKeyPair) MD5() string {
	p, _ := ssh.NewPublicKey(s.publicKeyRaw())

	return ssh.FingerprintLegacyMD5(p)
}

func (s *SSHKeyPair) SHA256() string {
	p, _ := ssh.NewPublicKey(s.publicKeyRaw())

	return ssh.FingerprintSHA256(p)
}

// New generates an SSHKeyPair, which contains a pair of SSH keys.
func New(conf *SSHKeyPairConfig) (*SSHKeyPair, error) {
	var err error

	if conf.Comment == "" {
		conf.Comment = GetSSHKeyComment()
	}

	skeypair := &SSHKeyPair{
		Type:       conf.Type,
		Passphrase: conf.Passphrase,
		Comment:    conf.Comment,
	}

	if conf.Bits == 0 && conf.Type == RSA {
		skeypair.Bits = RsaDefaultBits
	} else {
		skeypair.Bits = conf.Bits
	}

	switch conf.Type {
	case ED25519:
		err = skeypair.generateED25519Keys()
	case RSA:
		err = skeypair.generateRSAKeys()
	case ECDSA:
		err = skeypair.generateECDSAKeys(elliptic.P384())
	default:
		return nil, UnsupportedKeyTypeError{string(conf.Type)}
	}

	if err != nil {
		return nil, err
	}

	return skeypair, nil
}

// attaches a user@host suffix to a serialized public key. returns the original
// pubkey if we can't get the username or host.
func GetSSHKeyComment() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s@%s\n", usr.Username, hostname)
}
