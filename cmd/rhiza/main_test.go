package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/guide"
	"github.com/dphaener/rhiza/internal/rootguard"
	runpkg "github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/rogpeppe/go-internal/txtar"
)

func TestMain(m *testing.M) {
	// Started as rhiza (exec rhiza, a fake bin/ callback, ctx.Self), the
	// test binary is a rhiza process. RunMain gets no rhiza entry, so a
	// script's bare rhiza line reaches TestScripts' in-process command.
	if filepath.Base(os.Args[0]) == "rhiza" {
		os.Exit(rhizaProcess())
	}
	os.Exit(rootguard.Main(func() int {
		return testscript.RunMain(m, map[string]func() int{"rhiza-process": rhizaProcess})
	}))
}

// rhizaProcess is the test binary run as a rhiza process. Under a script
// (RHIZA_TEST_SCRIPT set by setup) it exits 2 unless the script is on
// processScripts: every other script runs rhiza in-process. A refusal is
// also written to the script's refusal marker ($WORK/processRefusedMarker),
// so the script fails at once (refusedCheck) instead of waiting on a
// caller that ignored the exit status.
func rhizaProcess() int {
	if name := os.Getenv("RHIZA_TEST_SCRIPT"); name != "" {
		if _, ok := processScripts[name]; !ok {
			msg := fmt.Sprintf("rhiza: script %s started rhiza as a process but is not on processScripts, the process list in cmd/rhiza/main_test.go: run its rhiza lines in-process (a bare rhiza line), or add it to the list with its reason\n", name)
			fmt.Fprint(os.Stderr, msg)
			if work := os.Getenv("WORK"); work != "" {
				if f, err := os.OpenFile(filepath.Join(work, processRefusedMarker), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
					fmt.Fprint(f, msg)
					f.Close()
				}
			}
			return 2
		}
	}
	if log := os.Getenv("RHIZA_TEST_PROCESS_LOG"); log != "" {
		// TestInProcessCommand counts the rhiza processes a script starts.
		if f, err := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			fmt.Fprintln(f, strings.Join(os.Args[1:], " "))
			f.Close()
		}
	}
	return profiled(os.Getenv("RHIZA_CPUPROFILE"), func() int {
		wd, _ := os.Getwd()
		// A real rhiza process: it adopts and reaps its own children.
		return run(os.Args[1:], wd, selfPath(), os.Stdin, os.Stdout, os.Stderr, os.Environ(), processChildren)
	})
}

// -update rewrites a script's archive file when a `cmp stdout <file>` fails,
// e.g. describe-golden.txtar's golden after a request or response type
// changes: go test ./cmd/rhiza -run TestScripts/describe-golden -update
var update = flag.Bool("update", false, "rewrite the archive files that cmp compares against")

func TestScripts(t *testing.T) {
	// Longest first by script-timings.tsv (order_test.go). testscript reads
	// Files only when Dir is empty.
	files, err := scriptFiles("testdata/script", os.Getenv("RHIZA_TEST_SHARD"), scriptTimingsFile)
	if err != nil {
		t.Fatal(err)
	}
	p := scriptParams(t, files)
	p.UpdateScripts = *update
	testscript.Run(t, p)
}

// scriptParams is TestScripts' testscript.Params over files: setup, the
// script commands and the in-process rhiza. It links rhiza on PATH first,
// so setup finds it for RHIZA_TEST_BIN.
func scriptParams(t *testing.T, files []string) testscript.Params {
	self := rhizaLink(t)
	paths := map[string]string{}
	for _, f := range files {
		paths[scriptName(f)] = f
	}
	return testscript.Params{
		Files: files,
		Setup: func(env *testscript.Env) error {
			if err := setup(env); err != nil {
				return err
			}
			return keepScriptEnv(env, paths[scriptWorkName(env.WorkDir)])
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"rhiza":   func(ts *testscript.TestScript, neg bool, args []string) { rhizaInProcess(ts, neg, args, self) },
			"fold-ok": foldOK, "quickstart-doc": quickstartDoc, "chtimes": chtimes, "health-server": healthServer, "run-lock": runLock, "file-server": fileServer, "root-daemon-lock": rootDaemonLock, "edit-sessions": editSessions,
		},
	}
}

var (
	rhizaLinkOnce sync.Once
	rhizaLinkPath string
	rhizaLinkErr  error
)

// rhizaLink hard-links, once per test binary, a rhiza beside the
// rhiza-process RunMain put in its bin directory on PATH, and returns the
// link's path: exec rhiza, fake bin/ callbacks and ctx.Self start the test
// binary under the name rhiza, which TestMain runs as rhiza. A hard link,
// not a symlink or a copy: os.Args[0] must name rhiza, and a copy would be
// a second binary to keep in step.
func rhizaLink(t *testing.T) string {
	rhizaLinkOnce.Do(func() {
		bin, err := exec.LookPath("rhiza-process")
		if err != nil {
			rhizaLinkErr = err
			return
		}
		link := filepath.Join(filepath.Dir(bin), "rhiza")
		if err := os.Link(bin, link); err != nil && !errors.Is(err, fs.ErrExist) {
			rhizaLinkErr = err
			return
		}
		rhizaLinkPath = link
	})
	if rhizaLinkErr != nil {
		t.Fatalf("link rhiza beside rhiza-process: %v", rhizaLinkErr)
	}
	return rhizaLinkPath
}

