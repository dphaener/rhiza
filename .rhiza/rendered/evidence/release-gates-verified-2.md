# release-gates-verified-2

evidence · `01M787B3JWRJ2EH8GHYSWWJG1H` · actor [claude-code-dc8ac117](../actor/claude-code-dc8ac117.md)

## result

PASS

## statement

At b0afd3c6d (detached checkout) I ran go test -trimpath -count=1 ./tools/release/... -run TestGates -v. All 14 TestGates subtests passed: all pass, CI red, CI missing, commit missing from the history (soak fails), soak too short, bug filed since, bug before the soak, milestone in flight, milestone blocked, second patch in 24h, patch after 24h, severe skips soak, severe skips batching and severe never skips CI. They use a fake clock and a fake CI value over the testdata/runs.jsonl run-history fixture, read through ReadRunHistory. TestGatesCheckReports also passed. It drives runCheck end to end in a scratch git repo with a fake CI func and a run-history file, and checks each gate's pass or FAIL line with its reason. The whole tools/release package passed, and go vet and gofmt are clean. I read the diff for review finding f2 (the milestone gate was vacuous). gateInputs now counts a change set serving a finished milestone as in flight when its land entry (Landing.Entry, a path under the log) is not in the candidate commit's tree (git ls-tree). So a land that is uncommitted, or an older --commit, fails the gate. The new test steps show that: an uncommitted CS6 land gives 'milestone FAIL change sets not landed in <head>: CS6', and the gate passes once the land is committed.

## Links

### Out

- verifies [release-gates](../capability/rhiza/distribution/release-gates.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
