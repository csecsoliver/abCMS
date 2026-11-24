package media

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// ErrUnsupportedFormat indicates that an uploaded file is not an allowed image format.
var ErrUnsupportedFormat = errors.New("unsupported image format")

var allowedFormats = map[string]bool{
	"JPEG": true,
	"JPG":  true,
	"PNG":  true,
	"GIF":  true,
	"WEBP": true,
	"BMP":  true,
	"TIFF": true,
	"TIF":  true,
	"ICO":  true,
	"CUR":  true,
	"SVG":  true,
	"SVGZ": true,
	"APNG": true,
	"AVIF": true,
	"JXL":  true,
	"HEIC": true,
	"HEIF": true,
}

// Service manages image uploads, thumbnails, and file cleanup under UploadDir.
type Service struct {
	UploadDir string
	URLPrefix string
}

// NewService creates a media Service rooted at uploadDir.
func NewService(uploadDir string) *Service {
	return &Service{
		UploadDir: uploadDir,
		URLPrefix: "/static/uploads/",
	}
}

// SaveUpload validates an image upload, converts HEIC/HEIF to JPEG, preserves metadata,
// writes the full-size image and a 1000x1000 JPEG thumbnail, and returns their URL paths.
func (s *Service) SaveUpload(ctx context.Context, origFilename string, data []byte) (string, string, error) {
	if len(data) == 0 {
		return "", "", ErrUnsupportedFormat
	}
	if err := os.MkdirAll(s.UploadDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create upload dir: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "abcms-upload-*")
	if err != nil {
		return "", "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return "", "", fmt.Errorf("write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", "", fmt.Errorf("close temp file: %w", err)
	}

	format, err := identifyFormat(ctx, tmpPath)
	if err != nil || !allowedFormats[format] {
		return "", "", ErrUnsupportedFormat
	}

	baseName := sanitizeFilename(origFilename)
	idPrefix := uuid.NewString()
	filename := idPrefix + "-" + baseName

	if format == "HEIC" || format == "HEIF" {
		if !strings.HasSuffix(strings.ToLower(filename), ".jpg") && !strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
			filename += ".jpg"
		}
		destPath := filepath.Join(s.UploadDir, filename)
		// Convert HEIC/HEIF without -strip so EXIF/metadata is intentionally preserved.
		if err := runMagick(ctx, tmpPath+"[0]", destPath); err != nil {
			return "", "", fmt.Errorf("convert %s to jpeg: %w", format, err)
		}
	} else {
		destPath := filepath.Join(s.UploadDir, filename)
		// Write original bytes verbatim to preserve image metadata.
		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			return "", "", fmt.Errorf("write upload file: %w", err)
		}
	}

	destPath := filepath.Join(s.UploadDir, filename)
	thumbFilename := "thumb-" + filename + ".jpg"
	thumbPath := filepath.Join(s.UploadDir, thumbFilename)
	if err := makeThumbnail(ctx, destPath, thumbPath); err != nil {
		_ = os.Remove(destPath)
		return "", "", fmt.Errorf("generate thumbnail: %w", err)
	}

	return s.URLPrefix + filename, s.URLPrefix + thumbFilename, nil
}

// GenerateThumbnail creates a 1000x1000 JPEG thumbnail for an existing uploaded image URL.
func (s *Service) GenerateThumbnail(ctx context.Context, imageURL string) (string, error) {
	srcPath, ok := s.resolveURLPath(imageURL)
	if !ok {
		return "", fmt.Errorf("invalid upload path %q", imageURL)
	}
	if err := os.MkdirAll(s.UploadDir, 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	thumbFilename := "thumb-" + uuid.NewString() + ".jpg"
	thumbPath := filepath.Join(s.UploadDir, thumbFilename)
	if err := makeThumbnail(ctx, srcPath, thumbPath); err != nil {
		return "", fmt.Errorf("generate thumbnail: %w", err)
	}
	return s.URLPrefix + thumbFilename, nil
}

// DeleteFile removes an uploaded image or thumbnail file if it lives within UploadDir.
func (s *Service) DeleteFile(imageURL string) {
	targetPath, ok := s.resolveURLPath(imageURL)
	if !ok {
		return
	}
	_ = os.Remove(targetPath)
}

// FileSizeBytes returns the byte size of an uploaded file if it exists within UploadDir.
func (s *Service) FileSizeBytes(imageURL string) (int64, bool) {
	targetPath, ok := s.resolveURLPath(imageURL)
	if !ok {
		return 0, false
	}
	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		return 0, false
	}
	return info.Size(), true
}

func (s *Service) resolveURLPath(imageURL string) (string, bool) {
	if imageURL == "" || !strings.HasPrefix(imageURL, s.URLPrefix) {
		return "", false
	}
	rel := strings.TrimPrefix(imageURL, s.URLPrefix)
	if rel == "" || strings.Contains(rel, "..") {
		return "", false
	}
	absUploadDir, err := filepath.Abs(s.UploadDir)
	if err != nil {
		return "", false
	}
	candidate := filepath.Join(absUploadDir, filepath.FromSlash(rel))
	relCheck, err := filepath.Rel(absUploadDir, candidate)
	if err != nil || relCheck == "." || strings.HasPrefix(relCheck, "..") {
		return "", false
	}
	return candidate, true
}

func identifyFormat(ctx context.Context, filePath string) (string, error) {
	cmd := exec.CommandContext(ctx, "magick", "identify", "-format", "%m", filePath+"[0]")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.ToUpper(strings.TrimSpace(string(out))), nil
}

func makeThumbnail(ctx context.Context, srcPath, destThumbPath string) error {
	return runMagick(ctx, srcPath+"[0]", "-auto-orient", "-resize", "1000x1000>", destThumbPath)
}

func runMagick(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "magick", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func sanitizeFilename(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" || base == "." || base == "/" || base == `\` {
		return "image"
	}
	var b strings.Builder
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '.', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	cleaned := strings.Trim(b.String(), "._")
	if cleaned == "" {
		return "image"
	}
	return cleaned
}