// scriptWorkName is the script name in a testscript work dir
// (.../script-<name>).
func scriptWorkName(workDir string) string {
	return strings.TrimPrefix(filepath.Base(workDir), "script-")
}

// processScripts is the one list of scripts allowed to start rhiza as a
// process (exec rhiza, a fake bin/ command that calls rhiza, rhiza with &,
// a stdin or ttyin feed, a command that starts ctx.Self, or a line whose
// run or daemon lock a later line must find held by an exited process),
// each with its reason. Every other script's rhiza lines run in-process,
// and a rhiza process it starts exits 2 naming this list.
var processScripts = map[string]string{
	// The fixtures of script_inprocess_test.go, compared with exec rhiza.
	"process-list-fixture": "TestProcessScripts' fixture on the list",
	"in-process-fixture":   "TestInProcessCommand compares its in-process lines with exec rhiza",

	"anchor-fields":                        "starts rhiza through sh -c",
	"apply":                                "a stdin feed",
	"arbiter-agent-tabs":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"arbiter-answer":                       "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"arbiter-question":                     "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"arbiter-step":                         "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"autonomous-gate-rule":                 "starts rhiza through sh -c",
	"autonomous-permissions":               "a fake bin/ rex that runs its tab's rhiza command",
	"autonomous-review":                    "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"autonomy-accept-rule":                 "starts rhiza through sh -c",
	"brief-files":                          "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"bugfix-workflow-autonomous":           "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"build-check":                          "starts rhiza through sh -c",
	"build-evidence-commit":                "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"build-leaf-fails":                     "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"build-leaves-status":                  "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"build-tabs":                           "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"build-verifier-checkout":              "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"build-verifier-findings":              "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"build-workflow":                       "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"carve-guards":                         "starts rhiza through sh -c",
	"change-new-workflow-flag":             "starts rhiza through sh -c",
	"change-request-kind":                  "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"change-set-ref":                       "starts rhiza through sh -c",
	"change-show-canonical-text":           "starts rhiza through sh -c",
	"change-show-serves":                   "starts rhiza through sh -c",
	"change-workflow-pin":                  "starts rhiza through sh -c",
	"check":                                "starts rhiza through sh -c",
	"check-evidence-actor":                 "starts rhiza through sh -c",
	"check-git":                            "starts rhiza through sh -c",
	"check-move":                           "starts rhiza through sh -c",
	"ci-install":                           "starts rhiza through sh -c",
	"custom-workflows-doc":                 "quickstart-doc runs the doc's rhiza commands as processes",
	"daemon-continue-slot":                 "starts rhiza through sh -c",
	"daemon-headless-runs":                 "rhiza daemon starts its run as a detached ctx.Self process",
	"daemon-lock":                          "starts rhiza through sh -c; rhiza with &",
	"daemon-milestone-request":             "starts rhiza through sh -c",
	"daemon-no-root":                       "starts rhiza through sh -c",
	"daemon-notify-hook":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; rhiza with &",
	"daemon-plans-proposals":               "starts rhiza through sh -c",
	"daemon-projects":                      "rhiza with &",
	"daemon-resume-landed":                 "starts rhiza through sh -c",
	"daemon-settings":                      "starts rhiza through sh -c",
	"daemon-start-stop":                    "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"daemon-version-in-health":             "starts rhiza through sh -c; rhiza with &",
	"describe":                             "a stdin feed",
	"feed":                                 "starts rhiza through sh -c",
	"feed-follow":                          "starts rhiza through sh -c; rhiza with &",
	"findings-concurrent":                  "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza",
	"findings-concurrent-rounds":           "starts rhiza through sh -c",
	"guide":                                "starts rhiza through sh -c",
	"guide-install":                        "starts rhiza through sh -c",
	"headless-run-agent-tabs":              "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"headless-run-blocked-notice":          "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"hook-commands":                        "rhiza with &",
	"import-records-built":                 "starts rhiza through sh -c",
	"inbox-advisor":                        "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"inbox-api-error":                      "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; rhiza with &",
	"inbox-blocker-advice":                 "starts rhiza through sh -c",
	"inbox-disk-low":                       "starts rhiza through sh -c",
	"inbox-gate-questions":                 "starts rhiza through sh -c",
	"inbox-landed-plan-run":                "starts rhiza through sh -c and its mkchanges.sh helper",
	"inbox-landed-resumable":               "starts rhiza through sh -c and its mkchanges.sh helper",
	"inbox-list":                           "starts rhiza through sh -c",
	"inbox-load-hold":                      "starts rhiza through sh -c",
	"inbox-main-stale":                     "starts rhiza through sh -c",
	"inbox-open":                           "starts rhiza through sh -c",
	"inbox-owner-run-accept":               "starts rhiza through sh -c",
	"inbox-permission":                     "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"inbox-projects":                       "starts rhiza through sh -c",
	"inbox-question-advisor":               "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; a fake bin/ gate that calls rhiza",
	"inbox-signoff":                        "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"inbox-stale":                          "starts rhiza through sh -c and its mkchanges.sh helper",
	"inbox-stalled":                        "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; rhiza with &",
	"inbox-uncovered-proposal":             "starts rhiza through sh -c and its mkcs.sh helper",
	"init":                                 "a stdin feed",
	"leaf-cycle":                           "starts rhiza through sh -c",
	"log-paged":                            "starts rhiza through sh -c",
	"ls":                                   "starts rhiza through sh -c",
	"machine-trust-bypass":                 "starts rhiza through sh -c",
	"machine-trust-unwritten":              "starts rhiza through sh -c",
	"old-folder":                           "starts rhiza through sh -c",
	"owner-run-gate-accept":                "starts rhiza through sh -c",
	"plan-milestone-backlog":               "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-milestone-check":                 "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-milestone-interview-waiting":     "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"plan-milestone-number":                "starts rhiza through sh -c",
	"plan-requests":                        "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"plan-review-visit-round":              "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-review-workflow":                 "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-run-milestone-env":               "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"plan-workflow":                        "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-autonomous":             "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-autonomous-left":        "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-ci-flake":               "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-ci-flake-red":           "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-ci-infra-timeout":       "a workflow step runs rhiza",
	"plan-workflow-ci-infra-wait":          "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; a fake bin/ rhiza that wraps the real one; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-merge-main-moved":       "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-merge-main-moved-often": "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-open":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"plan-workflow-pr":                     "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"printout-counts-used-tokens":          "a fake bin/ rex that runs its tab's rhiza command",
	"printout-gap-names-reason":            "a fake bin/ rex that runs its tab's rhiza command",
	"project-add-path":                     "starts rhiza through sh -c",
	"quickstart":                           "quickstart-doc runs the doc's rhiza commands as processes",
	"render":                               "starts rhiza through sh -c",
	"render-generate":                      "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"reopen-base":                          "starts rhiza through sh -c",
	"request-ref":                          "starts rhiza through sh -c",
	"review-check-worktree":                "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"review-checks-step":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"review-code-check":                    "starts rhiza through sh -c; a workflow step runs rhiza",
	"review-code-check-override":           "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"review-findings":                      "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza",
	"review-signoff":                       "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"review-signoff-autonomous":            "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"review-signoff-changes":               "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"review-summary":                       "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"review-workflow":                      "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"roadmap-paged":                        "starts rhiza through sh -c",
	"roadmap-projects":                     "starts rhiza through sh -c",
	"root-init":                            "starts rhiza through sh -c",
	"run-abandon":                          "starts rhiza through sh -c",
	"run-abandon-failed":                   "starts rhiza through sh -c",
	"run-abandon-shared":                   "starts rhiza through sh -c",
	"run-abandon-unmerged":                 "starts rhiza through sh -c",
	"run-abandon-unpushed":                 "starts rhiza through sh -c",
	"run-agent-done":                       "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"run-agent-step":                       "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"run-agent-tabs":                       "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"run-amend":                            "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"run-arbitrate-asked":                  "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-autonomous":                       "a fake bin/ rex that runs its tab's rhiza command",
	"run-blocked":                          "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"run-cleanup-steps":                    "starts rhiza through sh -c",
	"run-command-step":                     "starts rhiza through sh -c",
	"run-continue-queue":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza daemon starts its run as a detached ctx.Self process",
	"run-end-removes-leftovers":            "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza",
	"run-gate":                             "starts rhiza through sh -c",
	"run-handback":                         "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"run-log":                              "starts rhiza through sh -c; rhiza daemon starts its run as a detached ctx.Self process",
	"run-ls":                               "starts rhiza through sh -c",
	"run-max-visits":                       "starts rhiza through sh -c",
	"run-nested-workflow-session":          "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-one-live":                         "starts rhiza through sh -c",
	"run-records-dispatch":                 "starts rhiza through sh -c",
	"run-resume":                           "starts rhiza through sh -c",
	"run-resume-agent-session":             "starts rhiza through sh -c",
	"run-resume-landed-past-land":          "starts rhiza through sh -c",
	"run-resume-landed-plan":               "starts rhiza through sh -c",
	"run-resume-lost-fixer":                "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"run-resume-lost-session":              "starts rhiza through sh -c; rhiza with &",
	"run-resume-retry":                     "starts rhiza through sh -c",
	"run-resume-unaccepted":                "starts rhiza through sh -c",
	"run-review-check":                     "starts rhiza through sh -c",
	"run-review-check-handles":             "starts rhiza through sh -c",
	"run-review-check-landed":              "starts rhiza through sh -c",
	"run-review-journal":                   "starts rhiza through sh -c",
	"run-send-back":                        "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-change":                      "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-change-fail":                 "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-change-rebase":               "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-ci-fix":                      "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c and its setup.sh helper; a workflow step runs rhiza",
	"run-ship-ci-lint-fix":                 "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c and its setup.sh helper; a workflow step runs rhiza",
	"run-ship-conflict-fix":                "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c and its setup.sh helper; a workflow step runs rhiza",
	"run-ship-conflict-fix-more":           "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c and its setup.sh helper; rhiza with &; a workflow step runs rhiza",
	"run-ship-deploy-check":                "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-deploy-timeout":              "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; a workflow step runs rhiza",
	"run-ship-infra-wait":                  "a fake bin/ rhiza that wraps the real one; a workflow step runs rhiza",
	"run-ship-land-conflict":               "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; a fake bin/ headrhiza that calls rhiza; starts rhiza through sh -c and its setup.sh helper; a workflow step runs rhiza",
	"run-ship-landed-infra-wait":           "a fake bin/ rhiza that wraps the real one; a workflow step runs rhiza",
	"run-ship-main-dispatch":               "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-main-dispatch-fail":          "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-main-dispatch-resume":        "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-merge-draft":                 "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-merge-main-moved":            "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-model-conflict":              "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c and its setup.sh helper; a workflow step runs rhiza",
	"run-ship-needs-leaves":                "a fake bin/ rex that runs its tab's rhiza command",
	"run-ship-rerun-fail":                  "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-rerun-pass":                  "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; a workflow step runs rhiza",
	"run-ship-retry-ci-triage":             "starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-retry-pending":               "starts rhiza through sh -c; a workflow step runs rhiza",
	"run-ship-rspec-failure":               "starts rhiza through sh -c",
	"run-ship-signoff":                     "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-signoff-workflow-owner":           "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"run-spawn":                            "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"run-start-milestone-autonomous":       "starts rhiza through sh -c",
	"run-start-new":                        "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"run-start-new-main-clean":             "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c; a workflow step runs rhiza",
	"run-start-queue":                      "starts rhiza through sh -c; rhiza daemon starts its run as a detached ctx.Self process",
	"run-start-ready-dispatch":             "starts rhiza through sh -c",
	"run-start-skill-missing":              "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"run-status":                           "starts rhiza through sh -c; rhiza with &",
	"run-status-messages":                  "a stdin feed",
	"run-status-stalled":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; rhiza with &",
	"run-status-tokens":                    "a fake bin/ rex that runs its tab's rhiza command",
	"run-tail":                             "starts rhiza through sh -c; rhiza daemon starts its run as a detached ctx.Self process",
	"search":                               "starts rhiza through sh -c",
	"session-done-captures":                "a fake bin/ claude that calls rhiza",
	"session-hook-log-guard":               "a stdin feed",
	"session-hook-message":                 "a stdin feed",
	"session-interview":                    "starts rhiza through sh -c",
	"session-interview-landed":             "starts rhiza through sh -c",
	"session-interview-recorded":           "starts rhiza through sh -c",
	"session-journal":                      "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza through its hook command",
	"session-message":                      "starts rhiza through sh -c; rhiza with &; a stdin feed",
	"session-message-busy":                 "a stdin feed",
	"session-message-idle-fallback":        "a stdin feed",
	"session-message-permission":           "rhiza with &; a stdin feed",
	"session-pane-key":                     "starts rhiza through sh -c",
	"session-run":                          "a fake bin/ claude that calls rhiza; starts rhiza through sh -c; a stdin feed",
	"session-run-blocked-notice":           "a fake bin/ scenario that calls rhiza; starts rhiza through sh -c; rhiza with &",
	"session-run-tokens":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"session-start":                        "starts rhiza through sh -c",
	"set":                                  "starts rhiza through sh -c",
	"ship-change-one-ci":                   "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"ship-check-ci-test-findings":          "starts rhiza through sh -c",
	"ship-check-settles-conflict":          "starts rhiza through sh -c",
	"ship-close-out-forget-retry":          "starts rhiza through sh -c",
	"ship-close-out-log-copies":            "starts rhiza through sh -c",
	"ship-gate-head-check":                 "starts rhiza through sh -c",
	"ship-land-branch-build":               "a fake bin/ buildhead that calls rhiza",
	"show":                                 "starts rhiza through sh -c",
	"skill-check":                          "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"skill-check-roles":                    "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza",
	"slice-errors":                         "starts rhiza through sh -c",
	"source-add":                           "starts rhiza through sh -c",
	"source-raw-session":                   "starts rhiza through sh -c",
	"stacked":                              "starts rhiza through sh -c",
	"stacked-cycle":                        "starts rhiza through sh -c",
	"stacked-revised":                      "starts rhiza through sh -c",
	"stacked-revised-drop":                 "starts rhiza through sh -c",
	"stage-autonomy":                       "starts rhiza through sh -c",
	"status":                               "starts rhiza through sh -c",
	"text-refs-display":                    "starts rhiza through sh -c",
	"text-refs-engine-writes":              "starts rhiza through sh -c",
	"text-refs-write":                      "starts rhiza through sh -c",
	"tui-launch-usage":                     "starts rhiza through sh -c",
	"tui-open-usage":                       "starts rhiza through sh -c",
	"usage-hold":                           "starts rhiza through sh -c",
	"usage-live-record":                    "a stdin feed",
	"watchdog-health":                      "starts rhiza through sh -c",
	"why-paged":                            "starts rhiza through sh -c",
	"workflow-invoke":                      "a fake bin/ rex that runs its tab's rhiza command; starts rhiza through sh -c",
	"workflow-new":                         "starts rhiza through sh -c",
	"workflow-personal-override":           "a fake bin/ rex that runs its tab's rhiza command; a fake bin/ claude that calls rhiza; starts rhiza through sh -c",
	"workflow-stage-autonomy":              "starts rhiza through sh -c",
}

