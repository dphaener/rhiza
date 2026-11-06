package daemon

import (
	"errors"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"strings"
	"testing"
	"time"
)

// fakeRepo is a main checkout whose git and make are scripted.
type fakeRepo struct {
	head, remote string
	branch       string
	dirty        bool
	ahead        bool // checkout is ahead of origin/main
	ffErr, mkErr error
	calls, makes int
	fetches, ffs int
	logs         []string
}

func (f *fakeRepo) su() *SelfUpdate {
	return &SelfUpdate{
		Git: func(args ...string) (string, error) {
			f.calls++
			switch args[0] {
			case "fetch":
				f.fetches++
			case "rev-parse":
				if args[2] == "HEAD" {
					return f.head + "\n", nil
				}
				return f.remote + "\n", nil
			case "merge-base":
				if f.ahead {
					return "", nil
				}
				return "", errors.New("exit 1")
			case "symbolic-ref":
				return f.branch + "\n", nil
			case "status":
				if f.dirty {
					return " M x.go\n", nil
				}
				return "", nil
			case "merge":
				f.ffs++
				if f.ffErr != nil {
					return "", f.ffErr
				}
				f.head = f.remote
			}
			return "", nil
		},
		Make: func(args ...string) error {
			if len(args) != 1 || args[0] != "update" {
				panic("make " + strings.Join(args, " "))
			}
			f.makes++
			return f.mkErr
		},
		Log: func(format string, a ...any) { f.logs = append(f.logs, format) },
	}
}

func TestSelfUpdateMovedMainRunsOnce(t *testing.T) {
	f := &fakeRepo{head: "aaa", remote: "bbb", branch: "main"}
	u := f.su()
	if !u.Run() || f.makes != 1 || f.ffs != 1 || f.head != "bbb" {
		t.Fatalf("makes %d ffs %d head %s", f.makes, f.ffs, f.head)
	}
	u.Run()
	u.Run()
	if f.makes != 1 || f.fetches != 3 {
		t.Fatalf("makes %d fetches %d after unmoved polls", f.makes, f.fetches)
	}
	f.remote = "ccc"
	u.Run()
	if f.makes != 2 {
		t.Fatalf("makes %d after main moved again", f.makes)
	}
}

func TestSelfUpdateUnmovedRunsNothing(t *testing.T) {
	for name, f := range map[string]*fakeRepo{
		"equal": {head: "aaa", remote: "aaa", branch: "main"},
		"ahead": {head: "aaa", remote: "bbb", branch: "main", ahead: true},
	} {
		if f.su().Run() || f.makes != 0 || f.ffs != 0 {
			t.Errorf("%s: ran", name)
		}
	}
}

func TestSelfUpdateSkipsAndWaitsForMove(t *testing.T) {
	cases := map[string]*fakeRepo{
		"dirty":    {head: "aaa", remote: "bbb", branch: "main", dirty: true},
		"off main": {head: "aaa", remote: "bbb", branch: "topic"},
		"diverged": {head: "aaa", remote: "bbb", branch: "main", ffErr: errors.New("not possible to fast-forward")},
		"make":     {head: "aaa", remote: "bbb", branch: "main", mkErr: errors.New("exit status 2")},
	}
	for name, f := range cases {
		u := f.su()
		u.Run()
		if len(f.logs) == 0 {
			t.Errorf("%s: nothing logged", name)
		}
		ffs, makes := f.ffs, f.makes
		f.dirty, f.branch, f.ffErr, f.mkErr = false, "main", nil, nil
		u.Run()
		if f.ffs != ffs || f.makes != makes {
			t.Errorf("%s: retried before main moved", name)
		}
		f.remote = "ccc"
		u.Run()
		if f.makes != makes+1 {
			t.Errorf("%s: not retried after main moved", name)
		}
	}
}

func TestSelfUpdateFetchesAtMostEvery(t *testing.T) {
	f := &fakeRepo{head: "aaa", remote: "aaa", branch: "main"}
	u := f.su()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	u.Now = func() time.Time { return now }
	u.Every = time.Minute
	u.Run()
	now = now.Add(30 * time.Second)
	u.Run()
	if f.fetches != 1 {
		t.Fatalf("fetches %d within Every", f.fetches)
	}
	now = now.Add(30 * time.Second)
	u.Run()
	if f.fetches != 2 {
		t.Fatalf("fetches %d after Every", f.fetches)
	}
}

func TestSelfUpdateUsesFetch(t *testing.T) {
	f := &fakeRepo{head: "aaa", remote: "bbb", branch: "main"}
	u := f.su()
	u.Fetch = func() error { return errors.New("timed out") }
	if u.Run() || f.fetches != 0 || f.makes != 0 || len(f.logs) != 1 {
		t.Fatalf("fetches %d makes %d logs %v", f.fetches, f.makes, f.logs)
	}
}

func TestSelfUpdateProjectChoosesCheckout(t *testing.T) {
	root := "/r"
	reg := rootcfg.Config{
		SelfUpdateProject: "rhiza",
		Projects:          map[string]rootcfg.Project{"rhiza": {Path: "code/rhiza"}, "other": {Path: "o"}},
	}
	var dirs []string
	git := func(dir string, args ...string) (string, error) {
		dirs = append(dirs, "git:"+dir)
		switch args[0] {
		case "rev-parse":
			if args[2] == "HEAD" {
				return "aaa", nil
			}
			return "bbb", nil
		case "merge-base":
			return "", errors.New("exit 1")
		case "symbolic-ref":
			return "main", nil
		}
		return "", nil
	}
	mk := func(dir string, args ...string) error { dirs = append(dirs, "make:"+dir); return nil }
	u := &SelfUpdate{}
	if !u.ForCheckout(root, reg, git, mk) {
		t.Fatal("named project not resolved")
	}
	u.Run()
	for _, d := range dirs {
		if !strings.HasSuffix(d, ":/r/code/rhiza") {
			t.Fatalf("ran outside the project: %s", d)
		}
	}
	if len(dirs) == 0 || dirs[len(dirs)-1] != "make:/r/code/rhiza" {
		t.Fatalf("calls %v", dirs)
	}

	dirs = nil
	reg.SelfUpdateProject = ""
	u = &SelfUpdate{}
	if u.ForCheckout(root, reg, git, mk) {
		t.Fatal("unset still resolved")
	}
	u.Run()
	if len(dirs) != 0 {
		t.Fatalf("ran with none named: %v", dirs)
	}
}
