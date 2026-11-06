//go:build linux

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func TestReap(t *testing.T) {
	common := t.TempDir()
	list := filepath.Join(common, "rhiza", "drivers")

	// A plain os/exec child that calls setsid() itself, running meanwhile.
	plain := exec.Command("sh", "-c", "sleep 1; exit 3")
	plain.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := plain.Start(); err != nil {
		t.Fatal(err)
	}
	// A listed detached child nothing waits on.
	drv := exec.Command("true")
	drv.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := drv.Start(); err != nil {
		t.Fatal(err)
	}
	pid := drv.Process.Pid
	addDriver(list, pid)

	deadline := time.Now().Add(5 * time.Second)
	for !isZombie(pid) && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if !isZombie(pid) {
		t.Fatal("child did not become a zombie")
	}
	reapDrivers(common)
	if _, err := os.Stat("/proc/" + itoa(pid)); err == nil {
		t.Errorf("pid %d still exists after reap", pid)
	}
	if got := readDrivers(list); len(got) != 0 {
		t.Errorf("driver list = %v, want empty", got)
	}
	err := plain.Wait()
	ee, ok := err.(*exec.ExitError)
	if !ok || ee.ExitCode() != 3 {
		t.Errorf("plain child Wait = %v, want exit status 3", err)
	}
}

func isZombie(pid int) bool {
	b, err := os.ReadFile("/proc/" + itoa(pid) + "/stat")
	if err != nil {
		return false
	}
	s := string(b)
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == ')' {
			return len(s) > i+2 && s[i+2] == 'Z'
		}
	}
	return false
}

func itoa(n int) string { return strconv.Itoa(n) }
