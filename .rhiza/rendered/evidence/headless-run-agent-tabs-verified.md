# headless-run-agent-tabs-verified

evidence · `01M6FZAYAHF4T2PJ6EN747RXKJ` · actor [claude-code-5806efe4](../actor/claude-code-5806efe4.md)

## result

PASS

## statement

At f6dc3a5f I ran go test ./cmd/rhiza -run 'TestScripts/(headless-run-agent-tabs|run-agent-step|run-blocked|build-tabs|run-agent-tabs|arbiter-agent-tabs|review-checks-step|daemon-headless-runs)$' and it passed. headless-run-agent-tabs.txtar uses a fake rex and drives a headless run: rhiza daemon --once starts a detached rhiza run start --foreground child with REX_SESSION unset. Its -v log shows 4 rex window new calls, each with -s rhiza-repo and titled 'repo: Built · run <run id> · build · builder|verifier a|b #1'. No call contains split, so none is anchored on a driver block. The closed blocks sort equal to the opened ones (b1-b4), so each session closed only its own tab. A second --foreground run with REX_SESSION=session:mine still opened both tabs in rhiza-repo. In the diff, headlessSession (internal/api/run.go) calls ensureRepoSession on --foreground start and resume, and session.RunTabTitle names the run and step; this answers attempt 1's findings f1-f3. go test ./internal/... passed except TestGeneratedModel (internal/tui/modelapp), a read-timing check (median ~100ms, want <50ms) that fails the same way at the parent commit f2704ea1 under load average ~19. It is unrelated to this leaf.

## Links

### Out

- verifies [headless-run-agent-tabs](../capability/rhiza/runtime/rex-sessions/headless-run-agent-tabs.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
