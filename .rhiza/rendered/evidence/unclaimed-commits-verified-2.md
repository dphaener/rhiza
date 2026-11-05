# unclaimed-commits-verified-2

evidence · `01M6KM9VQYA5EW7QJN7G2NM3W1` · actor [claude-code-f59ef46b](../actor/claude-code-f59ef46b.md)

## result

PASS

## statement

Verifier (claude-code-f59ef46b) at abe1e4f8 (attempt 2: fix for the quality finding, adds --full-history to gitUnclaimed). Ran the leaf's acceptance script `go test ./cmd/rhiza -run TestScripts/check-unclaimed -count=1`: check-unclaimed-commits passes. check --json exits 0 with ok true and exactly 4 warn findings of the unclaimed rule: the no-trailer commit (missing), the unknown-ID commit (unknown change set 01M5ZZZ...), the prose-only commit and the Actor/Leaf-only commit. There are none for the paragraph-before-Co-Authored-By commit, the GitHub-squash-shaped commit, the .rhiza/-only commit, the merge or the pre-cutoff commit. (The acceptance's 'sixth' is a miscount for the fifth, the prose-only commit, as compliance already noted. The statement excludes .rhiza/-only commits.) check-unclaimed-cutoff and the new check-unclaimed-full-history pass. To test the new script, I removed --full-history in a scratch worktree, and it then fails (1 unclaimed match, want 3), so it pins the fix. The diff is the one flag plus a doc comment and the script. make fmt-check passes. make test passes for every package (cmd/rhiza, internal/api, internal/check ...), except the TestReadBudget timing budget, run with load avg ~23 on 4 cores (read path untouched). make self-check fails only on stale .rhiza/rendered actor files from earlier review-run 'record writes' commits. HEAD changes no .rhiza/ file.

## Links

### Out

- verifies [unclaimed-commits](../capability/rhiza/integrity/commit-trailers/unclaimed-commits.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
