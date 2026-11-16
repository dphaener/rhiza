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
// processScripts: every other script runs rhiza in-process.
func rhizaProcess() int {
	if name := os.Getenv("RHIZA_TEST_SCRIPT"); name != "" {
		if _, ok := processScripts[name]; !ok {
			fmt.Fprintf(os.Stderr, "rhiza: script %s started rhiza as a process but is not on processScripts, the process list in cmd/rhiza/main_test.go: run its rhiza lines in-process (a bare rhiza line), or add it to the list with its reason\n", name)
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
		return run(os.Args[1:], wd, selfPath(), os.Stdin, os.Stdout, os.Stderr, os.Environ())
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
	"apply":                                "rhiza apply reads its file from the process's working directory, not run()'s; a stdin feed",
	"apply-concurrent":                     "rhiza apply reads its file from the process's working directory, not run()'s",
	"arbiter-agent-tabs":                   "fake bin/ rex, claude, gh on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"arbiter-answer":                       "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"arbiter-question":                     "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's; rhiza with &",
	"arbiter-step":                         "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"autonomous-gate-rule":                 "starts rhiza through sh -c",
	"autonomous-permissions":               "fake bin/ rex, claude on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"autonomous-review":                    "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"autonomy-accept-rule":                 "starts rhiza through sh -c",
	"brief-files":                          "fake bin/ rex, claude on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"bugfix-workflow-autonomous":           "fake bin/ rex, claude, gh on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"build-check":                          "starts rhiza through sh -c",
	"build-evidence-commit":                "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"build-leaf-fails":                     "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"build-leaves-status":                  "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"build-tabs":                           "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"build-verifier-checkout":              "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"build-verifier-findings":              "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"build-workflow":                       "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"carve-guards":                         "starts rhiza through sh -c",
	"change-new-workflow-flag":             "starts rhiza through sh -c",
	"change-request-kind":                  "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"change-set-ref":                       "starts rhiza through sh -c",
	"change-show-canonical-text":           "starts rhiza through sh -c",
	"change-show-serves":                   "starts rhiza through sh -c",
	"change-workflow-pin":                  "starts rhiza through sh -c",
	"change-worktree-unique-slug":          "fake bin/ rex, claude on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"check":                                "starts rhiza through sh -c",
	"check-evidence-actor":                 "starts rhiza through sh -c",
	"check-git":                            "starts rhiza through sh -c",
	"check-move":                           "starts rhiza through sh -c",
	"ci-install":                           "starts rhiza through sh -c",
	"custom-workflows-doc":                 "quickstart-doc runs the doc's rhiza commands as processes",
	"daemon-continue-slot":                 "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-headless-runs":                 "fake bin/ rex on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-install":                       "fake bin/ systemctl, loginctl on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-load-hold":                     "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-lock":                          "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's; rhiza with &",
	"daemon-max-runs-pending":              "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-milestone-request":             "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-no-root":                       "starts rhiza through sh -c",
	"daemon-notify":                        "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-notify-hook":                   "fake bin/ rex, claude on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's; rhiza with &",
	"daemon-pause":                         "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-picks":                         "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-picks-backlog":                 "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-plans-proposals":               "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-projects":                      "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's; rhiza with &",
	"daemon-queued-start-slot":             "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-repo-session":                  "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-request-slot":                  "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-resume":                        "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-resume-landed":                 "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-resume-landed-plan":            "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-resume-slot":                   "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-rex-session":                   "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-root-ship-workflow":            "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-round-robin":                   "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-settings":                      "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-settings-scope":                "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-start-stop":                    "fake bin/ systemctl, loginctl, rex on its PATH, calling rhiza back; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-status":                        "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-status-projects":               "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's",
	"daemon-version-in-health":             "fake bin/ rex on its PATH; in-process run() would look them up on the test binary's PATH and start the host's; rhiza with &",
	"describe":                             "a stdin feed",
	"e2e":                                  notConverted,
	"evidence-actor-show":                  notConverted,
	"feed":                                 notConverted,
	"feed-follow":                          notConverted,
	"field-schema":                         notConverted,
	"findings-concurrent":                  notConverted,
	"findings-concurrent-rounds":           notConverted,
	"frame":                                notConverted,
	"free-bytes-default":                   notConverted,
	"guide":                                notConverted,
	"guide-install":                        notConverted,
	"headless-run-agent-tabs":              notConverted,
	"headless-run-blocked-notice":          notConverted,
	"health":                               notConverted,
	"health-ci":                            notConverted,
	"health-load-hold":                     notConverted,
	"health-machine-part":                  notConverted,
	"health-projects":                      notConverted,
	"health-risk-warning":                  notConverted,
	"health-runs-cap":                      notConverted,
	"hook-commands":                        notConverted,
	"import-autonomous-decisions":          notConverted,
	"import-friction-recorded":             notConverted,
	"import-owner-credited-decisions":      notConverted,
	"import-records-built":                 notConverted,
	"import-retires-docs":                  notConverted,
	"import-unbuilt-as-drafts":             notConverted,
	"inbox-advisor":                        notConverted,
	"inbox-api-error":                      notConverted,
	"inbox-blocker-advice":                 notConverted,
	"inbox-blocking-questions":             notConverted,
	"inbox-disk-low":                       notConverted,
	"inbox-gate-questions":                 notConverted,
	"inbox-landed-plan-run":                notConverted,
	"inbox-landed-resumable":               notConverted,
	"inbox-list":                           notConverted,
	"inbox-load-hold":                      notConverted,
	"inbox-main-stale":                     notConverted,
	"inbox-open":                           notConverted,
	"inbox-owner-run-accept":               notConverted,
	"inbox-permission":                     notConverted,
	"inbox-projects":                       notConverted,
	"inbox-question-advisor":               notConverted,
	"inbox-signoff":                        notConverted,
	"inbox-stale":                          notConverted,
	"inbox-stalled":                        notConverted,
	"inbox-uncovered-proposal":             notConverted,
	"init":                                 notConverted,
	"leaf-cycle":                           notConverted,
	"leaf-order":                           notConverted,
	"lifecycle":                            notConverted,
	"link":                                 notConverted,
	"log":                                  notConverted,
	"log-entries-read-only":                notConverted,
	"log-paged":                            notConverted,
	"ls":                                   notConverted,
	"machine-show":                         notConverted,
	"machine-trust-bypass":                 notConverted,
	"machine-trust-unwritten":              notConverted,
	"model-commit-untracked":               notConverted,
	"mv":                                   notConverted,
	"old-folder":                           notConverted,
	"overlaps":                             notConverted,
	"owner-run-gate-accept":                notConverted,
	"plan-milestone-backlog":               notConverted,
	"plan-milestone-check":                 notConverted,
	"plan-milestone-interview-waiting":     notConverted,
	"plan-milestone-number":                notConverted,
	"plan-milestone-workflow":              notConverted,
	"plan-requests":                        notConverted,
	"plan-review-visit-round":              notConverted,
	"plan-review-workflow":                 notConverted,
	"plan-run-milestone-env":               notConverted,
	"plan-workflow":                        notConverted,
	"plan-workflow-autonomous":             notConverted,
	"plan-workflow-autonomous-left":        notConverted,
	"plan-workflow-ci-flake":               notConverted,
	"plan-workflow-ci-flake-red":           notConverted,
	"plan-workflow-ci-infra-timeout":       notConverted,
	"plan-workflow-ci-infra-wait":          notConverted,
	"plan-workflow-merge-main-moved":       notConverted,
	"plan-workflow-merge-main-moved-often": notConverted,
	"plan-workflow-open":                   notConverted,
	"plan-workflow-pr":                     notConverted,
	"policy-check":                         notConverted,
	"printout-counts-used-tokens":          notConverted,
	"printout-gap-names-reason":            notConverted,
	"project-add-path":                     notConverted,
	"project-add-url":                      notConverted,
	"project-agent":                        notConverted,
	"project-ls-rm":                        notConverted,
	"question-answer":                      notConverted,
	"quickstart":                           notConverted,
	"ready":                                notConverted,
	"record":                               notConverted,
	"render":                               notConverted,
	"render-generate":                      notConverted,
	"render-refuses-damaged-log":           notConverted,
	"reopen":                               notConverted,
	"reopen-base":                          notConverted,
	"repair":                               notConverted,
	"request-ref":                          notConverted,
	"request-ref-risk-note":                notConverted,
	"request-workflow-flag":                notConverted,
	"retire":                               notConverted,
	"retire-after-update":                  notConverted,
	"review-check-worktree":                notConverted,
	"review-checks-step":                   notConverted,
	"review-code-check":                    notConverted,
	"review-code-check-override":           notConverted,
	"review-findings":                      notConverted,
	"review-signoff":                       notConverted,
	"review-signoff-autonomous":            notConverted,
	"review-signoff-changes":               notConverted,
	"review-summary":                       notConverted,
	"review-workflow":                      notConverted,
	"rex-install":                          notConverted,
	"rhiza-root":                           notConverted,
	"roadmap":                              notConverted,
	"roadmap-milestones":                   notConverted,
	"roadmap-paged":                        notConverted,
	"roadmap-projects":                     notConverted,
	"roadmap-statement":                    notConverted,
	"roadmap-umbrella":                     notConverted,
	"root-init":                            notConverted,
	"run-abandon":                          notConverted,
	"run-abandon-dirty":                    notConverted,
	"run-abandon-failed":                   notConverted,
	"run-abandon-shared":                   notConverted,
	"run-abandon-unmerged":                 notConverted,
	"run-abandon-unpushed":                 notConverted,
	"run-agent-done":                       notConverted,
	"run-agent-step":                       notConverted,
	"run-agent-tabs":                       notConverted,
	"run-amend":                            notConverted,
	"run-arbitrate-asked":                  notConverted,
	"run-autonomous":                       notConverted,
	"run-autonomy-source":                  notConverted,
	"run-blocked":                          notConverted,
	"run-cleanup-steps":                    notConverted,
	"run-command-step":                     notConverted,
	"run-command-step-output":              notConverted,
	"run-continue-queue":                   notConverted,
	"run-end-removes-leftovers":            notConverted,
	"run-failure-classifier":               notConverted,
	"run-gate":                             notConverted,
	"run-handback":                         notConverted,
	"run-journal":                          notConverted,
	"run-log":                              notConverted,
	"run-ls":                               notConverted,
	"run-max-visits":                       notConverted,
	"run-nested-workflow-session":          notConverted,
	"run-one-live":                         notConverted,
	"run-records-dispatch":                 notConverted,
	"run-resume":                           notConverted,
	"run-resume-agent-session":             notConverted,
	"run-resume-landed-past-land":          notConverted,
	"run-resume-landed-plan":               notConverted,
	"run-resume-lost-fixer":                notConverted,
	"run-resume-lost-session":              notConverted,
	"run-resume-retry":                     notConverted,
	"run-resume-retry-wait":                notConverted,
	"run-resume-unaccepted":                notConverted,
	"run-review-check":                     notConverted,
	"run-review-check-handles":             notConverted,
	"run-review-check-landed":              notConverted,
	"run-review-journal":                   notConverted,
	"run-send-back":                        notConverted,
	"run-ship-change":                      notConverted,
	"run-ship-change-fail":                 notConverted,
	"run-ship-change-rebase":               notConverted,
	"run-ship-ci-fix":                      notConverted,
	"run-ship-ci-lint-fix":                 notConverted,
	"run-ship-conflict-fix":                notConverted,
	"run-ship-conflict-fix-more":           notConverted,
	"run-ship-deploy-check":                notConverted,
	"run-ship-deploy-timeout":              notConverted,
	"run-ship-infra-wait":                  notConverted,
	"run-ship-land-conflict":               notConverted,
	"run-ship-landed-infra-wait":           notConverted,
	"run-ship-main-dispatch":               notConverted,
	"run-ship-main-dispatch-fail":          notConverted,
	"run-ship-main-dispatch-resume":        notConverted,
	"run-ship-merge-draft":                 notConverted,
	"run-ship-merge-main-moved":            notConverted,
	"run-ship-model-conflict":              notConverted,
	"run-ship-needs-leaves":                notConverted,
	"run-ship-rerun-fail":                  notConverted,
	"run-ship-rerun-pass":                  notConverted,
	"run-ship-retry-ci-triage":             notConverted,
	"run-ship-retry-pending":               notConverted,
	"run-ship-rspec-failure":               notConverted,
	"run-ship-signoff":                     notConverted,
	"run-signoff-workflow-owner":           notConverted,
	"run-spawn":                            notConverted,
	"run-start":                            notConverted,
	"run-start-agent-session":              notConverted,
	"run-start-milestone-autonomous":       notConverted,
	"run-start-new":                        notConverted,
	"run-start-new-main-clean":             notConverted,
	"run-start-queue":                      notConverted,
	"run-start-ready-dispatch":             notConverted,
	"run-start-skill-missing":              notConverted,
	"run-status":                           notConverted,
	"run-status-messages":                  notConverted,
	"run-status-stalled":                   notConverted,
	"run-status-tokens":                    notConverted,
	"run-tail":                             notConverted,
	"run-worktree-carries-accept":          notConverted,
	"runs-need-root":                       notConverted,
	"schedules-check":                      notConverted,
	"search":                               notConverted,
	"session-done-captures":                notConverted,
	"session-hook-log-guard":               notConverted,
	"session-hook-message":                 notConverted,
	"session-interview":                    notConverted,
	"session-interview-landed":             notConverted,
	"session-interview-recorded":           notConverted,
	"session-interview-source":             notConverted,
	"session-journal":                      notConverted,
	"session-message":                      notConverted,
	"session-message-busy":                 notConverted,
	"session-message-idle-fallback":        notConverted,
	"session-message-permission":           notConverted,
	"session-pane-key":                     notConverted,
	"session-run":                          notConverted,
	"session-run-blocked-notice":           notConverted,
	"session-run-tokens":                   notConverted,
	"session-run-usage-limit":              notConverted,
	"session-start":                        notConverted,
	"set":                                  notConverted,
	"ship-change-one-ci":                   notConverted,
	"ship-check-ci-test-findings":          notConverted,
	"ship-check-follows-moved-head":        notConverted,
	"ship-check-settles-conflict":          notConverted,
	"ship-close-out-forget-retry":          notConverted,
	"ship-close-out-log-copies":            notConverted,
	"ship-deploy-settings":                 notConverted,
	"ship-follows-pr-head":                 notConverted,
	"ship-gate-head-check":                 notConverted,
	"ship-land-branch-build":               notConverted,
	"ship-local-gate":                      notConverted,
	"ship-main-dispatch-select":            notConverted,
	"ship-prepush-check":                   notConverted,
	"show":                                 notConverted,
	"skill-check":                          notConverted,
	"skill-check-roles":                    notConverted,
	"slice":                                notConverted,
	"slice-budget":                         notConverted,
	"slice-change":                         notConverted,
	"slice-errors":                         notConverted,
	"source-add":                           notConverted,
	"source-raw-session":                   notConverted,
	"stacked":                              notConverted,
	"stacked-cycle":                        notConverted,
	"stacked-dropped":                      notConverted,
	"stacked-revised":                      notConverted,
	"stacked-revised-disjoint":             notConverted,
	"stacked-revised-drop":                 notConverted,
	"stacked-revised-rewrite":              notConverted,
	"stage-autonomy":                       notConverted,
	"status":                               notConverted,
	"status-terminal-current":              notConverted,
	"testscript-root":                      notConverted,
	"text-refs-display":                    notConverted,
	"text-refs-engine-writes":              notConverted,
	"text-refs-write":                      notConverted,
	"tree":                                 notConverted,
	"tui-autonomy-app":                     notConverted,
	"tui-inbox-app":                        notConverted,
	"tui-launch-usage":                     notConverted,
	"tui-old-apps-gone":                    notConverted,
	"tui-open-usage":                       notConverted,
	"tui-request-app":                      notConverted,
	"tui-runs-app":                         notConverted,
	"unknown":                              notConverted,
	"usage":                                notConverted,
	"usage-hold":                           notConverted,
	"usage-live-record":                    notConverted,
	"usage-settings":                       notConverted,
	"version":                              notConverted,
	"version-channel":                      notConverted,
	"wait-timeout":                         notConverted,
	"watchdog-health":                      notConverted,
	"watchdog-policy":                      notConverted,
	"why":                                  notConverted,
	"why-paged":                            notConverted,
	"workflow-check":                       notConverted,
	"workflow-check-invoke":                notConverted,
	"workflow-check-loop":                  notConverted,
	"workflow-check-merge":                 notConverted,
	"workflow-check-risk":                  notConverted,
	"workflow-describe":                    notConverted,
	"workflow-env":                         notConverted,
	"workflow-invoke":                      notConverted,
	"workflow-ls":                          notConverted,
	"workflow-new":                         notConverted,
	"workflow-personal-override":           notConverted,
	"workflow-plan-autonomy":               notConverted,
	"workflow-stage-autonomy":              notConverted,
}

const notConverted = "not converted yet"

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

// rhizaInProcess is the script command `rhiza args...`: it runs run() in
// the test process with the script's working directory, environment and
// streams, self as its own path, and fails or passes as `exec rhiza` would.
func rhizaInProcess(ts *testscript.TestScript, neg bool, args []string, self string) {
	if len(args) > 0 && strings.HasPrefix(args[len(args)-1], "&") {
		ts.Fatalf("rhiza: an in-process rhiza cannot run in the background; use exec rhiza ... & and put the script on processScripts")
	}
	var mu sync.Mutex
	stdout := &lockedWriter{mu: &mu, w: ts.Stdout()}
	stderr := &lockedWriter{mu: &mu, w: ts.Stderr()}
	code := run(args, ts.MkAbs("."), self, strings.NewReader(""), stdout, stderr, scriptEnv(ts))
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
	// The rhiza process entry checks the script against processScripts.
	env.Setenv("RHIZA_TEST_SCRIPT", scriptWorkName(env.WorkDir))
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