// notConverted is the reason a script kept on processScripts until its
// rhiza lines ran in-process. No entry may keep it: unconverted fails one.
const notConverted = "not converted yet"

// unconverted returns the scripts on list whose reason is notConverted,
// sorted.
func unconverted(list map[string]string) []string {
	var names []string
	for name, reason := range list {
		if reason == notConverted {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// scriptEnvKey is the key of a script's environment names in env.Values.
type scriptEnvKey struct{}

// keepScriptEnv keeps, for the in-process rhiza, the names of the script's
// environment: setup's (env.Vars) and those the script's own env lines set
// (file, the script's txtar; none when it is unknown).
func keepScriptEnv(env *testscript.Env, file string) error {
	var names []string
	seen := map[string]bool{}
	add := func(n string) {
		if n != "" && !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	for _, kv := range env.Vars {
		n, _, _ := strings.Cut(kv, "=")
		add(n)
	}
	if file != "" {
		a, err := txtar.ParseFile(file)
		if err != nil {
			return err
		}
		for _, n := range envLineNames(string(a.Comment)) {
			add(n)
		}
	}
	env.Values[scriptEnvKey{}] = names
	return nil
}

// envLineNames is every NAME a script's `env NAME=...` lines set, in
// order, conditions and ! prefixes skipped.
func envLineNames(script string) []string {
	var names []string
	for _, line := range strings.Split(script, "\n") {
		f := strings.Fields(line)
		for len(f) > 0 && (f[0] == "!" || strings.HasPrefix(f[0], "[")) {
			f = f[1:]
		}
		if len(f) < 2 || f[0] != "env" {
			continue
		}
		for _, a := range f[1:] {
			if strings.HasPrefix(a, "#") {
				break
			}
			a = strings.Trim(a, "'")
			if n, _, ok := strings.Cut(a, "="); ok && n != "" {
				names = append(names, n)
			}
		}
	}
	return names
}

// scriptEnv is the environment the in-process rhiza gets: every name
// keepScriptEnv kept, with its current value (ts.Getenv). An empty value
// is kept as NAME=, as an exec'd rhiza gets it after `env NAME=`.
func scriptEnv(ts *testscript.TestScript) []string {
	names, _ := ts.Value(scriptEnvKey{}).([]string)
	env := make([]string, 0, len(names))
	for _, n := range names {
		env = append(env, n+"="+ts.Getenv(n))
	}
	return env
}

// processRefusedMarker is the file under a script's $WORK where
// rhizaProcess records each refused rhiza process.
const processRefusedMarker = ".rhiza-process-refused"

// refusedEvery is how often an in-process rhiza line looks for the
// refusal marker while run() works.
const refusedEvery = 100 * time.Millisecond

// processRefused returns what rhizaProcess recorded under work ("" when it
// refused no rhiza process there).
func processRefused(work string) string {
	b, err := os.ReadFile(filepath.Join(work, processRefusedMarker))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// refusedCheck fails the script when one of its rhiza processes was
// refused: an off-list script fails at once, never waiting on a caller (a
// fake bin/ callback, a hook) that ignored the exit status 2.
func refusedCheck(ts *testscript.TestScript) {
	if msg := processRefused(ts.Getenv("WORK")); msg != "" {
		ts.Fatalf("%s", msg)
	}
}

// rhizaInProcess is the script command `rhiza args...`: it runs run() in
// the test process with the script's working directory, environment and
// streams, self as its own path, and fails or passes as `exec rhiza` would.
// A refused rhiza process of the script (refusedCheck), before the line or
// while run() works, fails the line at once; run() is then left to finish
// on its own, writing nowhere.
func rhizaInProcess(ts *testscript.TestScript, neg bool, args []string, self string) {
	if len(args) > 0 && strings.HasPrefix(args[len(args)-1], "&") {
		ts.Fatalf("rhiza: an in-process rhiza cannot run in the background; use exec rhiza ... & and put the script on processScripts")
	}
	refusedCheck(ts)
	var mu sync.Mutex
	stdout := &lockedWriter{mu: &mu, w: ts.Stdout()}
	stderr := &lockedWriter{mu: &mu, w: ts.Stderr()}
	done := make(chan int, 1)
	go func() {
		done <- run(args, ts.MkAbs("."), self, strings.NewReader(""), stdout, stderr, scriptEnv(ts), ownChildren{})
	}()
	tick := time.NewTicker(refusedEvery)
	defer tick.Stop()
	work := ts.Getenv("WORK")
	var code int
wait:
	for {
		select {
		case code = <-done:
			break wait
		case <-tick.C:
			if msg := processRefused(work); msg != "" {
				mu.Lock()
				stdout.w, stderr.w = io.Discard, io.Discard
				mu.Unlock()
				ts.Fatalf("%s", msg)
			}
		}
	}
	refusedCheck(ts)
	switch {
	case code == 0 && neg:
		ts.Fatalf("unexpected command success")
	case code != 0:
		ts.Logf("[exit status %d]", code)
		if !neg {
			ts.Fatalf("unexpected command failure")
		}
	}
}

// lockedWriter serialises writes to w: run() may write from more than one
// goroutine, and ts.Stdout() is a strings.Builder.
type lockedWriter struct {
	mu *sync.Mutex
	w  io.Writer
}

func (l *lockedWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(p)
}

// shardNames returns the files shard (RHIZA_TEST_SHARD, "i/n", 1-based)
// selects, split by the seconds in times so the n shards take about as long:
// each file, longest first, goes to the shard with the least total so far,
// ties to the lower shard. A file times does not list weighs the median of
// the timed files' seconds (1 when none is timed), for placement only. The shards
// are disjoint and together cover every file; shard i comes back in
// longestFirst's order.
func shardNames(files []string, times map[string]float64, shard string) ([]string, error) {
	var i, n int
	var rest string
	k, _ := fmt.Sscanf(shard, "%d/%d%s", &i, &n, &rest)
	if k != 2 || n < 1 || i < 1 || i > n {
		return nil, fmt.Errorf("RHIZA_TEST_SHARD=%q: want i/n with 1 <= i <= n", shard)
	}
	median := medianTime(files, times)
	weight := func(f string) (float64, bool) {
		if s, ok := times[scriptName(f)]; ok {
			return s, true
		}
		return median, false
	}
	// Longest first by weight; on equal weight a timed file goes first, then
	// by name, so the placement is the same on every machine.
	order := slices.Clone(files)
	sort.SliceStable(order, func(a, b int) bool {
		wa, okA := weight(order[a])
		wb, okB := weight(order[b])
		if wa != wb {
			return wa > wb
		}
		if okA != okB {
			return okA
		}
		return scriptName(order[a]) < scriptName(order[b])
	})
	totals := make([]float64, n)
	out := []string{}
	for _, f := range order {
		least := 0
		for s := 1; s < n; s++ {
			if totals[s] < totals[least] {
				least = s
			}
		}
		w, _ := weight(f)
		totals[least] += w
		if least == i-1 {
			out = append(out, f)
		}
	}
	return longestFirst(out, times), nil
}

// medianTime is the median of the seconds times gives files, or 1 when none
// of them is timed, so untimed scripts still spread across the shards.
func medianTime(files []string, times map[string]float64) float64 {
	vs := []float64{}
	for _, f := range files {
		if v, ok := times[scriptName(f)]; ok {
			vs = append(vs, v)
		}
	}
	if len(vs) == 0 {
		return 1
	}
	sort.Float64s(vs)
	m := len(vs) / 2
	if len(vs)%2 == 1 {
		return vs[m]
	}
	return (vs[m-1] + vs[m]) / 2
}

// Every script starts with a fixed clock and ID seed, so its output bytes
// are stable. A script that wants the real clock or random IDs clears them
// with `env RHIZA_TEST_NOW=` and `env RHIZA_TEST_SEED=`.
const (
	testNow  = "2026-10-24T12:00:00Z"
	testSeed = "1"
	// testFreeBytes is the free disk every script sees (1 TiB), so the
	// daemon's minFree guard never reads the host's disk. A script that
	// wants another figure sets RHIZA_TEST_FREE_BYTES itself.
	testFreeBytes = "1099511627776"
	// testClockScale is the scale of every script's clock (internal/run's
	// ScaledClock): a wait's sleep of d takes d/1000 of real time, at
	// least 50 ms.
	testClockScale = "1000"
)

func setup(env *testscript.Env) error {
	if err := checkWorkOutsideRepo(env.WorkDir); err != nil {
		return err
	}
	// The rhiza process entry checks the script against processScripts,
	// and a refusal it records fails the script at its end even when no
	// later rhiza line sees it.
	env.Setenv("RHIZA_TEST_SCRIPT", scriptWorkName(env.WorkDir))
	work := env.WorkDir
	env.Defer(func() {
		if msg := processRefused(work); msg != "" {
			env.T().Fatal(msg)
		}
	})
	env.Setenv("RHIZA_TEST_NOW", testNow)
	env.Setenv("RHIZA_TEST_SEED", testSeed)
	env.Setenv("RHIZA_TEST_FREE_BYTES", testFreeBytes)
	// An idle machine with 1 TiB available: no script reads the host's load
	// or memory. A script that wants other readings sets them itself.
	env.Setenv("RHIZA_TEST_LOAD1", "0")
	env.Setenv("RHIZA_TEST_LOAD5", "0")
	env.Setenv("RHIZA_TEST_MEM_AVAILABLE", testFreeBytes)
	// Every script's waits run on the scaled clock, so a builtin
	// workflow's 30 s poll or 2 min backoff costs milliseconds. A script
	// that wants RealClock clears it with `env RHIZA_TEST_CLOCK_SCALE=`.
	env.Setenv(runpkg.ClockScaleEnv, testClockScale)
	// RHIZA_TEST_BIN is the directory holding the in-process rhiza, so a
	// script can set PATH to it alone and hide the host's programs (e.g.
	// systemctl) from rhiza.
	if bin, err := exec.LookPath("rhiza"); err == nil {
		env.Setenv("RHIZA_TEST_BIN", filepath.Dir(bin))
	}
	// Every script's $WORK is its own Rhiza root with a project named repo
	// at $WORK/repo (written directly, not through rhiza project add).
	// $NOROOT is a folder outside any root, for the outside-a-root cases.
	err := os.WriteFile(filepath.Join(env.WorkDir, "rhiza.toml"),
		[]byte("[projects.repo]\npath = \"repo\"\n\n[daemon]\n\"rhiza.daemon.failureClassifier\" = \"\"\n"), 0o644)
	if err != nil {
		return err
	}
	noroot, err := os.MkdirTemp("", "rhiza-noroot-")
	if err != nil {
		return err
	}
	env.Defer(func() { os.RemoveAll(noroot) })
	env.Setenv("NOROOT", noroot)
	if err := hermeticGit(env); err != nil {
		return err
	}
	return stubGuideCommands(env.Getenv("HOME"))
}

// stubGuideCommands writes, under home's .claude/commands, a command file
// for each guide text (rhiza, then rhiza-<name> for the rest), so every
// script launches its fake agent as if the guide were installed: a session
// whose start command is not installed never launches. Taking the names
// from guide.List keeps every workflow's start skill (the review checks,
// critic and fixers included) covered as topics are added. A script about
// a missing skill removes them.
func stubGuideCommands(home string) error {
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, t := range guide.List() {
		name := "rhiza"
		if t.Name != "rhiza" {
			name += "-" + t.Name
		}
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte("stub\n"), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// checkWorkOutsideRepo fails when dir resolves inside a git work tree (an
// ancestor holds .git) or inside a git dir (a .git path element, or a bare
// repository's HEAD, objects and refs), so no script can write into a real
// repository. A run's GOTMPDIR under the repo puts $WORK there.
func checkWorkOutsideRepo(dir string) error {
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		real = dir
	}
	real, err = filepath.Abs(real)
	if err != nil {
		return err
	}
	for d := real; ; d = filepath.Dir(d) {
		if filepath.Base(d) == ".git" {
			return fmt.Errorf("testscript $WORK %s is inside a git dir: set GOTMPDIR outside the repository", real)
		}
		if _, err := os.Lstat(filepath.Join(d, ".git")); err == nil {
			return fmt.Errorf("testscript $WORK %s is inside the git work tree %s: set GOTMPDIR outside the repository", real, d)
		}
		bare := true
		for _, n := range []string{"HEAD", "objects", "refs"} {
			if _, err := os.Lstat(filepath.Join(d, n)); err != nil {
				bare = false
				break
			}
		}
		if bare {
			return fmt.Errorf("testscript $WORK %s is inside the git dir %s: set GOTMPDIR outside the repository", real, d)
		}
		if filepath.Dir(d) == d {
			return nil
		}
	}
}

// foldOK is the script command `fold-ok [dir]`: it reads the .rhiza folder
// dir (default $WORK/.rhiza) with store.ReadLog, folds it with
// changeset.FoldLog and prints every log and fold problem to stdout, one per
// line. It fails if there is any; `! fold-ok` fails if there is none.
// chtimes sets the modification time of files: chtimes <RFC 3339 time> file...
func chtimes(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) < 2 {
		ts.Fatalf("usage: chtimes time file...")
	}
	t, err := time.Parse(time.RFC3339, args[0])
	if err != nil {
		ts.Fatalf("chtimes: %v", err)
	}
	for _, f := range args[1:] {
		ts.Check(os.Chtimes(ts.MkAbs(f), t, t))
	}
}

// runLock is the script command `run-lock dir id`: it holds the run lock of
// run id in the runs directory dir until the script ends, so a fixture
// journal there counts as a live run.
func runLock(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) != 2 {
		ts.Fatalf("usage: run-lock dir id")
	}
	unlock, ok, err := runpkg.NewStore(ts.MkAbs(args[0])).TryLock(args[1])
	if err != nil || !ok {
		ts.Fatalf("run-lock %s: held=%v err=%v", args[1], !ok, err)
	}
	ts.Defer(unlock)
}

// editSessions is the script command `edit-sessions <edit> [args...]`: it
// applies one hand edit to every session record of the repo in the working
// directory, through the session store's Update, so it runs under the
// store's lock and a concurrent locked writer (session start recording the
// tab, a hook report) cannot drop it. The edits:
//
//	hook-at <time>                   the hook report's time
//	hook-state <state> [kind detail] the hook report's state (and kind, detail)
//	tab <title> <block> <session>    the record's tab
//	rex <state> <time>               a Rex report
//	autonomous                       mark the session autonomous
//	done <time>                      the session's own word that it is done
func editSessions(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) == 0 {
		ts.Fatalf("usage: edit-sessions edit [args...]")
	}
	at := func(v string) time.Time {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			ts.Fatalf("edit-sessions: bad time %q", v)
		}
		return t
	}
	var edit func(r *session.Record)
	a := args[1:]
	switch {
	case args[0] == "hook-at" && len(a) == 1:
		edit = func(r *session.Record) {
			if r.Status.Hook != nil {
				r.Status.Hook.At = at(a[0])
			}
		}
	case args[0] == "hook-state" && (len(a) == 1 || len(a) == 3):
		edit = func(r *session.Record) {
			if r.Status.Hook == nil {
				return
			}
			r.Status.Hook.State = a[0]
			if len(a) == 3 {
				r.Status.Hook.Kind, r.Status.Hook.Detail = a[1], a[2]
			}
		}
	case args[0] == "tab" && len(a) == 3:
		edit = func(r *session.Record) { r.Tab = session.Tab{Title: a[0], Block: a[1], Session: a[2]} }
	case args[0] == "rex" && len(a) == 2:
		edit = func(r *session.Record) { r.Status.Rex = &session.Report{State: a[0], At: at(a[1])} }
	case args[0] == "autonomous" && len(a) == 0:
		edit = func(r *session.Record) { r.Autonomous = true }
	case args[0] == "done" && len(a) == 1:
		edit = func(r *session.Record) { r.Status.Done = &session.Report{State: session.StateDone, At: at(a[0])} }
	default:
		ts.Fatalf("usage: edit-sessions edit [args...]: unknown edit or wrong arguments %q", args)
	}
	st := session.NewStore(session.Dir(ts.MkAbs(".git")))
	recs, err := st.List()
	if err != nil {
		ts.Fatalf("edit-sessions: %v", err)
	}
	for _, r := range recs {
		if _, err := st.Update(r.ID, edit); err != nil {
			ts.Fatalf("edit-sessions: %v", err)
		}
	}
}

// rootDaemonLock is the script command `root-daemon-lock <root>`: it takes
// the root daemon's lock in the test process, so every rhiza the script
// runs sees a root daemon running, until the script ends.
func rootDaemonLock(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) != 1 {
		ts.Fatalf("usage: root-daemon-lock root")
	}
	l, err := daemon.TryLockRoot(ts.MkAbs(args[0]), time.Now())
	if err != nil {
		ts.Fatalf("root-daemon-lock: %v", err)
	}
	ts.Defer(l.Unlock)
}

// healthServer is the script command `health-server <var> <log> <status>...`:
// it starts a fake deploy health server on loopback, in the test process, so
// it outlives each rhiza the script runs, and sets <var> to its URL. The
// server answers the statuses in order, the last one from then on, and
// appends each status it answers to <log>, one per line. It closes when the
// script ends.
func healthServer(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) < 3 {
		ts.Fatalf("usage: health-server var log status...")
	}
	var statuses []int
	for _, a := range args[2:] {
		n, err := strconv.Atoi(a)
		if err != nil {
			ts.Fatalf("health-server: status %q: %v", a, err)
		}
		statuses = append(statuses, n)
	}
	log := ts.MkAbs(args[1])
	var mu sync.Mutex
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		code := statuses[min(hits, len(statuses)-1)]
		hits++
		if f, err := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			fmt.Fprintln(f, code)
			f.Close()
		}
		mu.Unlock()
		w.WriteHeader(code)
	}))
	ts.Defer(srv.Close)
	ts.Setenv(args[0], srv.URL+"/health")
}

