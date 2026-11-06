package api

// This file holds run.tail: the last lines of a headless run's driver log.
// It is the read-only part of run.log (no paging cursor, no Rex layer) that
// the web UI's run page uses.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.tail",
		Summary:  "Print the last lines of a headless run's driver log; exits 2 for a run with no log",
		Request:  reflect.TypeFor[RunTailRequest](),
		Response: reflect.TypeFor[RunTailResponse](),
		ReadOnly: true,
		Run:      runRunTail,
	})
}

const (
	runTailDefaultLimit = 100
	runTailMaxLimit     = 1000
	// runTailWindow bounds how much of the file's end a call reads.
	runTailWindow = 256 * 1024
)

// RunTailRequest names the run and how many lines.
type RunTailRequest struct {
	Run   string `json:"run,omitempty" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run"`
	Limit int    `json:"limit,omitempty" desc:"The most lines to return, 1 to 1000; default 100"`
}

// RunTailResponse is the log's last lines.
type RunTailResponse struct {
	Run   string   `json:"run" desc:"The run ID"`
	Lines []string `json:"lines" desc:"The last lines of the log, oldest first"`
}

// Human prints the lines.
func (r RunTailResponse) Human(w io.Writer) {
	for _, l := range r.Lines {
		fmt.Fprintln(w, l)
	}
}

func runRunTail(ctx *Ctx, req any) (any, error) {
	r := req.(RunTailRequest)
	limit := r.Limit
	if limit == 0 {
		limit = runTailDefaultLimit
	}
	if limit < 0 || limit > runTailMaxLimit {
		return nil, Usagef("limit %d is out of range: give 1 to %d", r.Limit, runTailMaxLimit)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	var j run.Journal
	if ref := strings.TrimSpace(r.Run); ref != "" {
		if j, err = readRun(rs, ref); err != nil {
			return nil, err
		}
	} else if j, err = latestRun(ctx, rs); err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(runLogDir(common), j.ID+".log"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NotFoundf("run %s has no log: only a run the daemon started has one", j.ID)
		}
		return nil, Errorf("read log: %v", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, Errorf("read log: %v", err)
	}
	off := max(fi.Size()-runTailWindow, 0)
	buf := make([]byte, fi.Size()-off)
	if _, err := f.ReadAt(buf, off); err != nil && err != io.EOF {
		return nil, Errorf("read log: %v", err)
	}
	text := string(buf)
	if off > 0 { // the window starts mid-line: drop the partial one
		if i := strings.IndexByte(text, '\n'); i >= 0 {
			text = text[i+1:]
		} else {
			text = ""
		}
	}
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	if len(lines) > limit {
		lines = lines[len(lines)-limit:]
	}
	if lines == nil {
		lines = []string{}
	}
	return RunTailResponse{Run: j.ID, Lines: lines}, nil
}
