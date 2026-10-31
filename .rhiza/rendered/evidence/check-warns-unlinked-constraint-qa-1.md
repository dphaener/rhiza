# check-warns-unlinked-constraint-qa-1

evidence · `01M69X0NMVFD8MC1CYX469FZE4` · actor [claude-code-fc2ab29d](../actor/claude-code-fc2ab29d.md)

## result

PASS

## statement

QA cycle 1 with a binary built from the branch head, in a scratch git repo: constraints linked by constrains to a capability and to an intent gave no warning; one with only relates-to, and one with only a derived-from link to a source, each gave exactly one 'warn link ... constrains nothing' line with a next hint; a retired constraint and an unlanded draft-change-set constraint gave none; check printed ok and exited 0, and --json had ok true with severity warn. Running the next hint (rhiza link loose constrains thing) silenced that warning. rhiza check on this repo stayed ok with no warnings; go test of TestRunUnlinkedConstraint and the check-constraint/check-git/check-move scripts passed.

## Links

### Out

- verifies [check-warns-unlinked-constraint](../capability/rhiza/integrity/check/check-warns-unlinked-constraint.md)

## History

- `01M5SWTG70W2JC17JR7VCMT74R` genesis check warns about a constraint that is linked to nothing.
