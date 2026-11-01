package daemon

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultMinFree is the free space under which the daemon starts and
// resumes nothing: 10 GiB.
const DefaultMinFree uint64 = 10 << 30

// ParseSize reads a byte count: plain bytes, or a number with a unit
// (K, M, G, T, optionally followed by iB or B; all powers of 1024).
func ParseSize(s string) (uint64, error) {
	t := strings.ToUpper(strings.TrimSpace(s))
	t = strings.TrimSuffix(t, "B")
	t = strings.TrimSuffix(t, "I")
	mult := uint64(1)
	if n := len(t); n > 0 {
		if i := strings.IndexByte("KMGT", t[n-1]); i >= 0 {
			mult = 1 << (10 * uint(i+1))
			t = t[:n-1]
		}
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("%q is not a size: use bytes or a number with K, M, G or T", s)
	}
	return uint64(v * float64(mult)), nil
}

// FormatSize writes n bytes in the largest binary unit, e.g. 1.0 GiB.
func FormatSize(n uint64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	f := float64(n)
	i := 0
	for f >= 1024 && i < len(units)-1 {
		f /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d B", n)
	}
	return fmt.Sprintf("%.1f %s", f, units[i])
}
