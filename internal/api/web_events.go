package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// inboxEvery bounds how often the inbox is worked out for a stream while the
// log and the journals stand still: it is the dearest source.
const inboxEvery = 5 * time.Second

// WebEventSources are the sources of the web UI's /events stream, read from
// the repository base works in: the change-set log, the runs' journals and
// driver logs, the inbox and the daemon's state file. Each looks at the
// files as they are when called. With health, the daemon event also moves
// when a new health report lands.
func WebEventSources(base *Ctx, health *WebHealth) web.Sources {
	common := func() string {
		c := *base
		d, err := gitCommonDir(&c)
		if err != nil {
			return ""
		}
		return d
	}
	var mu sync.Mutex
	var inboxAt time.Time
	var inboxSig string
	var inboxLast map[string]string
	sources := web.Sources{
		web.EventChange: func() map[string]string {
			return dirFingerprints(filepath.Join(base.Root, ".rhiza", "log"), "", true)
		},
		web.EventRun: func() map[string]string {
			if c := common(); c != "" {
				return dirFingerprints(run.Dir(c), ".jsonl", false)
			}
			return nil
		},
		web.EventLog: func() map[string]string {
			if c := common(); c != "" {
				return dirFingerprints(filepath.Join(c, "rhiza", "logs"), ".log", false)
			}
			return nil
		},
		web.EventInbox: func() map[string]string {
			mu.Lock()
			defer mu.Unlock()
			c := common()
			sig := fmt.Sprint(dirFingerprints(filepath.Join(base.Root, ".rhiza", "log"), "", true),
				dirFingerprints(run.Dir(c), ".jsonl", false),
				dirFingerprints(daemon.Dir(c), daemon.StateFile, false))
			if inboxLast != nil && sig == inboxSig && time.Since(inboxAt) < inboxEvery {
				return inboxLast
			}
			cc := *base
			cc.Record, cc.DryRun, cc.Out, cc.JSON = false, false, nil, true
			cc.Stdin, cc.Rex, cc.Start, cc.Detach, cc.Interrupts = nil, nil, nil, nil, nil
			got := map[string]string{}
			if resp, err := runInbox(&cc, nil); err == nil {
				for _, it := range resp.(InboxResponse).Items {
					got[it.ID] = it.Kind + " " + it.Detail + " " + it.Question + " " + it.Recommendation
				}
			} else if inboxLast != nil {
				got = inboxLast
			}
			inboxAt, inboxSig, inboxLast = time.Now(), sig, got
			return got
		},
		web.EventDaemon: func() map[string]string {
			got := map[string]string{}
			if health != nil {
				got["health"] = health.cache.fingerprint()
			}
			c := common()
			if c == "" {
				return got
			}
			s, ok := daemon.ReadState(c)
			if !ok {
				got["daemon"] = "none"
				return got
			}
			var runs []string
			for _, r := range s.Runs {
				runs = append(runs, r.ID+":"+r.State)
			}
			paused := ""
			for _, v := range s.Settings {
				if v.Key == "rhiza.daemon.paused" {
					paused = v.Value
				}
			}
			got["daemon"] = fmt.Sprintf("%s|%v|paused=%s|max=%d|runs=%s|disk=%v|usage=%v|stopped=%v|skips=%d",
				s.Process, s.Autonomous, paused, s.MaxRuns, strings.Join(runs, ","), s.DiskLow != nil, s.UsageHold != nil, !s.Stopped.IsZero(), len(s.Skips))
			return got
		},
	}
	return sources
}

// dirFingerprints fingerprints the entries of dir: with suffix, the files
// ending in it (by name without it) by size and mtime; with dirs, the
// subdirectories, by their entries' count and last name, so an entry
// appended to a change set's existing directory moves it. A missing dir has
// none.
func dirFingerprints(dir, suffix string, dirs bool) map[string]string {
	got := map[string]string{}
	es, err := os.ReadDir(dir)
	if err != nil {
		return got
	}
	for _, e := range es {
		n := e.Name()
		if strings.HasPrefix(n, ".") {
			continue
		}
		if dirs {
			if e.IsDir() {
				sub, _ := os.ReadDir(filepath.Join(dir, n))
				last := ""
				if len(sub) > 0 {
					last = sub[len(sub)-1].Name()
				}
				got[n] = fmt.Sprintf("%d:%s", len(sub), last)
			}
			continue
		}
		if e.IsDir() || !strings.HasSuffix(n, suffix) {
			continue
		}
		if fi, err := e.Info(); err == nil {
			got[strings.TrimSuffix(n, suffix)] = fmt.Sprintf("%d@%d", fi.Size(), fi.ModTime().UnixNano())
		}
	}
	return got
}

// DiskLowInboxIDs are the IDs of the inbox items a disk-low state makes;
// tests use it to name the event they wait for.
func DiskLowInboxIDs(low *daemon.DiskLow) []string {
	var ids []string
	for _, it := range diskLowItems(low) {
		ids = append(ids, it.ID)
	}
	return ids
}
