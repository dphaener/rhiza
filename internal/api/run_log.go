package api

// This file holds run.log: a headless run's driver log, paged. The daemon
// starts a run outside Rex with its output in a file under
// <git common dir>/rhiza/logs: <change set>.log for the start,
// <run>.log for a resume. The operation reads those files only.

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.log",
		Summary:  "Page a headless run's driver log a page of lines at a time, or open it in a Rex layer with --layer; exits 2 for a run with no log",
		Request:  reflect.TypeFor[RunLogRequest](),
		Response: reflect.TypeFor[RunLogResponse](),
		Run:      runRunLog,
	})
}

const (
	runLogDefaultLimit = 50
	runLogMaxLimit     = 1000
)

// RunLogRequest names the run and the page.
type RunLogRequest struct {
	Run    string `json:"run,omitempty" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run"`
	Limit  int    `json:"limit,omitempty" desc:"The most lines to return, 1 to 1000; default 50"`
	Cursor string `json:"cursor,omitempty" desc:"Where to go on: the next_cursor of the page before"`
	All    bool   `json:"all,omitempty" flag:"all" desc:"Print the whole log, ignoring --limit and --cursor"`
	Layer  bool   `json:"layer,omitempty" flag:"layer" desc:"Open the log in a Rex layer, in the session the caller runs in, instead of printing it"`
}

// RunLogResponse is one page of the log, or says the layer was opened.
type RunLogResponse struct {
	Run        string   `json:"run" desc:"The run ID"`
	Lines      []string `json:"lines" desc:"The log's lines of this page, oldest first"`
	NextCursor string   `json:"next_cursor,omitempty" desc:"The cursor of the next page; absent when no more remain"`
	Opened     bool     `json:"opened,omitempty" desc:"True when --layer opened the log in a Rex layer"`
}

// Human prints the lines, then the command for the next page.
func (r RunLogResponse) Human(w io.Writer) {
	if r.Opened {
		fmt.Fprintf(w, "Opened the log of run %s in a Rex layer\n", r.Run)
		return
	}
	for _, l := range r.Lines {
		fmt.Fprintln(w, l)
	}
	if r.NextCursor != "" {
		fmt.Fprintf(w, "\nnext: rhiza run log %s --cursor %s\n", r.Run, r.NextCursor)
	}
}

// runLogLayerLua opens the log in a layer: a shell running `<self> run log
// <run> --all` through less, at the end. Arguments: self, run (Lua string
// literals).
const runLogLayerLua = `
local res, err = rex.session.new_layer{
  bounds = { 0.05, 0.05, 0.9, 0.9 },
  layout = rex.layout.block{
    flavor = "com.superlogical.terminal.shell",
    label = "run log",
    options = { command = { "sh", "-c", 'exec 2>&1; "$0" run log "$1" --all | less +G', %[1]s, %[2]s } },
  },
  focus = true,
}
if err ~= nil then error(err, 0) end
return res
`

func runRunLog(ctx *Ctx, req any) (any, error) {
	r := req.(RunLogRequest)
	limit := r.Limit
	if limit == 0 {
		limit = runLogDefaultLimit
	}
	if limit < 0 || limit > runLogMaxLimit {
		return nil, Usagef("limit %d is out of range: give 1 to %d", r.Limit, runLogMaxLimit)
	}
	start := 0
	if r.Cursor != "" {
		n, err := strconv.Atoi(r.Cursor)
		if err != nil || n < 0 {
			return nil, Usagef("bad cursor %q: give a next_cursor from run log", r.Cursor)
		}
		start = n
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
	st := j.State()
	var paths []string
	for _, name := range []string{st.Change, j.ID} {
		if name == "" {
			continue
		}
		p := filepath.Join(common, "rhiza", "logs", name+".log")
		if _, err := os.Stat(p); err == nil {
			paths = append(paths, p)
		}
	}
	if len(paths) == 0 {
		return nil, Usagef("run %s has no log: only a run the daemon started has one", j.ID).WithNext("rhiza run status " + j.ID)
	}
	resp := RunLogResponse{Run: j.ID, Lines: []string{}}
	if r.Layer {
		if ctx.Exec == nil || ctx.Self == "" || ctx.RexSession == "" {
			return nil, Errorf("cannot open a layer: no Rex session")
		}
		if ctx.DryRun {
			return resp, nil
		}
		lua := fmt.Sprintf(runLogLayerLua, strconv.Quote(ctx.Self), strconv.Quote(j.ID))
		if _, err := ctx.Exec("rex", []string{"do", "-s", ctx.RexSession, "-e", lua}, ctx.Root, nil); err != nil {
			return nil, Errorf("rex could not open the layer: %v", err)
		}
		resp.Opened = true
		return resp, nil
	}
	n := 0
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, Errorf("read log: %v", err)
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
		for sc.Scan() {
			if r.All || (n >= start && len(resp.Lines) < limit) {
				resp.Lines = append(resp.Lines, sc.Text())
			} else if n >= start+limit {
				resp.NextCursor = strconv.Itoa(start + limit)
				f.Close()
				return resp, nil
			}
			n++
		}
		err = sc.Err()
		f.Close()
		if err != nil {
			return nil, Errorf("read log: %v", err)
		}
	}
	return resp, nil
}
