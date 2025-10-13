package session

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"strings"
)

func Name(path string) string {
	basename := filepath.Base(path)
	sessionPrefix := strings.ReplaceAll(basename, ".", "-")
	return strings.Join([]string{sessionPrefix, hash(path)}, "-")
}

func hash(path string) string {
	hash := sha1.New()
	hash.Write([]byte(path))
	hashByteSlice := hash.Sum(nil)
	return fmt.Sprintf("%x", hashByteSlice)[:4]
}
