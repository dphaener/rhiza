# auto-release-qa-1

evidence · `01M78QJP4A27HE4ZACMJQZSS2H` · actor [claude-code-ddf76b91](../actor/claude-code-ddf76b91.md)

## result

FAIL

## statement

QA at 903604b3b with the built tools/release binary, real git against a scratch clone and bare origin, a fake gh (CI/release list only) and the real root daemon. record-run against the live daemon appended {commit dbb826dd..., at 2026-11-12T20:16:37Z}; a second run added nothing; with no -runs and RHIZA_RELEASE_RUNS empty it exits 1 with 'no run history file'. cut with every gate passing pushed an annotated v0.2.0 at the candidate commit to origin; a repeat run left origin unchanged; with CI failing it printed the gate and 'v0.2.0 is blocked' and made no tag locally or on origin. systemd-analyze --user verify accepts rhiza-release.timer and .service; agents/vm.md names rhiza-release.timer. FAIL: the annotated tag's message has none of the notes' '## Breaking changes' / '## Milestone ...' headings (git cat-file -p refs/tags/v0.2.0 shows 0 lines starting with #): git tag -a -F strips #-lines under the default cleanup, so the release body that release.yml takes from %(contents:body) loses every section heading.

## Links

### Out

- verifies [auto-release](../capability/rhiza/distribution/auto-release.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
