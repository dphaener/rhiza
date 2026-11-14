# builtin-conductor-default-verified

evidence · `01M7C4C5N35ARQCKM5DJTDK7XS` · actor [claude-code-2054dcef](../actor/claude-code-2054dcef.md)

## result

PASS

## statement

At 124c2bbf1 re-ran the Checks: block: go test -trimpath ./internal/workflow, ./internal/api -run 'TestBuildCheck|TestShipRebase|TestDescribe', ./internal/run ./internal/daemon, ./cmd/rhiza -run TestScripts/describe-golden, rhiza ship generate (0 files changed) and rhiza render (no diff), all exit 0. Acceptance: go test -trimpath ./internal/workflow -run TestBuiltinShipChange PASS (asserts the 10 steps in order, conduct checks plan,code,qa, gate kind gate with build-check/review-check/rhiza check/ship gate, every non-pass route of conduct..merge = handback:conduct, merge only pass->main-ci and fail, main-ci until rhiza ship check --main, no non-pass or hand-back route after the merge, close-out/ship/deploy-check commands, handback_cap 3, prompt naming rhiza ship rebase --no-push, no leaves or checks step); go test -trimpath ./internal/api -run TestBuildCheck PASS (fails a leaf with no Evidence); a scratch build's rhiza workflow describe ship-change in a fresh repo (no repo override) lists conduct, gate, signoff, push, wait-ci, merge, main-ci, close-out, ship, deploy-check with those routes (this repo's own .rhiza/workflows/ship-change.toml still shadows it, which is the rhiza-repo-conductor-default leaf). Diff read: ship-change.toml rewrite, workflow.go accepts signoff on.changes = handback:<agent step>, ship rebase --no-push with TestShipRebaseNoPush, regenerated cli.md and describe-golden; no TODOs or stubs.

## Links

### Out

- verifies [builtin-conductor-default](../capability/rhiza/workflows/built-in-workflows/builtin-conductor-default.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
