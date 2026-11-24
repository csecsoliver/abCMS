package media

import (
	"fmt"
	"os"
	"syscall"
)

// FreeSpaceText reports available disk space on the filesystem containing path.
func FreeSpaceText(path string) string {
	target := path
	if _, err := os.Stat(target); err != nil {
		target = "."
	}
	var stat syscall.Statfs_t
	if err := syscall.Statfs(target, &stat); err != nil {
		return "Free disk space: unavailable"
	}
	availBytes := stat.Bavail * uint64(stat.Bsize)
	return "Free disk space: " + formatBytes(availBytes)
}

// DBSizeBytes returns the size of the SQLite database file in bytes.
func DBSizeBytes(dbPath string) (int64, bool) {
	info, err := os.Stat(dbPath)
	if err != nil {
		return 0, false
	}
	return info.Size(), true
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit && exp < 4; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}
	return fmt.Sprintf("%.1f %s", float64(b)/float64(div), suffixes[exp])
}
