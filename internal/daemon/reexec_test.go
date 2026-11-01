package daemon

import (
	"testing"
	"time"
)

func fakeReexec(t *testing.T, common string, img *Image, execs *int) *Reexec {
	r := &Reexec{
		Path: "/bin/rhiza", Args: []string{"rhiza", "daemon"}, Env: nil,
		Stat: func(string) (Image, error) { return *img, nil },
		Exec: func(path string, argv, env []string) error {
			*execs++
			// The new image takes the lock: it must be free by now.
			l, err := TryLock(common, time.Now())
			if err != nil {
				t.Fatalf("new image cannot take the lock: %v", err)
			}
			l.Unlock()
			return nil
		},
	}
	r.Start(*img)
	return r
}

func TestReexec(t *testing.T) {
	common := t.TempDir()
	now := time.Now()
	img := Image{Size: 10, Mod: now}
	execs := 0
	r := fakeReexec(t, common, &img, &execs)

	if r.Due(0) {
		t.Fatal("unchanged executable is due")
	}
	img.Size = 11
	if r.Due(1) {
		t.Fatal("a pending start must delay the exec")
	}
	if !r.Due(0) {
		t.Fatal("changed executable with nothing pending is not due")
	}
	img = Image{Size: 10, Mod: now.Add(time.Second)}
	if !r.Due(0) {
		t.Fatal("changed mtime is not due")
	}

	lock, err := TryLock(common, now)
	if err != nil {
		t.Fatal(err)
	}
	lock, err = r.Do(common, lock, now)
	if err != nil || lock == nil || execs != 1 {
		t.Fatalf("Do = %v, %v, execs %d", lock, err, execs)
	}
	if _, ok := Running(common); !ok {
		t.Fatal("lock not held again after a returned exec")
	}
	lock.Unlock()
}
