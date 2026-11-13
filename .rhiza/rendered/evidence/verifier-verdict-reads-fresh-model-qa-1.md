# verifier-verdict-reads-fresh-model-qa-1

evidence · `01M79705CEX83KBAW6N3288717` · actor [claude-code-a371af78](../actor/claude-code-a371af78.md)

## result

PASS

## statement

Built rhiza at b9d9f97e3 (go build ./cmd/rhiza) and, three times, drove a real two-leaf build run with it in a hermetic scratch Rhiza root (own HOME, git config, rhiza.toml; fake rex and claude copied from headless-run-agent-tabs.txtar, the fake claude verifier writing Evidence via rhiza add evidence + rhiza link verifies): rhiza run start --permission-mode manual --foreground --workflow ship-change. All three runs succeeded with both leaves 'passed: attempt 1: verifier passed', the two verifiers writing concurrently and no 'verifier ended with neither Evidence nor findings'. (The daemon --once path was held by the host's real load guard, so the run was started in the foreground.) Also go test -trimpath ./internal/api -run TestEvidenceReadsPastSessionWrites ok and go test -trimpath ./cmd/rhiza -run 'TestScripts/headless-run-agent-tabs$' -count=3 ok.

## Links

### Out

- verifies [verifier-verdict-reads-fresh-model](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step/verifier-verdict-reads-fresh-model.md)

## History

- `01M7942GRFNQ82YV0V8F4A9VNA` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172; reproduce it locally, repeatedly, then fix it or quarantine it
