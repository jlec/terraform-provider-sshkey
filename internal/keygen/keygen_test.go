package keygen_test

import (
	"testing"

	"github.com/jlec/terraform-provider-sshkey/internal/keygen"
)

func TestNewSSHKeyPair(t *testing.T) {
	t.Parallel()

	c := keygen.SSHKeyPairConfig{
		Passphrase: []byte(""),
		Type:       keygen.RSA,
	}

	if _, err := keygen.New(&c); err != nil {
		t.Errorf("error creating SSH key pair: %v", err)
	}
}

func TestGenerateEd25519Keys(t *testing.T) {
	t.Parallel()

	c := keygen.SSHKeyPairConfig{
		Type: keygen.ED25519,
	}
	key, err := keygen.New(&c)

	t.Run("test generate SSH keys", func(t *testing.T) {
		t.Parallel()

		if err != nil {
			t.Errorf("error creating SSH key pair: %v", err)
		}

		// TODO: is there a good way to validate these? Lengths seem to vary a bit,
		// so far now we're just asserting that the keys indeed exist.
		if len(key.PrivateKeyPEM()) == 0 {
			t.Error("error creating SSH private key PEM; key is 0 bytes")
		}
		if len(key.PublicKey()) == 0 {
			t.Error("error creating SSH public key; key is 0 bytes")
		}
	})
}

func TestGenerateECDSAKeys(t *testing.T) {
	t.Parallel()

	// Create temp directory for keys
	conf := keygen.SSHKeyPairConfig{
		Type: keygen.ECDSA,
	}
	key, _ := keygen.New(&conf)

	t.Run("test generate SSH keys", func(t *testing.T) {
		t.Parallel()

		// TODO: is there a good way to validate these? Lengths seem to vary a bit,
		// so far now we're just asserting that the keys indeed exist.
		if len(key.PrivateKeyPEM()) == 0 {
			t.Error("error creating SSH private key PEM; key is 0 bytes")
		}
		if len(key.PublicKey()) == 0 {
			t.Error("error creating SSH public key; key is 0 bytes")
		}
	})
}

func TestGeneratePublicKeyWithEmptyDir(t *testing.T) {
	t.Parallel()

	for _, keyType := range keygen.SSSHKeyTypes {
		func(t *testing.T) {
			t.Helper()

			conf := keygen.SSHKeyPairConfig{Passphrase: nil, Type: keyType}
			key, err := keygen.New(&conf)

			// FIXME: Do some testing
			_ = key

			if err != nil {
				t.Fatalf("error creating SSH key pair: %v", err)
			}
		}(t)
	}
}

func TestGenerateKeyWithPassphrase(t *testing.T) {
	t.Parallel()

	for _, keyType := range keygen.SSSHKeyTypes {
		tpass := "testpass"

		func(t *testing.T) {
			t.Helper()

			conf := keygen.SSHKeyPairConfig{
				Passphrase: []byte(tpass),
				Type:       keyType,
			}

			_, err := keygen.New(&conf)
			if err != nil {
				t.Fatalf("error creating SSH key pair: %v", err)
			}

			conf = keygen.SSHKeyPairConfig{Passphrase: []byte(tpass), Type: keyType}

			key, err := keygen.New(&conf)
			if err != nil {
				t.Fatalf("error reading SSH key pair: %v", err)
			}

			// FIXME: Do some testing
			_ = key
		}(t)
	}
}

func TestReadingKeyWithPassphrase(t *testing.T) {
	t.Parallel()

	for _, keyType := range keygen.SSSHKeyTypes {
		c := keygen.SSHKeyPairConfig{Passphrase: []byte("test"), Type: keyType}
		if _, err := keygen.New(&c); err != nil {
			t.Fatalf("error reading SSH key pair: %v", err)
		}
	}
}
