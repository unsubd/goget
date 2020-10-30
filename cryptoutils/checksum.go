package cryptoutils

import (
	"crypto/sha256"
	"fmt"
	"goget/logging"
	"io"
	"os"
)

func FileChecksumSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logging.LogError("FILE_CHECKSUM_SHA256", err, filePath)
		return "", err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		logging.LogError("FILE_CHECKSUM_SHA256", err, filePath)
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
