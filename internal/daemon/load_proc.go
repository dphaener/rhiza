//go:build !darwin

package daemon

import "os"

// readLoadavg reads the 1- and 5-minute load averages from /proc/loadavg,
// where there is one.
func readLoadavg() (l1, l5 float64, ok bool) {
	b, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0, false
	}
	l1, l5, err = ParseLoadavg(string(b))
	return l1, l5, err == nil
}
