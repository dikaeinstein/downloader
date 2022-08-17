package hash

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// RemoteHasher generates a checksum from a remote source.
type RemoteHasher struct {
	client *http.Client
}

// Hash fetches the checksum at the given url and returns it as a string.
func (r RemoteHasher) Hash(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %v", url, res.Status)
	}

	checksum, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %v", url, err)
	}

	return string(checksum), nil
}

// NewRemoteHasher returns an initialized RemoteHasher
func NewRemoteHasher(client *http.Client) RemoteHasher {
	return RemoteHasher{client}
}

// LocalHasher is a Hasher that generates a checksum from a local source.
type LocalHasher struct{}

// Hash returns the checksum from the file at the given path.
func (LocalHasher) Hash(ctx context.Context, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	checksum, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(checksum), nil
}

type StringHasher string

func (StringHasher) Hash(ctx context.Context, checksum string) (string, error) {
	return checksum, nil
}

// FakeHasher is a Hasher that always returns the same hash: "fakehash".
type FakeHasher struct{}

func (FakeHasher) Hash(ctx context.Context, path string) (string, error) {
	return "fakehash", nil
}

// Verifier is a type that can verify a hash.
type Verifier struct{}

// VerifyHash reports whether the io.Reader has contents with
// SHA-256 that matches the given hex value.
func (Verifier) Verify(input io.Reader, hex string) error {
	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		return err
	}

	if !strings.EqualFold(hex, fmt.Sprintf("%x", hash.Sum(nil))) {
		return fmt.Errorf(
			"%s corrupt? does not have expected SHA-256 of %v",
			input, hex,
		)
	}

	return nil
}
