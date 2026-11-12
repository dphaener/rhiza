# release-next-version-verified-2

evidence · `01M78752NDJRVJ0A85NCPQK222` · actor [claude-code-54b27abc](../actor/claude-code-54b27abc.md)

## result

PASS

## statement

At 777aa4c32 ran go test -trimpath -count=1 ./tools/release/... -run 'TestNextVersion|TestReadHistory' -v: all pass, TestNextVersion covers major from v1, breaking-before-v1-is-minor, milestone minor, fixes patch, fixes beside unfinished features still a patch, features alone none, nothing since the tag none, no tag breaking minor; TestReadHistoryMilestones shows done points and product intents are not milestones and a milestone counts the servers of intents under it, following node.update set/unset and node.move. Whole tools/release package passes, vet and gofmt clean. Read the diff: ReadHistory resolves each served intent to its milestone ancestor (op shapes match real log node.update/node.move). Probe via go test -overlay over the real .rhiza/log: only milestone-N intents appear, and unfinished milestone-26/32 are excluded from Inputs; go run ./tools/release check prints v0.2.0 minor from v0.1.0 (pre-v1 breaking).

## Links

### Out

- verifies [release-next-version](../capability/rhiza/distribution/release-next-version.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
