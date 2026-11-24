package media

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makeTestPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 64, 48))
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{R: 120, G: 200, B: 80, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func TestSaveUploadAndThumbnail(t *testing.T) {
	ctx := context.Background()
	uploadDir := filepath.Join(t.TempDir(), "uploads")
	svc := NewService(uploadDir)

	pngData := makeTestPNG(t)
	imageURL, thumbURL, err := svc.SaveUpload(ctx, "../../test image.png", pngData)
	if err != nil {
		t.Fatalf("SaveUpload error: %v", err)
	}
	if !strings.HasPrefix(imageURL, "/static/uploads/") || !strings.HasPrefix(thumbURL, "/static/uploads/thumb-") {
		t.Fatalf("unexpected URLs: image=%q thumb=%q", imageURL, thumbURL)
	}
	if strings.Contains(imageURL, "..") {
		t.Fatalf("path traversal leaked into imageURL: %q", imageURL)
	}

	// Verify full-size image was written verbatim (preserving metadata).
	savedName := strings.TrimPrefix(imageURL, "/static/uploads/")
	savedBytes, err := os.ReadFile(filepath.Join(uploadDir, savedName))
	if err != nil {
		t.Fatalf("read saved upload: %v", err)
	}
	if !bytes.Equal(savedBytes, pngData) {
		t.Fatal("saved image bytes differ from original upload")
	}

	if size, ok := svc.FileSizeBytes(imageURL); !ok || size != int64(len(pngData)) {
		t.Fatalf("FileSizeBytes(image) = (%d, %v), want (%d, true)", size, ok, len(pngData))
	}
	if size, ok := svc.FileSizeBytes(thumbURL); !ok || size == 0 {
		t.Fatalf("FileSizeBytes(thumb) = (%d, %v), want (>0, true)", size, ok)
	}

	lazyThumb, err := svc.GenerateThumbnail(ctx, imageURL)
	if err != nil {
		t.Fatalf("GenerateThumbnail error: %v", err)
	}
	if _, ok := svc.FileSizeBytes(lazyThumb); !ok {
		t.Fatal("lazy thumbnail file not found on disk")
	}

	svc.DeleteFile(imageURL)
	svc.DeleteFile(thumbURL)
	svc.DeleteFile(lazyThumb)

	if _, ok := svc.FileSizeBytes(imageURL); ok {
		t.Fatal("expected deleted image to be removed")
	}
}

func TestSaveUploadRejectsInvalidFormat(t *testing.T) {
	svc := NewService(t.TempDir())
	_, _, err := svc.SaveUpload(context.Background(), "not-an-image.txt", []byte("plain text not an image"))
	if !errors.Is(err, ErrUnsupportedFormat) {
		t.Fatalf("expected ErrUnsupportedFormat, got %v", err)
	}
}

func TestFreeSpaceText(t *testing.T) {
	txt := FreeSpaceText(t.TempDir())
	if !strings.HasPrefix(txt, "Free disk space: ") {
		t.Fatalf("unexpected FreeSpaceText: %q", txt)
	}
}
