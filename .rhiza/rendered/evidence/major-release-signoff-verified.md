# major-release-signoff-verified

evidence · `01M78Q6NRFZ4WKZPJAWHPWM0P9` · actor [claude-code-d2777e42](../actor/claude-code-d2777e42.md)

## result

PASS

## statement

At bd2464cb8 (attempt 2) I ran `go test -trimpath -count=1 ./tools/release/... -run TestMajor -v`: PASS. The whole tools/release package also passes. TestMajor covers these cases. With no sign-off, Cut drafts one GitHub release, returns SignOff and pushes nothing. When the draft is published, the fake makes a lightweight v2.0.0 at the candidate on the remote, and fetch brings it here. Cut then writes the annotated tag (message: version, then notes) over it and force-pushes once. A repeat returns Existing with no push. Signed off with no remote tag, it does a plain tag and push. A published tag at another commit is an error and nothing is pushed. A failing gate holds the major with no draft. This addresses attempt 1's finding f1. I read the diff of cut.go, cut_test.go and release.yml. For a major, only an annotated tag counts as cut. release.yml skips GoReleaser, attestation and the manifest for a lightweight tag. I checked the fake against real git in a scratch bare repo. ls-remote with refs/tags/v and refs/tags/v^{} gives one line for a lightweight tag and an added peeled ^{} line for an annotated one. After fetch, cat-file -t gives commit, and the tag lists. tag -a -f plus push --force replaces it, and the remote then shows it annotated. So the sign-off is waited for, then the tag is made and pushed as for a minor. Side note, not a failure: another clone that still has the lightweight tag fails a plain `git fetch --tags` after the force-push (would clobber existing tag).

## Links

### Out

- verifies [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
