# plan-source-verified

evidence · `01M5XMSYRWZJ3N0S9298EKZ3N9` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

Testscript source-raw-session launches a fake agent whose env log has RHIZA_SESSION=<session id>; source add --raw-session gives raw_sha256 equal to shasum -a 256 of the fake transcript and a session field naming the session, and grep finds neither the transcript's marker text nor its path under .rhiza/; an unknown session is not_found, exit 4; --raw-session with --raw or --raw-sha256 is usage, exit 2, and nothing is written. TestPlanText (plan-source) finds --raw-session "$RHIZA_SESSION" in rhiza guide plan; describe-golden passes. Hands-on: in a scratch change set the scratch build's source add --raw-session 00000000-... exited 4 (no session matches), with --raw too exited 2, and show plan stayed not found; guide plan printed --raw-session "$RHIZA_SESSION".

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/source-raw-session$'; go test -count=1 ./cmd/rhiza -run 'TestScripts/source-add$'; go test -count=1 ./cmd/rhiza -run 'TestScripts/describe-golden$'; go test -count=1 ./internal/guide -run 'TestPlanText'; hands-on: rhiza source add --raw-session (unknown -> exit 4, with --raw -> exit 2)

## Links

### Out

- verifies [plan-source](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-source.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
