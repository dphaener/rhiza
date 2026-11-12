# auto-release-verified-2

evidence · `01M78QZ5HSSGVBAPBHAXXX5VKK` · actor [claude-code-4873dec4](../actor/claude-code-4873dec4.md)

## result

PASS

## statement

Verifier at db0802806 (attempt 2): go test -trimpath ./tools/release/... -run TestCut -count=1 passed: TestCut with a fake git tags v0.2.0 annotated on the candidate and pushes exactly once as the last call when every gate passes, returns exists with no push when the version is tagged locally or on the remote, and returns blocked with no tag or push when CI is red, the soak is under 24h or the commit is not in the run history; TestCutRealGitKeepsHeadings cuts with real git and a bare origin, on the plain path and on the tag -a -f replace path, and reads '## ' headings (incl. '## Breaking changes') back via %(contents:body). The 'never for a major' case now lives in TestMajor (moved by the major-release-signoff leaf), which passed: two cuts return sign-off with zero pushes and one draft. go test -trimpath ./tools/release/... -run TestRecordDaemonRun passed (appends {commit,at} only for a new commit). TestReleaseWorkflowAndTimer, TestParseDaemonStatus, the whole tools/release package and go test ./cmd/rhiza -run TestVMDocLessons passed. Diff: cut.go adds --cleanup=verbatim to both tag -a calls, fixing the confirmed review finding (headings stripped). Real git 2.x check: a 'v0.2.0\n\n## Features\n\n- thing' message read back with release.yml's `git tag -l --format='%(contents:body)'` gives only '- thing' without the flag and '## Features\n\n- thing' with it. release.yml writes that annotation body to release-notes.md and passes it to goreleaser --release-notes, skipping non-annotated tags. script/systemd/rhiza-release.timer (OnCalendar=hourly) and rhiza-release.service (WorkingDirectory %h/code/rhiza, record-run as ExecStartPre, then cut) are in the repo and agents/vm.md names rhiza-release.timer. The confirmed lightweight-tag finding (describe --tags takes a published major's lightweight tag as the last release) is on the major-release-signoff leaf, not this one, and is raised as the question major-signoff-lightweight-tag.

## Links

### Out

- verifies [auto-release](../capability/rhiza/distribution/auto-release.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
