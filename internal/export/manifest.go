package export

import (
	"fmt"
	"strings"
	"time"
)

type Manifest struct {
	Name        string
	Format      string
	GeneratedAt time.Time
	Rows        int
	Checksum    string
}

func NewManifest(name, format string, rows int, checksum string, now time.Time) Manifest {
	return Manifest{Name: strings.TrimSpace(name), Format: strings.ToLower(strings.TrimSpace(format)), GeneratedAt: now.UTC(), Rows: rows, Checksum: checksum}
}

func (m Manifest) Valid() bool {
	return m.Name != "" && IsSupported(m.Format) && m.Rows >= 0 && !m.GeneratedAt.IsZero()
}

func (m Manifest) String() string {
	return fmt.Sprintf("%s.%s rows=%d checksum=%s", m.Name, m.Format, m.Rows, m.Checksum)
}
