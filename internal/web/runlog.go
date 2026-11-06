package web

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// A headless run's log tail on its page: the last lines of the driver log in
// the always-dark terminal block, a live section that re-renders when a log
// grows. A run without a log says so.

const logTailLines = 100

type runLogView struct {
	Empty *Empty
	Term  Term
}

func runLogData(r *http.Request, id string) runLogView {
	var resp struct {
		Lines []string `json:"lines"`
	}
	if err := apiGet(r, "/api/run.tail?limit="+strconv.Itoa(logTailLines)+"&run="+url.QueryEscape(id), &resp); err != nil {
		if _, ok := err.(*notFoundError); ok {
			return runLogView{Empty: &Empty{Title: "No log", Text: "Only a run the daemon started headless has a driver log."}}
		}
		return runLogView{Empty: &Empty{Title: "The log is unavailable", Text: err.Error()}}
	}
	t := Term{Label: "Run log", Head: "log · tail", Live: true}
	for _, l := range resp.Lines {
		t.Lines = append(t.Lines, TermLine{Text: strings.TrimRight(l, "\r")})
	}
	if len(t.Lines) == 0 {
		t.Lines = []TermLine{{Text: "(empty)"}}
	}
	return runLogView{Term: t}
}

func init() {
	RegisterFragment("run-log", &Fragment{
		Events:   []string{EventLog},
		Template: template.Must(template.New("run-log").ParseFS(templateFS, "templates/partials.html", "templates/run-log.html")).Lookup("run-log"),
		Data:     func(r *http.Request) any { return runLogData(r, r.URL.Query().Get("run")) },
	})
}
