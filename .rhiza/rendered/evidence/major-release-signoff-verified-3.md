# major-release-signoff-verified-3

evidence · `01M78S060VMS7KSKACN38T5DC6` · actor [claude-code-17e982f7](../actor/claude-code-17e982f7.md)

## result

PASS

## statement

Verifier, attempt 5, at 6ebde15cd. Ran the acceptance `go test -trimpath -count=1 ./tools/release/... -run TestMajor`: PASS. TestMajor shows two cuts return SignOff with no push and exactly one draft at the candidate; after the release is published (fake git models the lightweight tag publishing makes on the remote) one cut tags and pushes once, with the tag message the version then the notes; a repeat returns Existing; a published tag at another commit errors and pushes nothing; a failing gate blocks with no draft. Also ran the whole package `go test -trimpath -count=1 ./tools/release/...`: ok, incl. TestCutRealGitKeepsHeadings, which with real git and a bare origin replaces a published lightweight v2.0.0 by an annotated tag (tag -a -f, forced push) whose body keeps '## Breaking changes'. Read cut.go: the lightweight tag of a major is not treated as already released, the sign-off is majorSignOff (published GitHub release), and next.go/breaking.go run git describe without --tags, so only annotated v* tags count as releases (Decision release-annotated-tags-only, which answers question major-signoff-lightweight-tag). make self-check from a clean detached worktree at this commit with env -i: ok (in the run env it fails only on this verifier's own untracked actor).

## Links

### Out

- verifies [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
