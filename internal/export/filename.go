package export

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func SafeFilename(prefix, extension string, now time.Time) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "export"
	}
	extension = strings.TrimPrefix(strings.TrimSpace(extension), ".")
	if extension == "" {
		extension = "dat"
	}
	return filepath.Join("exports", fmt.Sprintf("%s-%s.%s", prefix, now.UTC().Format("20060102T150405Z"), extension))
}

func IsSupported(format string) bool { return format == "csv" || format == "json" }
