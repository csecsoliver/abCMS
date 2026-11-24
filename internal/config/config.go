package config

import (
	"crypto/rand"
	"log"
	"os"
	"strings"
	"time"
)

// Config holds runtime configuration for the abCMS server.
type Config struct {
	Addr           string
	DBPath         string
	UploadDir      string
	StaticDir      string
	SessionSecret  []byte
	SessionTTL     time.Duration
	MaxUploadBytes int64
	SecureCookies  bool
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	addr := strings.TrimSpace(os.Getenv("ABCMS_ADDR"))
	if addr == "" {
		if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
			if strings.HasPrefix(port, ":") {
				addr = port
			} else {
				addr = ":" + port
			}
		} else {
			addr = ":8080"
		}
	}

	dbPath := strings.TrimSpace(os.Getenv("ABCMS_DB_PATH"))
	if dbPath == "" {
		dbPath = "abcms.sqlite"
	}

	staticDir := strings.TrimSpace(os.Getenv("ABCMS_STATIC_DIR"))
	if staticDir == "" {
		staticDir = "static"
	}

	uploadDir := strings.TrimSpace(os.Getenv("ABCMS_UPLOAD_DIR"))
	if uploadDir == "" {
		uploadDir = "static/uploads"
	}

	var secret []byte
	if rawSecret := os.Getenv("ABCMS_SESSION_SECRET"); rawSecret != "" {
		secret = []byte(rawSecret)
	} else {
		secret = make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			log.Fatalf("failed to generate ephemeral session secret: %v", err)
		}
		log.Println("WARNING: ABCMS_SESSION_SECRET is not set; using ephemeral session secret (sessions will reset on restart)")
	}

	secureEnv := strings.ToLower(strings.TrimSpace(os.Getenv("ABCMS_SECURE_COOKIES")))
	secureCookies := secureEnv == "true" || secureEnv == "1"

	return Config{
		Addr:           addr,
		DBPath:         dbPath,
		UploadDir:      uploadDir,
		StaticDir:      staticDir,
		SessionSecret:  secret,
		SessionTTL:     100 * time.Hour, // 360000 seconds, matching legacy Lapis config
		MaxUploadBytes: 15 << 20,        // 15 MiB, matching legacy nginx.conf
		SecureCookies:  secureCookies,
	}
}
