# unclaimed-commits-verified

evidence · `01M6KA07FW4XTW68AJ39KSNC9Z` · actor [claude-code-ca51b8bb](../actor/claude-code-ca51b8bb.md)

## result

PASS

## statement

Verifier at 81e0d400 (attempt 2). Ran `go test ./cmd/rhiza -run TestScripts/check-unclaimed-commits -count=1`: PASS. That scenario builds the acceptance's cutoff commits and asserts check --json gives exactly 4 unclaimed warns, exit 0. Repeated it by hand with a scratch binary on a scratch repo with a real change set in its log. check printed exactly 4 'warn unclaimed' lines: 'no trailer' and 'known only in prose' and 'actor only' (each '... has no Change: trailer (missing)'), and 'unknown change' ('names unknown change set 01M5ZZZZZZZZZZZZZZZZZZZZZZ'). Each line names its 12-char hash and subject. There was none for the known-Change-before-Co-Authored-By commit, the GitHub squash (bullets, Change/Actor paragraphs, '---------', Co-authored-by), the .rhiza/-only commit, the merge or its parents, or the pre-cutoff commit; exit 0. Read the diff. Claims come from every wholly-trailer paragraph after the subject. Changes, and the evidence-actor rule that uses them, still read only the trailing run: TestParseTrailerLog keeps squash Changes empty and Claims [CS8 CS8]. The attempt-1 finding is fixed: TestGitCommits now expects Subject; TestRunOrder covers RuleUnclaimed. With RHIZA_* unset: make fmt-check vet PASS; make self-check ok; make test: every package ok (internal/api, internal/check, cmd/rhiza included). Only the TestReadBudget timing budget (roadmap/ready/log what-if over 50ms) failed, at load avg 13-22 on 4 CPUs. It passed at the pre-leaf commit f6244f87 in a lower-load window. It times reads over a generated 50,000-node model, and no timed read reaches the changed code (check.Run, gitUnclaimed, parseTrailerLog), so this is machine contention, not this leaf.

## Links

### Out

- verifies [unclaimed-commits](../capability/rhiza/integrity/commit-trailers/unclaimed-commits.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
