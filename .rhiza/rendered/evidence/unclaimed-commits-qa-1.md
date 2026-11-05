# unclaimed-commits-qa-1

evidence · `01M6KESJJ44W94983RVX75A36B` · actor [claude-code-dc452748](../actor/claude-code-dc452748.md)

## result

PASS

## statement

QA at 0bee2696 with a scratch build on a scratch repo: cutoff at an untrailered commit, then commits 'no trailer', 'unknown change' (Change: 01M5ZZZ...), 'known before coauthor' (Change/Actor paragraph then Co-Authored-By), a GitHub-squash shape (bullets, two Change/Actor paragraphs, '---------', Co-authored-by), 'prose only' (see Change: <id> in prose), a .rhiza/-only commit, 'actor only' (Actor:/Leaf:), a lowercase change: trailer and a --no-ff merge of a claimed side commit. rhiza check printed ok, exit 0, and exactly 4 'warn unclaimed' lines, each with a 12-char hash and subject: no trailer, prose only, actor only (missing) and unknown change (names unknown change set 01M5ZZZ...). None for the squash, the known commit, .rhiza/-only, the merge, side or the pre-cutoff commit; check --json listed the same 4 with rule unclaimed, severity warn.

## Links

### Out

- verifies [unclaimed-commits](../capability/rhiza/integrity/commit-trailers/unclaimed-commits.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
