# daemon-picks-verified

evidence · `01M6281864PKQ2777DGD7AXQBD` · actor [claude-code-eaa50737](../actor/claude-code-eaa50737.md)

## result

PASS

## statement

At 890b190 I ran go test ./internal/api -run TestDaemonPicks -count=1 and it passed. Over a fixture with accepted milestone change sets first and second, a leafless one, a backlog one, a draft one, and one running, one interrupted and one failed run on three more milestone change sets, daemon.Pick picks [first] at max 2, nothing at max 1, and [first second] at max 4. With the failed run abandoned it picks [first] at max 2 and [first second failed] at max 4, and daemonCandidates leaves out the leafless, backlog and draft change sets. go test ./cmd/rhiza -run TestScripts/daemon-picks -count=1 passed: rhiza daemon --once with a fake rex prints exactly one 'Started a run of', the rex log has one 'rex window new ... run start' line, and --max-runs=-1 is refused. From the diff, interrupted runs block the pick but take no slot, failed runs that are not abandoned block it (a retry reopens the same journal), and List sorts runs oldest first so the latest run per change set wins. make fmt-check, go vet and make test all pass. make self-check reports a stale render only for actor claude-code-eaa50737, which is this verifier session's actor from the in-flight change set and not part of the commit.

## Links

### Out

- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
