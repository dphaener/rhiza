# check-warns-unlinked-constraint

capability · `01M69T0MZJZFJNQS12RYP79XYG` · parent [check](../check.md)

## acceptance

A Go test in internal/check with one constraint linked by constrains, one with only relates-to and one with no links gives two warnings and OK true; a check-constraint.txtar shows the warning and exit 0; this repo's make self-check stays clean

## statement

rhiza check warns, under the link rule, once for each live as-is constraint with no live constrains link; derived-from or relates-to links alone still count as unlinked

## Links

### Out

- decided-by [unlinked-constraint-under-link-rule](../../../../decision/unlinked-constraint-under-link-rule.md)

### In

- verifies [check-warns-unlinked-constraint-verified](../../../../evidence/check-warns-unlinked-constraint-verified.md)
- verifies [check-warns-unlinked-constraint-qa-1](../../../../evidence/check-warns-unlinked-constraint-qa-1.md)
- verifies [pr-142-genesis-check-warns-about-a](../../../../evidence/pr-142-genesis-check-warns-about-a.md)

## History

- `01M5SWTG70W2JC17JR7VCMT74R` genesis check warns about a constraint that is linked to nothing.
