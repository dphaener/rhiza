//go:build darwin

package daemon

import "golang.org/x/sys/unix"

// readLoadavg reads the 1- and 5-minute load averages from sysctl
// vm.loadavg.
func readLoadavg() (l1, l5 float64, ok bool) {
	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		return 0, 0, false
	}
	l1, l5, err = ParseSysctlLoadavg(b)
	return l1, l5, err == nil
}
