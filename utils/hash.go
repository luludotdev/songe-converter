package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// CalculateMD5 calculates an MD5 hex string
func CalculateMD5(bytes []byte) string {
	sum := md5.Sum(bytes)
	return hex.EncodeToString(sum[:])
}

// CalculateSHA1 calculates a SHA1 hex string
func CalculateSHA1(bytes []byte) string {
	sum := sha1.Sum(bytes)
	return hex.EncodeToString(sum[:])
}
