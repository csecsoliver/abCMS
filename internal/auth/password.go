package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// HashPassword hashes a plaintext password with bcrypt at cost 12.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash failed: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a plaintext password against a stored bcrypt or legacy Argon2 PHC hash.
func VerifyPassword(password, encodedHash string) bool {
	if encodedHash == "" {
		return false
	}
	if strings.HasPrefix(encodedHash, "$argon2") {
		return verifyArgon2(password, encodedHash)
	}
	return bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password)) == nil
}

func verifyArgon2(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	// Standard PHC format: ["", "argon2id"|"argon2i", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(parts) != 6 {
		return false
	}
	variant := parts[1]
	if variant != "argon2id" && variant != "argon2i" {
		return false
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false
	}

	var memory, iterations uint32
	var parallelism uint8
	for _, kv := range strings.Split(parts[3], ",") {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			return false
		}
		val, err := strconv.ParseUint(pair[1], 10, 32)
		if err != nil {
			return false
		}
		switch pair[0] {
		case "m":
			memory = uint32(val)
		case "t":
			iterations = uint32(val)
		case "p":
			if val == 0 || val > 255 {
				return false
			}
			parallelism = uint8(val)
		default:
			return false
		}
	}
	if memory == 0 || iterations == 0 || parallelism == 0 {
		return false
	}

	salt, err := decodePHCBase64(parts[4])
	if err != nil || len(salt) == 0 {
		return false
	}
	expectedHash, err := decodePHCBase64(parts[5])
	if err != nil || len(expectedHash) == 0 {
		return false
	}

	keyLen := uint32(len(expectedHash))
	var derived []byte
	switch variant {
	case "argon2id":
		derived = argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)
	case "argon2i":
		derived = argon2.Key([]byte(password), salt, iterations, memory, parallelism, keyLen)
	}

	return subtle.ConstantTimeCompare(derived, expectedHash) == 1
}

func decodePHCBase64(s string) ([]byte, error) {
	if strings.ContainsRune(s, '=') {
		return base64.StdEncoding.DecodeString(s)
	}
	return base64.RawStdEncoding.DecodeString(s)
}
