package hash

import (
	"fmt"
	"os"

	"github.com/twmb/murmur3"
)

// Returns the murmur3 hash of path
func File(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := murmur3.SeedSum64(1337, data)
	return fmt.Sprintf("%X", hash), nil
}

// Returns the murmur3 hash of s
func String(s string) string {
	hash := murmur3.SeedStringSum64(1337, s)
	return fmt.Sprintf("%X", hash)
}
