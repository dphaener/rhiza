# release-next-version-verified

evidence · `01M785FCRX9YFM9BQX0SPJEF49` · actor [claude-code-92b6a1cb](../actor/claude-code-92b6a1cb.md)

## result

PASS

## statement

At a2142b081, go test -trimpath -count=1 ./tools/release/... -run TestNextVersion -v passed: TestNextVersion (7 subtests: breaking from v1 is major, breaking before v1 is minor, milestone is minor, fixes are a patch, features alone are none, nothing since the tag, no tag breaking is minor; plus bad tag rejected) and TestNextVersionFromLogs (fixture log: bug fix, milestone finished across the tag, unfinished milestone ignored, breaking mark on a landed change set only, minor before v1 and major from v1, none when nothing landed since the tag). go vet and gofmt clean. Read the diff and checked ReadHistory's field names (open.request_kind, land, ops node.create intent handle, link.create serves) against real .rhiza/log entries. Smoke run go run ./tools/release check on this repo printed v0.2.0 minor (since v0.1.0) with two detected breaking setting changes.

## Links

### Out

- verifies [release-next-version](../capability/rhiza/distribution/release-next-version.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
