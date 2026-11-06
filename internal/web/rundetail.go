package web

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// A run's page, /runs/<id>: its steps from the journal as a timeline with
// each step's outcome and time, the agent sessions of each step, and for a
// review run its checks' findings. The body is a live section that
// re-renders when a run's journal moves. The headless driver log (run.log)
// is shown as a log tail beneath (runlog.go).

type runStepItem struct {
	Step    string `json:"step"`
	Outcome string `json:"outcome"`
	Detail  string `json:"detail"`
	At      string `json:"at"`
}

type runFindingItem struct {
	Severity string `json:"severity"`
	Place    string `json:"place"`
	Reason   string `json:"reason"`
	Verdict  string `json:"verdict"`
}

// runInfo is the part of a run.status response the page uses.
type runInfo struct {
	Title     string        `json:"title"`
	Workflow  string        `json:"workflow"`
	State     string        `json:"state"`
	Current   string        `json:"current"`
	Detail    string        `json:"detail"`
	StartedAt string        `json:"started_at"`
	Steps     []runStepItem `json:"steps"`
	Sessions  []struct {
		Step     string `json:"step"`
		NoTokens string `json:"no_tokens"`
		Tokens   *struct {
			Session struct {
				Input         int64 `json:"input"`
				Output        int64 `json:"output"`
				CacheCreation int64 `json:"cache_creation"`
				CacheRead     int64 `json:"cache_read"`
			} `json:"session"`
		} `json:"tokens"`
	} `json:"sessions"`
	Checks []struct {
		Check   string `json:"check"`
		Outcome string `json:"outcome"`
		Rounds  []struct {
			Findings []runFindingItem `json:"findings"`
		} `json:"rounds"`
	} `json:"checks"`
}

type reviewCheck struct {
	Check    string
	Pill     Pill
	Findings []runFindingItem
}

type runDetailView struct {
	Empty    *Empty
	NoSteps  *Empty
	Pill     Pill
	Workflow string
	Steps    []Step
	Review   []reviewCheck
}

type runPageData struct {
	Kicker, Title, Frag, LogFrag string
	View                         runDetailView
	Log                          runLogView
}

func stepPill(state string) Pill {
	switch state {
	case "done":
		return Pill{"done", "done"}
	case "failed":
		return Pill{"failed", "failed"}
	case "running":
		return Pill{"running", "running"}
	case "waiting":
		return Pill{"waiting", "waiting on you"}
	}
	return Pill{}
}

func fmtDuration(d time.Duration) string {
	switch {
	case d < 0:
		return ""
	case d < time.Minute:
		return strconv.Itoa(int(d.Seconds())) + "s"
	case d < time.Hour:
		return strconv.Itoa(int(d.Minutes())) + "m " + strconv.Itoa(int(d.Seconds())%60) + "s"
	}
	return strconv.Itoa(int(d.Hours())) + "h " + strconv.Itoa(int(d.Minutes())%60) + "m"
}

// buildRunDetail lays a run's journal out as a timeline.
func buildRunDetail(ri runInfo, href string) runDetailView {
	v := runDetailView{Pill: runPill(ri.State), Workflow: ri.Workflow,
		NoSteps: &Empty{Title: "No steps yet", Text: "Steps appear here as the run moves."}}
	sessions := func(step string) []Session {
		var out []Session
		n := 0
		for _, s := range ri.Sessions {
			if s.Step != step {
				continue
			}
			n++
			se := Session{Href: href, Who: "Session " + strconv.Itoa(n)}
			switch {
			case s.Tokens != nil:
				t := s.Tokens.Session
				se.Usage = "finished · " + strconv.FormatInt(t.Input+t.Output+t.CacheCreation+t.CacheRead, 10) + " tokens"
			case strings.Contains(s.NoTokens, "running"):
				se.Usage = "running"
			default:
				se.Usage = "finished"
			}
			out = append(out, se)
		}
		return out
	}
	prev, _ := time.Parse(time.RFC3339, ri.StartedAt)
	for _, s := range ri.Steps {
		st := Step{Name: s.Step, State: "done", Note: s.Detail, Sessions: sessions(s.Step)}
		if s.Outcome != "pass" {
			st.State = "failed"
		}
		st.Pill = stepPill(st.State)
		if st.State == "failed" {
			st.Pill.Text = s.Outcome
		}
		if t, err := time.Parse(time.RFC3339, s.At); err == nil {
			st.Duration = fmtDuration(t.Sub(prev))
			prev = t
			when := t.UTC().Format("Jan 2 15:04")
			if st.Note == "" {
				st.Note = when
			} else {
				st.Note = when + " · " + st.Note
			}
		}
		v.Steps = append(v.Steps, st)
	}
	if ri.Current != "" {
		st := Step{Name: ri.Current, State: "running", Sessions: sessions(ri.Current)}
		if ri.State == "blocked" {
			st.State = "waiting"
		} else if ri.State == "stopped" {
			st.State = "pending"
		}
		st.Pill = stepPill(st.State)
		v.Steps = append(v.Steps, st)
	}
	for _, c := range ri.Checks {
		rc := reviewCheck{Check: c.Check}
		switch c.Outcome {
		case "passed":
			rc.Pill = Pill{"done", "passed"}
		case "failed":
			rc.Pill = Pill{"failed", "failed"}
		default:
			rc.Pill = Pill{"running", "open"}
		}
		if n := len(c.Rounds); n > 0 {
			rc.Findings = c.Rounds[n-1].Findings
		}
		v.Review = append(v.Review, rc)
	}
	return v
}

// runDetail reads run id; the status is the HTTP status to answer with.
func runDetail(r *http.Request, id string) (runPageData, int) {
	var ri runInfo
	if err := apiGet(r, "/api/run.status?run="+url.QueryEscape(id), &ri); err != nil {
		if _, ok := err.(*notFoundError); ok {
			return runPageData{}, http.StatusNotFound
		}
		return runPageData{Kicker: "Run", Title: "Run", View: runDetailView{Empty: &Empty{Title: "This run is unavailable", Text: err.Error()}}}, http.StatusOK
	}
	title := ri.Title
	if title == "" {
		title = "Untitled change set"
	}
	return runPageData{
		Kicker: "Run", Title: title,
		Frag: "run-detail?run=" + url.QueryEscape(id), LogFrag: "run-log?run=" + url.QueryEscape(id),
		View: buildRunDetail(ri, "/runs/"+id),
	}, http.StatusOK
}

func servRun(w http.ResponseWriter, r *http.Request) {
	d, code := runDetail(r, r.PathValue("id"))
	if code == http.StatusNotFound {
		http.NotFound(w, r)
		return
	}
	if d.LogFrag != "" {
		d.Log = runLogData(r, r.PathValue("id"))
	}
	if d.Frag == "" {
		d.Frag = "run-detail?run=" + url.QueryEscape(r.PathValue("id"))
		d.LogFrag = "run-log?run=" + url.QueryEscape(r.PathValue("id"))
	}
	render(w, "run", page{Title: d.Title, Nav: "runs", Data: d})
}

func init() {
	RegisterFragment("run-detail", &Fragment{
		Events:   []string{EventRun},
		Template: template.Must(template.New("run-detail").ParseFS(templateFS, "templates/partials.html", "templates/run-detail.html")).Lookup("run-detail"),
		Data: func(r *http.Request) any {
			d, code := runDetail(r, r.URL.Query().Get("run"))
			if code == http.StatusNotFound {
				return runDetailView{Empty: &Empty{Title: "No such run"}}
			}
			return d.View
		},
	})
}
