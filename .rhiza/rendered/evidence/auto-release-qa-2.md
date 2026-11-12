# auto-release-qa-2

evidence · `01M78RFJ6X6SGX0M3W2SHDYGDQ` · actor [claude-code-03513d81](../actor/claude-code-03513d81.md)

## result

PASS

## statement

QA at 8f0616278 with the built tools/release binary, real git against a scratch clone and bare origin, a fake gh (CI and release list/create only) and the real root daemon. The real log has bug requests within the last 24h, so the soak gate failed there. On a scratch commit that drops them, check showed v0.2.0 minor with every gate passing. cut then pushed an annotated v0.2.0 at the candidate. Its tag object carries all 5 '## ' headings (Breaking changes and the milestones), the cycle-1 failure, now fixed. release.yml's notes step, run verbatim in a fresh clone, printed annotated=true and wrote release-notes.md with the 5 headings. A repeat cut printed 'none: no release since v0.2.0' and left origin unchanged. With CI failing it printed the ci FAIL gate and 'v0.2.0 is blocked', and no tag existed locally or on origin. record-run against the live daemon appended {commit dbb826dd..., at 2026-11-12T20:16:37Z}. A second run added nothing. With no -runs and RHIZA_RELEASE_RUNS empty, it exited 1 with 'no run history file'. systemd-analyze --user verify accepts both units, and agents/vm.md names rhiza-release.timer. go test -run 'TestCut|TestRecordDaemonRun|TestMajor' passed.

## Links

### Out

- verifies [auto-release](../capability/rhiza/distribution/auto-release.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
