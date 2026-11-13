package api

// This file holds the hidden daemon.profile operation: it asks the running
// root daemon to capture a CPU profile of itself, without restarting it
// (a restart would drop the warm caches and the long-uptime state the
// profile is wanted for). The daemon binds no socket, so the request is a
// file in the root's state folder that the daemon reads at each poll; the
// daemon writes the profile, then a result file the operation waits for.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/rootcfg"
)

// The files of a capture, in the root's state folder.
const (
	profileRequestFile = "profile-request.json"
	profileResultFile  = "profile-result.json"
	// profileDefaultSeconds is how long a capture lasts without --seconds.
	profileDefaultSeconds = 30
)

// DaemonProfileRequest asks the running daemon for a CPU profile of itself.
type DaemonProfileRequest struct {
	Seconds int    `json:"seconds,omitempty" desc:"How long to profile, in seconds; default 30"`
	Out     string `json:"out,omitempty" desc:"The file to write the gzip-compressed pprof profile to; default profile-<time>.pprof in the root's .rhiza-root folder"`
}

// DaemonProfileResponse names the profile written.
type DaemonProfileResponse struct {
	Path    string `json:"path" desc:"The profile file, for go tool pprof"`
	Seconds int    `json:"seconds" desc:"How long the daemon was profiled, in seconds"`
}

// Human prints the file's path.
func (r DaemonProfileResponse) Human(w io.Writer) { fmt.Fprintln(w, r.Path) }

func init() {
	register(Operation{
		Name:     "daemon.profile",
		Summary:  "Capture a CPU profile of the running daemon",
		Hidden:   true,
		Request:  reflect.TypeFor[DaemonProfileRequest](),
		Response: reflect.TypeFor[DaemonProfileResponse](),
		Run:      runDaemonProfile,
	})
}

// profileWaitEvery is how often the operation looks for the result.
var profileWaitEvery = 200 * time.Millisecond

// profileSlack is how long past the capture the operation waits for the
// daemon's next poll to pick the request up and finish it.
var profileSlack = 2*daemonPoll + 10*time.Second

type profileAsk struct {
	Seconds int    `json:"seconds"`
	Out     string `json:"out"`
	At      string `json:"at"`
}

type profileResult struct {
	Path  string `json:"path,omitempty"`
	Error string `json:"error,omitempty"`
}

func runDaemonProfile(ctx *Ctx, req any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("daemon profile writes no log entry and has no dry run")
	}
	r := req.(DaemonProfileRequest)
	if r.Seconds < 0 {
		return nil, Usagef("--seconds must be at least 1")
	}
	if r.Seconds == 0 {
		r.Seconds = profileDefaultSeconds
	}
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		return nil, Usagef("%s is not inside a Rhiza root", ctx.Wd).WithNext("rhiza root init")
	}
	if _, running := daemon.RunningRoot(root); !running {
		return nil, Usagef("no daemon is running at %s", root).WithNext("rhiza daemon start")
	}
	dir := daemon.RootStateDir(root)
	out := r.Out
	if out == "" {
		out = filepath.Join(dir, "profile-"+ctx.Now().UTC().Format("20060102T150405Z")+".pprof")
	} else if !filepath.IsAbs(out) {
		out = filepath.Join(ctx.Wd, out)
	}
	ask, err := json.Marshal(profileAsk{Seconds: r.Seconds, Out: out, At: ctx.Now().UTC().Format(time.RFC3339)})
	if err != nil {
		return nil, err
	}
	_ = os.Remove(filepath.Join(dir, profileResultFile))
	// The request file is created exclusively: a second request while one
	// is in progress finds it there.
	f, err := os.OpenFile(filepath.Join(dir, profileRequestFile), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return nil, Conflictf("a profile capture is already in progress").WithNext("rhiza daemon status")
	}
	if err != nil {
		return nil, Errorf("write the profile request: %v", err)
	}
	_, werr := f.Write(ask)
	if cerr := f.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		_ = os.Remove(filepath.Join(dir, profileRequestFile))
		return nil, Errorf("write the profile request: %v", werr)
	}
	wait := time.Duration(r.Seconds)*time.Second + profileSlack
	deadline := time.Now().Add(wait)
	for {
		b, err := os.ReadFile(filepath.Join(dir, profileResultFile))
		if err == nil {
			_ = os.Remove(filepath.Join(dir, profileResultFile))
			var res profileResult
			if err := json.Unmarshal(b, &res); err != nil {
				return nil, Errorf("read the profile result: %v", err)
			}
			if res.Error != "" {
				return nil, Errorf("the daemon could not capture a profile: %s", res.Error)
			}
			return DaemonProfileResponse{Path: res.Path, Seconds: r.Seconds}, nil
		}
		if time.Now().After(deadline) {
			// Leave the request: the daemon may still pick it up, and a
			// new request is refused until it has.
			return nil, Errorf("the daemon did not finish the capture within %s; the request stays queued", wait.Round(time.Second))
		}
		time.Sleep(profileWaitEvery)
	}
}

