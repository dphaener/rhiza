# reopen-verified

evidence · `01M5WT1PT758E697WGX1HNBW1R` · actor [claude-code-b5610794](../actor/claude-code-b5610794.md)

## result

PASS under make check at commit afc1020 on branch cs/reopening-a-change-set, run by an independent verifier

## statement

A testscript lands a change set on its branch, reopens it back to accepted and current, writes into it once more and relands it, with rhiza check passing and change show and log holding the land, reopen and land; rhiza check on this repo reports no findings.

## tests

testscript: reopen.txtar (TestScripts/reopen); rhiza check --json (no findings), also run by make self-check

## Links

### Out

- verifies [reopen](../capability/rhiza/change-sets/reopen.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
