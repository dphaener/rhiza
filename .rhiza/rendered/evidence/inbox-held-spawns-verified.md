# inbox-held-spawns-verified

evidence · `01M7QJ90CRSFP9TKF8CAHWSHSX` · actor [claude-code-aafdf611](../actor/claude-code-aafdf611.md)

## result

PASS

## statement

At c9080a8cc, re-ran the Checks: block under a capped scope: go test -trimpath -p 2 ./internal/api -run 'Inbox|NeedsYou|SpawnHeld' exit 0 (TestInboxSpawnHeldItems covers child and step-session holds, release by spawn, ended and unheld runs, a stable ID and the human output); go test -count=1 -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(describe-golden|inbox-spawn-held)$' PASS for both (inbox-spawn-held checks --json and human output with the detail, a step-session hold, release by spawn and session entries, ended and unlocked runs); go vet ./internal/api exit 0; rhiza ship generate exit 0 with 0 changed files; rhiza render exit 0, its only drift being the verifier's own actor. The diff adds spawnHeldItems in inboxOne only, so the daemon notices and needs-you count skip it, and leaves the load-high, memory-low, disk-low and tmp-low items unchanged. Each hold gets one ID hashed from the run, step, parent and start time. describe-golden is regenerated. No stubs or TODOs.

## Links

### Out

- verifies [inbox-held-spawns](../capability/rhiza/orchestration/inbox/inbox-list/inbox-held-spawns.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
