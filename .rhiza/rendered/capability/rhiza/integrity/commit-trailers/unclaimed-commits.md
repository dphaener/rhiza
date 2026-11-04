# unclaimed-commits

capability · `01M6HM30GNXZNEK2YGYEBWFK54` · parent [commit-trailers](../commit-trailers.md)

## acceptance

A testscript with a cutoff commits: one code commit with no trailer, one whose Change: names an unknown ID, one with a known Change: in a paragraph before Co-Authored-By, one shaped like a GitHub multi-commit squash (bullet paragraphs each followed by a known Change:/Actor: paragraph, then a '---------' line, then a Co-authored-by paragraph), one with 'Change: <known ID>' only inside a prose paragraph, one touching only .rhiza/, one with only Actor:/Leaf: trailers, a merge commit, and one untrailered commit before the cutoff; check --json lists exactly one warn finding of the unclaimed rule for the first, second, sixth and seventh, each naming its commit and saying missing or unknown, none for the squash-shaped commit, and exits 0

## statement

With a cutoff set, rhiza check warns once on each non-merge commit in <cutoff>..HEAD that changes a file outside .rhiza/ and has no Change: trailer naming a change set in the log, a Change: trailer being read from any wholly-trailer paragraph of the message after the subject (so a GitHub multi-commit squash body, ending with '---------' and a Co-authored-by paragraph, is claimed by its branch commits' Change: paragraphs); the warning says whether the trailer is missing or names an unknown change set; Actor:, Leaf: and Review: trailers alone do not claim a commit, and parseTrailerLog and the evidence-actor rule keep reading only the trailing trailer run

## Links

### Out

- depends-on [unclaimed-cutoff](unclaimed-cutoff.md)
- decided-by [unclaimed-changes-cutoff](../../../../decision/unclaimed-changes-cutoff.md)

### In

- depends-on [unclaimed-on-here](unclaimed-on-here.md)
- verifies [unclaimed-commits-verified](../../../../evidence/unclaimed-commits-verified.md)
- verifies [unclaimed-commits-qa-1](../../../../evidence/unclaimed-commits-qa-1.md)
- verifies [unclaimed-commits-verified-2](../../../../evidence/unclaimed-commits-verified-2.md)
- verifies [pr-248-unclaimed-changes](../../../../evidence/pr-248-unclaimed-changes.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
