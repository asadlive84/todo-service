package helper

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"
)

func ComputeFileHash(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}

func ComputeFileHashFromReader(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}



// Redis usually returns []byte. This converts everything to string safely.
func ToString(i interface{}) string {
	switch v := i.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}