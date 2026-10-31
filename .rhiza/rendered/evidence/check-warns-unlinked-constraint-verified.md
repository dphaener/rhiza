# check-warns-unlinked-constraint-verified

evidence · `01M69WXH5X2TTG8DBJW94CDFYM` · actor [claude-code-3fcc9d2c](../actor/claude-code-3fcc9d2c.md)

## result

PASS

## statement

At d70f621c, go test -count=1 ./internal/check -run TestRunUnlinkedConstraint passes: one constraint linked by constrains, one with only relates-to and one bare give exactly two link-rule warnings (bare, related) and OK true. go test ./cmd/rhiza -run TestScripts/check-constraint (plus the edited check-git and check-move) passes: rhiza check prints ok, one 'warn link ... constraint loose constrains nothing' line with a next hint, exits 0, --json has ok true and no fail; linking it silences the warning. make check exits 0 and make self-check prints only ok with no warnings. Diff read: internal/check/unlinked.go skips retired constraints and counts only live constrains links from the constraint.

## Links

### Out

- verifies [check-warns-unlinked-constraint](../capability/rhiza/integrity/check/check-warns-unlinked-constraint.md)

## History

- `01M5SWTG70W2JC17JR7VCMT74R` genesis check warns about a constraint that is linked to nothing.