// fileServer serves a directory over HTTP: file-server var dir sets var to the
// server's base URL. It closes when the script ends.
func fileServer(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) != 2 {
		ts.Fatalf("usage: file-server var dir")
	}
	srv := httptest.NewServer(http.FileServer(http.Dir(ts.MkAbs(args[1]))))
	ts.Defer(srv.Close)
	ts.Setenv(args[0], srv.URL)
}

func foldOK(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) > 1 {
		ts.Fatalf("usage: fold-ok [dir]")
	}
	dir := ".rhiza"
	if len(args) == 1 {
		dir = args[0]
	}
	dir = ts.MkAbs(dir)
	var problems []string
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		problems = append(problems, fmt.Sprintf("%s is not a folder", dir))
	}
	l := store.ReadLog(os.DirFS(dir))
	_, foldProblems := changeset.FoldLog(l)
	for _, p := range l.Problems {
		problems = append(problems, fmt.Sprintf("log: %s: %v", p.Path, p.Err))
	}
	for _, p := range foldProblems {
		problems = append(problems, "fold: "+p.String())
	}
	for _, p := range problems {
		fmt.Fprintln(ts.Stdout(), p)
	}
	switch {
	case neg && len(problems) == 0:
		ts.Fatalf("fold-ok: %d entries, no problems", len(l.Entries))
	case !neg && len(problems) > 0:
		ts.Fatalf("fold-ok: %d problems", len(problems))
	}
}

// testGitConfig is the global git config every script runs with.
const testGitConfig = `[user]
	name = Rhiza Test
	email = test@rhiza.invalid
[init]
	defaultBranch = main
[commit]
	gpgsign = false
`

// hermeticGit points git at a sandboxed HOME and global config so scripts
// never read the developer's real configuration. XDG_CACHE_HOME goes there
// too, so os.UserCacheDir, and with it every script's index, is under the
// script's work dir on linux as on darwin (where it is $HOME/Library/Caches)
// and never in the real user cache. A missing git binary is
// left to fail the scripts that use it rather than skipping them.
func hermeticGit(env *testscript.Env) error {
	home := filepath.Join(env.WorkDir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		return err
	}
	gitconfig := filepath.Join(home, ".gitconfig")
	if err := os.WriteFile(gitconfig, []byte(testGitConfig), 0o644); err != nil {
		return err
	}
	env.Setenv("HOME", home)
	env.Setenv("XDG_CACHE_HOME", filepath.Join(home, ".cache"))
	env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	env.Setenv("GIT_CONFIG_GLOBAL", gitconfig)
	return nil
}
