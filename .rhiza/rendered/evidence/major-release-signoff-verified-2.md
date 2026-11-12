# major-release-signoff-verified-2

evidence · `01M78R3ETMC50J2XNYV0V1K1VW` · actor [claude-code-69c2bc2d](../actor/claude-code-69c2bc2d.md)

## result

PASS

## statement

Verified at 4d496116f, the fix for the confirmed blocker that a signed-off major was never cut. go test -trimpath -count=1 ./tools/release/... -run 'TestMajor|TestCutRealGit' -v: PASS. TestMajor holds the major as a draft while unsigned. Once the draft is published, the fake fetch lands the lightweight v2.0.0 in the fixture repo, and the run tags and pushes once. A plain tag pushes once, a published tag at another commit is never moved, and a failing gate blocks with no draft. TestCutRealGitKeepsHeadings/replacing_a_published_lightweight_tag uses real git with a bare origin and its fetch. It ends with origin's v2.0.0 as an annotated tag that keeps the '## Breaking changes' heading. With next.go and breaking.go reverted to HEAD~1 in a scratch worktree, both tests fail with 'none: no release since v2.0.0'. The full go test -trimpath ./tools/release/... passes. In the diff, evaluate and breakingAgainstLastTag call describe without --tags, so they read annotated tags only.

## Links

### Out

- verifies [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
