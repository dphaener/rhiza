# unclaimed-changes-cutoff

decision · `01M6HM078ED04YAZ7XD9WRQZ15`

## reasoning

Location: store.Open rejects any repo.json that is not exactly canonical format 1 ({format,id,remotes}). A new key there would need a format bump, or would make every older rhiza binary fail to open the repo, including the shared ~/bin slot and other worktrees. A separate optional file keeps it committed (CI sees it), and older binaries just skip the warn-only rule. Claimed: the change-trailers decision says commits carry Change: <id> and a missing trailer flags an unclaimed change, so the Change: trailer is what claims a commit. Requiring the ID to resolve catches typos at no cost. Where the trailer is read: plan review round 1 (f1) found that GitHub's multi-commit squash body ends with '---------' and a Co-authored-by paragraph, so trailerRun, which stops at the first paragraph that is not wholly trailers, gives parseTrailerLog no Change: for main commits 2449d7a3, 5542390b, bb443656 and ff72e605; only single-commit squashes (6eaad5d4, ebee728a) keep Change: in the trailing paragraph. Reading Change: from every wholly-trailer paragraph after the subject claims both shapes with no change to how the repo merges, while prose lines like 'see Change: X' inside a mixed paragraph still do not count. The evidence-actor rule keeps the trailing run, because it must not credit an earlier branch commit's Actor:. Merge commits (origin/main merged into PR branches, the usual way to resolve conflicts here) carry no trailer and add no unclaimed work of their own, so counting them would only add noise. Turn on here: the recent main history checked (2449d7a3, 5542390b, bb443656, ff72e605, 6391d25a, 585ecf3a, 0a46afe2, 7faa706f, 6eaad5d4, ebee728a) carries a Change: trailer paragraph on every commit, in the trailing run or in the squash body, so pinning the cutoff to the branch's main base raises no warnings today. The exact squash commit the change set lands as cannot be known before merge. Any main commits that land in between carry the same trailers, so the base commit is a safe stand-in.

## statement

(1) The cutoff lives in a new committed file .rhiza/check.json, {"unclaimed_since": "<full commit hash>"}; repo.json stays the strict format-1 file, untouched. No file, or no key, means the unclaimed-changes rule is off; a hash that does not resolve is a check error naming the file. (2) The rule looks at the commits in <cutoff>..HEAD that change a file outside .rhiza/ (the same pathspec the evidence-actor rule uses), skipping merge commits (--no-merges: their content arrives through commits that are checked themselves). Such a commit is claimed when a Change: trailer in any trailer paragraph of its message after the subject (a paragraph that is wholly git trailers, as trailerParagraph reads one, not only the trailing run that parseTrailerLog reads) names a change set in the log. This is what claims a GitHub multi-commit squash merge, whose body lists each branch commit's Change: paragraph and ends with a '---------' line and a Co-authored-by paragraph. parseTrailerLog and the evidence-actor rule keep reading only the trailing run. Actor:, Leaf: or Review: trailers alone do not claim a commit. Each unclaimed commit gets one warning (severity warn, not error), and the message says whether the trailer is missing or names an unknown change set. (3) Yes: this change set turns the rule on for this repo by committing .rhiza/check.json, with the cutoff set to the origin/main commit the branch is built on when the file is written.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-ab09ebc9-1](../question/question-ab09ebc9-1.md)
- derived-from [plan-unclaimed-changes](../source/plan-unclaimed-changes.md)

### In

- decided-by [commit-trailers](../capability/rhiza/integrity/commit-trailers.md)
- decided-by [unclaimed-cutoff](../capability/rhiza/integrity/commit-trailers/unclaimed-cutoff.md)
- decided-by [unclaimed-commits](../capability/rhiza/integrity/commit-trailers/unclaimed-commits.md)
- decided-by [unclaimed-on-here](../capability/rhiza/integrity/commit-trailers/unclaimed-on-here.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
