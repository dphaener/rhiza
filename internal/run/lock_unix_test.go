//go:build unix

package run

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestTerminateStopsGroup: terminating a process that leads its group
// stops the children it started too.
func TestTerminateStopsGroup(t *testing.T) {
	cmd := exec.Command("sh", "-c", "sleep 300 & echo $!; wait")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	line, err := bufio.NewReader(out).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	child, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Kill(child, syscall.SIGKILL)
	if err := Terminate(cmd.Process.Pid); err != nil {
		t.Fatal(err)
	}
	cmd.Wait()
	for deadline := time.Now().Add(5 * time.Second); PidAlive(child); {
		if time.Now().After(deadline) {
			t.Fatalf("child %d outlived its group leader", child)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