// daemonProfiler is the daemon's side: Poll, called at each poll, starts
// the capture a request file asks for. The capture ends on its own timer, so
// polling goes on meanwhile.
type daemonProfiler struct {
	dir string
	// Start and Stop begin and end the CPU profile; After is the capture's
	// clock. Tests replace them.
	Start func(io.Writer) error
	Stop  func()
	After func(time.Duration) <-chan time.Time

	mu        sync.Mutex
	capturing bool
}

func newDaemonProfiler(dir string) *daemonProfiler {
	return &daemonProfiler{dir: dir, Start: pprof.StartCPUProfile, Stop: pprof.StopCPUProfile, After: time.After}
}

// Reset drops a request and result an earlier daemon left behind.
func (p *daemonProfiler) Reset() {
	_ = os.Remove(filepath.Join(p.dir, profileRequestFile))
	_ = os.Remove(filepath.Join(p.dir, profileResultFile))
}

func (p *daemonProfiler) finish(res profileResult) {
	b, _ := json.Marshal(res)
	tmp := filepath.Join(p.dir, profileResultFile+".tmp")
	if os.WriteFile(tmp, b, 0o644) == nil {
		_ = os.Rename(tmp, filepath.Join(p.dir, profileResultFile))
	}
	_ = os.Remove(filepath.Join(p.dir, profileRequestFile))
}

// Poll starts the capture a request file asks for, unless one runs.
func (p *daemonProfiler) Poll() {
	p.mu.Lock()
	if p.capturing {
		p.mu.Unlock()
		return
	}
	b, err := os.ReadFile(filepath.Join(p.dir, profileRequestFile))
	if err != nil || len(b) == 0 {
		p.mu.Unlock()
		return
	}
	var ask profileAsk
	if err := json.Unmarshal(b, &ask); err != nil || ask.Seconds < 1 || ask.Out == "" {
		p.mu.Unlock()
		p.finish(profileResult{Error: "unreadable profile request"})
		return
	}
	f, err := os.Create(ask.Out)
	if err != nil {
		p.mu.Unlock()
		p.finish(profileResult{Error: err.Error()})
		return
	}
	if err := p.Start(f); err != nil {
		p.mu.Unlock()
		f.Close()
		_ = os.Remove(ask.Out)
		p.finish(profileResult{Error: err.Error()})
		return
	}
	p.capturing = true
	p.mu.Unlock()
	timer := p.After(time.Duration(ask.Seconds) * time.Second)
	go func() {
		<-timer
		p.Stop()
		err := f.Close()
		// Clear the request before the flag, so no poll in between takes
		// the finished request for a new one.
		if err != nil {
			p.finish(profileResult{Error: err.Error()})
		} else {
			p.finish(profileResult{Path: ask.Out})
		}
		p.mu.Lock()
		p.capturing = false
		p.mu.Unlock()
	}()
}
