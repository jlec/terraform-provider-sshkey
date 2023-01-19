package keygen_test

import (
	"testing"

	"github.com/jlec/terraform-provider-sshkey/internal/keygen"
)

func TestNewSSHKeyPair(t *testing.T) {
	c := keygen.SSHKeyPairConfig{
		Passphrase: []byte(""),
		Type:       keygen.RSA,
	}

	if _, err := keygen.New(&c); err != nil {
		t.Errorf("error creating SSH key pair: %v", err)
	}
}

func TestGenerateEd25519Keys(t *testing.T) {
	c := keygen.SSHKeyPairConfig{
		Type: keygen.ED25519,
	}
	k, err := keygen.New(&c)

	t.Run("test generate SSH keys", func(t *testing.T) {
		if err != nil {
			t.Errorf("error creating SSH key pair: %v", err)
		}

		// TODO: is there a good way to validate these? Lengths seem to vary a bit,
		// so far now we're just asserting that the keys indeed exist.
		if len(k.PrivateKeyPEM()) == 0 {
			t.Error("error creating SSH private key PEM; key is 0 bytes")
		}
		if len(k.PublicKey()) == 0 {
			t.Error("error creating SSH public key; key is 0 bytes")
		}
	})
}

func TestGenerateECDSAKeys(t *testing.T) {
	// Create temp directory for keys
	c := keygen.SSHKeyPairConfig{
		Type: keygen.ECDSA,
	}
	k, _ := keygen.New(&c)

	t.Run("test generate SSH keys", func(t *testing.T) {
		// TODO: is there a good way to validate these? Lengths seem to vary a bit,
		// so far now we're just asserting that the keys indeed exist.
		if len(k.PrivateKeyPEM()) == 0 {
			t.Error("error creating SSH private key PEM; key is 0 bytes")
		}
		if len(k.PublicKey()) == 0 {
			t.Error("error creating SSH public key; key is 0 bytes")
		}
	})
}

func TestGeneratePublicKeyWithEmptyDir(t *testing.T) {
	for _, keyType := range keygen.SshKeyTypes {
		func(t *testing.T) {
			c := keygen.SSHKeyPairConfig{Passphrase: nil, Type: keyType}
			k, err := keygen.New(&c)

			// FIXME: Do some testing
			_ = k

			if err != nil {
				t.Fatalf("error creating SSH key pair: %v", err)
			}
		}(t)
	}
}

func TestGenerateKeyWithPassphrase(t *testing.T) {
	for _, keyType := range keygen.SshKeyTypes {
		ph := "testpass"

		func(t *testing.T) {
			c := keygen.SSHKeyPairConfig{
				Passphrase: []byte(ph),
				Type:       keyType,
			}

			_, err := keygen.New(&c)
			if err != nil {
				t.Fatalf("error creating SSH key pair: %v", err)
			}

			c = keygen.SSHKeyPairConfig{Passphrase: []byte(ph), Type: keyType}

			k, err := keygen.New(&c)
			if err != nil {
				t.Fatalf("error reading SSH key pair: %v", err)
			}

			// FIXME: Do some testing
			_ = k
		}(t)
	}
}

func TestReadingKeyWithPassphrase(t *testing.T) {
	for _, keyType := range keygen.SshKeyTypes {
		c := keygen.SSHKeyPairConfig{Passphrase: []byte("test"), Type: keyType}
		if _, err := keygen.New(&c); err != nil {
			t.Fatalf("error reading SSH key pair: %v", err)
		}
	}
}
