/*
Copyright 2022-2024 Justin Lecher

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package keygen_test

import (
	"strings"
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

		if strings.Contains(string(key.PublicKey()), "\n") {
			t.Error("Line break in public key")
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

	for _, keyType := range keygen.SSHKeyTypes {
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

	for _, keyType := range keygen.SSHKeyTypes {
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

	for _, keyType := range keygen.SSHKeyTypes {
		c := keygen.SSHKeyPairConfig{Passphrase: []byte("test"), Type: keyType}
		if _, err := keygen.New(&c); err != nil {
			t.Fatalf("error reading SSH key pair: %v", err)
		}
	}
}
